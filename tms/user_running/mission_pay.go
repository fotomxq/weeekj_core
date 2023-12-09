package TMSUserRunning

import (
	"database/sql"
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	FinancePayCreate "gitee.com/weeekj/weeekj_core/v5/finance/pay_create"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ArgsUpdateRunPayPrice 修改跑腿单跑腿费用参数
type ArgsUpdateRunPayPrice struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//追加缴费的费用
	Price int64 `db:"run_wait_price" json:"runWaitPrice" check:"price" empty:"true"`
}

// UpdateRunPayPrice 修改跑腿单跑腿费用
func UpdateRunPayPrice(args *ArgsUpdateRunPayPrice) (err error) {
	//意外退出处理
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint("recover, ", e))
			return
		}
	}()
	//启动事务
	var tx *sqlx.Tx
	tx = Router2SystemConfig.MainDB.MustBegin()
	if err != nil {
		return
	}
	//检查当前任务是否存在未完成支付，且支付状态必须是未完成
	var data FieldsMission
	err = tx.Get(&data, "SELECT id, run_pay_id, run_pay_at, run_wait_price FROM tms_user_running_mission WHERE id = $1 AND ($2 < 1 OR user_id = $2)", args.ID, args.UserID)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	// 如果尚未完成支付，且支付ID>0，则检查支付状态
	if data.RunPayID > 0 && data.RunPayAt.Unix() < 1000000 {
		var checkResults []FinancePay.DataCheckFinish
		checkResults, err = FinancePay.CheckFinishByIDs(&FinancePay.ArgsCheckFinishByIDs{
			IDs: []int64{data.RunPayID},
		})
		if err == nil {
			for _, v := range checkResults {
				if v.IsFinish {
					//将当前的价格，移动到已经支付部分
					var result sql.Result
					result, err = tx.Exec("UPDATE tms_user_running_mission SET run_price = run_price + $1, run_wait_price = 0, run_pay_at = to_timestamp(0), run_pay_id = 0 WHERE id = $2", data.RunWaitPrice, data.ID)
					err = CoreSQL.LastRowsAffected(tx, result, err)
					if err != nil {
						return
					}
				}
			}
		}
	}
	//追加价格
	var newLog string
	newLog, err = getLogData(fmt.Sprint("追加配送费金额: ￥", args.Price/100), []int64{})
	if err != nil {
		_ = tx.Rollback()
		return
	}
	var result sql.Result
	result, err = tx.Exec("UPDATE tms_user_running_mission SET update_at = NOW(), run_wait_price = $1, logs = logs || $2 WHERE id = $3", args.Price, newLog, data.ID)
	err = CoreSQL.LastRowsAffected(tx, result, err)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	//执行事务
	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//反馈
	return
}

// ArgsPayRunPay 缴纳跑腿费用参数
type ArgsPayRunPay struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//支付方式
	// 如果为退单，则为付款方式
	PaymentChannel CoreSQLFrom.FieldsFrom `json:"paymentChannel"`
}

// PayRunPay 缴纳跑腿费用参数
func PayRunPay(args *ArgsPayRunPay) (payData FinancePay.FieldsPayType, errCode string, err error) {
	//获取数据包
	var data FieldsMission
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, user_id, org_id, run_pay_id, run_pay_list, run_pay_at, run_wait_price FROM tms_user_running_mission WHERE id = $1 AND ($2 < 1 OR user_id = $2)", args.ID, args.UserID)
	if err != nil || data.ID < 1 {
		errCode = "no_data"
		err = errors.New("no data")
		return
	}
	//检查金额必须>0
	if data.RunWaitPrice < 1 {
		errCode = "price_less_1"
		err = errors.New("price less 1")
		return
	}
	//创建支付请求
	payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
		UserID:         data.UserID,
		OrgID:          data.OrgID,
		IsRefund:       false,
		Currency:       86,
		Price:          data.RunWaitPrice,
		PaymentChannel: args.PaymentChannel,
		ExpireAt:       CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
		Des:            "支付跑腿费用",
	})
	if err != nil {
		return
	}
	//修改任务的支付ID
	var newLog string
	newLog, err = getLogData(fmt.Sprint("发起支付跑腿费用(", data.RunWaitPrice, "元)"), []int64{})
	if err != nil {
		errCode = "log"
		return
	}
	data.RunPayList = append(data.RunPayList, payData.ID)
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET run_pay_id = :run_pay_id, logs = logs || :logs, run_pay_list = :run_pay_list WHERE id = :id", map[string]interface{}{
		"run_pay_id":   payData.ID,
		"run_pay_list": data.RunPayList,
		"id":           data.ID,
		"logs":         newLog,
	})
	if err != nil {
		errCode = "update"
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	return
}

// ArgsUpdateOrderPrice 修改跑腿单订单费用参数
type ArgsUpdateOrderPrice struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//订单费用
	OrderPrice int64 `db:"order_price" json:"orderPrice" check:"price" empty:"true"`
	//跑腿单核对订单数据包
	OrderDesFiles pq.Int64Array `db:"order_des_files" json:"orderDesFiles" check:"ids" empty:"true"`
	//跑腿单追加订单描述
	OrderDes string `db:"order_des" json:"orderDes" check:"des" min:"1" max:"3000" empty:"true"`
}

// UpdateOrderPrice 修改跑腿单订单费用
func UpdateOrderPrice(args *ArgsUpdateOrderPrice) (err error) {
	//修改任务订单费用
	var newLog string
	newLog, err = getLogData(fmt.Sprint("修改订单费用(", args.OrderPrice, "元)"), []int64{})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET order_wait_price = :order_wait_price, order_pay_at = to_timestamp(0), order_pay_id = 0, order_des_files = :order_des_files, order_des = :order_des, logs = logs || :logs WHERE id = :id AND (:role_id < 1 OR role_id = :role_id)", map[string]interface{}{
		"order_wait_price": args.OrderPrice,
		"order_des_files":  args.OrderDesFiles,
		"order_des":        args.OrderDes,
		"logs":             newLog,
		"id":               args.ID,
		"role_id":          args.RoleID,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//反馈
	return
}

// ArgsPayOrder 缴纳订单费用参数
type ArgsPayOrder struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//支付方式
	// 如果为退单，则为付款方式
	PaymentChannel CoreSQLFrom.FieldsFrom `json:"paymentChannel"`
}

// PayOrder 缴纳订单费用
func PayOrder(args *ArgsPayOrder) (payData FinancePay.FieldsPayType, errCode string, err error) {
	//获取数据包
	data := getMissionID(args.ID)
	if data.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//检查金额必须>0
	if data.OrderWaitPrice < 1 {
		errCode = "price_less_1"
		err = errors.New("price less 1")
		return
	}
	//创建支付请求
	payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
		UserID:         data.UserID,
		OrgID:          data.OrgID,
		IsRefund:       false,
		Currency:       86,
		Price:          data.OrderWaitPrice,
		PaymentChannel: args.PaymentChannel,
		ExpireAt:       CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
		Des:            "支付跑腿单订单费用",
	})
	if err != nil {
		return
	}
	//修改任务的支付ID
	var newLog string
	newLog, err = getLogData(fmt.Sprint("发起支付订单费用(", float64(data.OrderWaitPrice)/100, "元)"), []int64{})
	if err != nil {
		errCode = "log"
		return
	}
	data.RunPayList = append(data.RunPayList, payData.ID)
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET order_pay_id = :order_pay_id, logs = logs || :logs, run_pay_list = :run_pay_list WHERE id = :id", map[string]interface{}{
		"order_pay_id": payData.ID,
		"run_pay_list": data.RunPayList,
		"logs":         newLog,
		"id":           data.ID,
	})
	if err != nil {
		errCode = "update"
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//反馈
	return
}

// ArgsPayRunAndOrder 融合发起支付参数
type ArgsPayRunAndOrder struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//支付方式
	// 如果为退单，则为付款方式
	PaymentChannel CoreSQLFrom.FieldsFrom `json:"paymentChannel"`
}

// PayRunAndOrder 融合发起支付
// 支付跑腿单和订单费用
func PayRunAndOrder(args *ArgsPayRunAndOrder) (payData FinancePay.FieldsPayType, errCode string, err error) {
	//获取数据包
	data := getMissionID(args.ID)
	if data.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//合计费用
	waitPrice := data.RunWaitPrice + data.OrderWaitPrice
	//检查金额必须>0
	if waitPrice < 1 {
		errCode = "price_less_1"
		err = errors.New("price less 1")
		return
	}
	//创建支付请求
	payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
		UserID:         data.UserID,
		OrgID:          data.OrgID,
		IsRefund:       false,
		Currency:       86,
		Price:          waitPrice,
		PaymentChannel: args.PaymentChannel,
		ExpireAt:       CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
		Des:            "支付跑腿单订单费用",
	})
	if err != nil {
		return
	}
	//修改任务的支付ID
	var newLog string
	newLog, err = getLogData(fmt.Sprint("发起支付跑腿单和订单费用(", float64(waitPrice)/100, "元)"), []int64{})
	if err != nil {
		errCode = "log"
		return
	}
	data.RunPayList = append(data.RunPayList, payData.ID)
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET run_pay_id = :order_pay_id, order_pay_id = :order_pay_id, logs = logs || :logs, run_pay_list = :run_pay_list WHERE id = :id", map[string]interface{}{
		"order_pay_id": payData.ID,
		"run_pay_list": data.RunPayList,
		"logs":         newLog,
		"id":           data.ID,
	})
	if err != nil {
		errCode = "update"
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//反馈
	return
}

// 完成跑腿订单缴费部分
func payMissionOrder(id int64, haveTMS bool) (err error) {
	//获取跑腿单
	missionData := getMissionID(id)
	if missionData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//更新数据库
	if haveTMS {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET run_wait_price = 0, run_price = run_price + run_wait_price, run_pay_at = NOW(), order_wait_price = 0, order_price = order_price + order_wait_price ,order_pay_at = NOW() WHERE id = :id", map[string]interface{}{
			"id": id,
		})
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET order_wait_price = 0, order_price = order_price + order_wait_price ,order_pay_at = NOW() WHERE id = :id", map[string]interface{}{
			"id": id,
		})
	}
	if err != nil {
		return
	}
	//清理缓冲
	deleteMissionCache(id)
	//通知缴纳成功
	if haveTMS {
		pushNatsStatusUpdate("pay_order", id, "通过订单缴纳所有费用")
	} else {
		pushNatsStatusUpdate("pay_order", id, "缴纳订单部分费用")
	}
	//反馈
	return
}

// 完成跑腿单缴费
func payMissionPay(payID int64) (err error) {
	//获取配置
	var tmsRunningCurPrice int64
	tmsRunningCurPrice, err = BaseConfig.GetDataInt64("TMSRunningCurPrice")
	if err != nil {
		tmsRunningCurPrice = 0
	}
	//获取支付ID相关的数据包
	var dataList []FieldsMission
	dataList, err = getMissionListByPayID(payID)
	if err != nil {
		err = errors.New("no data")
		return
	}
	for _, vMission := range dataList {
		//识别支付类型
		if vMission.RunPayID == payID && vMission.RunPayAt.Unix() < 1000000 {
			//构建日志
			var newLog string
			var servicePrice int64 = 0
			if vMission.RunWaitPrice+vMission.RunPrice > 0 {
				if tmsRunningCurPrice > 0 {
					servicePrice = int64(float64(vMission.RunWaitPrice+vMission.RunPrice) * (float64(tmsRunningCurPrice) / 10000))
				}
			}
			if newLog == "" {
				newLog, err = getLogData(fmt.Sprint("完成跑腿费用支付[", vMission.ID, "]"), []int64{})
			} else {
				newLog, err = getLogData(fmt.Sprint("完成跑腿费用支付[", vMission.ID, "], 平台抽取服务费", servicePrice, "(", float64(tmsRunningCurPrice)/100, "%)"), []int64{})
			}
			if err != nil {
				err = errors.New(fmt.Sprint("get log failed, ", err))
				break
			}
			_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET run_pay_id = 0, run_pay_at = NOW(), run_price = :run_price, run_wait_price = 0, service_price = :service_price, logs = logs || :logs WHERE id = :id", map[string]interface{}{
				"id":            vMission.ID,
				"run_price":     vMission.RunPrice + vMission.RunWaitPrice,
				"service_price": servicePrice,
				"logs":          newLog,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("update failed mission id: ", vMission.ID, ", err: ", err))
				break
			}
		}
		if vMission.OrderPayID == payID && vMission.OrderPayAt.Unix() < 1000000 {
			var newLog string
			newLog, err = getLogData(fmt.Sprint("完成订单支付[", vMission.ID, "]"), []int64{})
			if err != nil {
				err = errors.New(fmt.Sprint("get log failed, ", err))
				break
			}
			_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET order_wait_price = 0, order_pay_id = 0, order_pay_at = NOW(), logs = logs || :logs WHERE id = :id", map[string]interface{}{
				"id":   vMission.ID,
				"logs": newLog,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("update failed mission id: ", vMission.ID, ", err: ", err))
				break
			}
		}
		//清理缓冲
		deleteMissionCache(vMission.ID)
		//通知缴纳成功
		pushNatsStatusUpdate("pay_run", vMission.ID, "缴纳跑腿费部分费用")
	}
	//叠加支付ID
	if err != nil {
		err = errors.New(fmt.Sprint(err, ", pay id: ", payID))
		return
	}
	//反馈
	return
}

// 完成跑腿单缴费
func payMissionFailed(payID int64, logDes string) (err error) {
	//获取支付ID相关的数据包
	var dataList []FieldsMission
	dataList, err = getMissionListByPayID(payID)
	if err != nil {
		err = errors.New("no data")
		return
	}
	var ids pq.Int64Array
	for _, v := range dataList {
		ids = append(ids, v.ID)
	}
	for _, vMission := range dataList {
		//生成日志
		var newLog string
		//识别支付类型
		if vMission.RunPayID == payID && vMission.RunPayAt.Unix() < 1000000 {
			newLog, err = getLogData(fmt.Sprint("支付跑腿单费用失败[", vMission.ID, "], ", logDes), []int64{})
			if err != nil {
				err = errors.New(fmt.Sprint("get log failed, ", err))
				break
			}
		}
		if vMission.OrderPayID == payID && vMission.OrderPayAt.Unix() < 1000000 {
			newLog, err = getLogData(fmt.Sprint("支付订单失败[", vMission.ID, "], ", logDes), []int64{})
			if err != nil {
				err = errors.New(fmt.Sprint("get log failed, ", err))
				break
			}
		}
		if newLog != "" {
			_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET logs = logs || :logs WHERE id = :id", map[string]interface{}{
				"id":   vMission.ID,
				"logs": newLog,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("update failed mission id: ", vMission.ID, ", err: ", err))
				break
			}
		}
	}
	//叠加支付ID
	if err != nil {
		err = errors.New(fmt.Sprint(err, ", pay id: ", payID))
		return
	}
	//反馈
	return
}
