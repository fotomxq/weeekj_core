package TMSUserRunning

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserRole "github.com/fotomxq/weeekj_core/v5/user/role"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateMissionRunner 分配跑腿员到任务参数
type ArgsUpdateMissionRunner struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//描述信息
	Des string `json:"des" check:"des" min:"1" max:"600"`
}

// UpdateMissionRunner 分配跑腿员到任务
func UpdateMissionRunner(args *ArgsUpdateMissionRunner) (err error) {
	var newLog string
	newLog, err = getLogData(args.Des, []int64{})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET update_at = NOW(), take_at = NOW(), role_id = :role_id, logs = logs || :logs WHERE id = :id", map[string]interface{}{
		"id":      args.ID,
		"role_id": args.RoleID,
		"logs":    newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//通知成功
	pushNatsStatusUpdate("pick", args.ID, "跑腿员接单并正在前往取货")
	//反馈
	return
}

// UpdateMissionRunnerSelf 分配跑腿员到任务，跑腿员自己接单
func UpdateMissionRunnerSelf(args *ArgsUpdateMissionRunner) (err error) {
	var newLog string
	newLog, err = getLogData(args.Des, []int64{})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET update_at = NOW(), role_id = :role_id, logs = logs || :logs WHERE id = :id AND role_id < 1", map[string]interface{}{
		"id":      args.ID,
		"role_id": args.RoleID,
		"logs":    newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//通知成功
	pushNatsStatusUpdate("pick", args.ID, "跑腿员接单并正在前往取货")
	//反馈
	return
}

// ArgsUpdateMissionReject 拒绝跑腿单参数
type ArgsUpdateMissionReject struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//描述信息
	Des string `json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// UpdateMissionReject 拒绝跑腿单
func UpdateMissionReject(args *ArgsUpdateMissionReject) (err error) {
	var newLog string
	newLog, err = getLogData(args.Des, []int64{})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET update_at = NOW(), role_id = 0, logs = logs || :logs WHERE id = :id AND role_id = :role_id", map[string]interface{}{
		"id":      args.ID,
		"logs":    newLog,
		"role_id": args.RoleID,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//通知成功
	pushNatsStatusUpdate("reject", args.ID, "跑腿员拒绝接单")
	//反馈
	return
}

// ArgsUpdateMission 修改跑腿单信息参数
type ArgsUpdateMission struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//期望上门时间
	WaitAt string `db:"wait_at" json:"waitAt" check:"isoTime"`
	//物品类型
	GoodType string `db:"good_type" json:"goodType" check:"mark"`
	//取货时间
	TakeAt string `db:"take_at" json:"takeAt"`
	//是否完结
	FinishAt string `db:"finish_at" json:"finishAt"`
	//取货码
	TakeCode string `db:"take_code" json:"takeCode"`
	//跑腿单类型
	// 0 帮我送 ; 1 帮我买; 2 帮我取
	RunType int `db:"run_type" json:"runType" check:"intThan0" empty:"true"`
	//关联组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//关联订单ID
	// 可能没有关联订单
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//跑腿付费时间
	RunPayAt string `db:"run_pay_at" json:"runPayAt" check:"isoTime" empty:"true"`
	//是否完成跑腿费支付
	RunPayID int64 `db:"run_pay_id" json:"runPayID" check:"id" empty:"true"`
	//跑腿费用总计
	// 已经支付的部分
	RunPrice int64 `db:"run_price" json:"runPrice" check:"price" empty:"true"`
	//等待缴纳的费用
	RunWaitPrice int64 `db:"run_wait_price" json:"runWaitPrice" check:"price" empty:"true"`
	//跑腿费是否货到付款
	RunPayAfter bool `db:"run_pay_after" json:"runPayAfter" check:"bool"`
	//订单是否货到付款
	OrderPayAfter bool `db:"order_pay_after" json:"orderPayAfter" check:"bool"`
	//订单费用
	OrderPrice int64 `db:"order_price" json:"orderPrice" check:"price" empty:"true"`
	//订单是否已经支付
	OrderPayAt string `db:"order_pay_at" json:"orderPayAt" check:"isoTime" empty:"true"`
	//订单支付ID
	OrderPayID int64 `db:"order_pay_id" json:"orderPayID" check:"id" empty:"true"`
	//跑腿单描述信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1200" empty:"true"`
	//跑腿单核对订单数据包
	OrderDesFiles pq.Int64Array `db:"order_des_files" json:"orderDesFiles" check:"ids" empty:"true"`
	//跑腿单追加订单描述
	OrderDes string `db:"order_des" json:"orderDes" check:"des" min:"1" max:"3000" empty:"true"`
	//物品重量
	GoodWidget int `db:"good_widget" json:"goodWidget" check:"intThan0" empty:"true"`
	//发货地址
	FromAddress CoreSQLAddress.FieldsAddress `db:"from_address" json:"fromAddress" check:"address_data" empty:"true"`
	//送货地址
	ToAddress CoreSQLAddress.FieldsAddress `db:"to_address" json:"toAddress" check:"address_data" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// UpdateMission 修改跑腿单信息
func UpdateMission(args *ArgsUpdateMission) (err error) {
	var newLog string
	newLog, err = getLogData("修改跑腿任务信息", []int64{})
	if err != nil {
		return
	}
	var waitAt, finishAt, runPayAt, orderPayAt time.Time
	if args.WaitAt != "" {
		waitAt, err = CoreFilter.GetTimeByISO(args.WaitAt)
		if err != nil {
			return
		}
	} else {
		waitAt = CoreFilter.GetNowTime()
	}
	if args.FinishAt != "" {
		finishAt, err = CoreFilter.GetTimeByISO(args.FinishAt)
		if err != nil {
			return
		}
	}
	if args.RunPayAt != "" {
		runPayAt, err = CoreFilter.GetTimeByISO(args.RunPayAt)
		if err != nil {
			return
		}
	}
	if args.OrderPayAt != "" {
		orderPayAt, err = CoreFilter.GetTimeByISO(args.OrderPayAt)
		if err != nil {
			return
		}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET update_at = NOW(), wait_at = :wait_at, good_type = :good_type, take_at = :take_at, finish_at = :finish_at, take_code = :take_code, run_type = :run_type, org_id = :org_id, user_id = :user_id, order_id = :order_id, role_id = :role_id, run_pay_at = :run_pay_at, run_pay_id = :run_pay_id, run_price = :run_price, run_wait_price = :run_wait_price, run_pay_after = :run_pay_after, order_pay_after = :order_pay_after, order_price = :order_price, order_pay_at = :order_pay_at, order_pay_id = :order_pay_id, des = :des, order_des_files = :order_des_files, order_des = :order_des, good_widget = :good_widget, from_address = :from_address, to_address = :to_address, logs = logs || :logs, params = :params WHERE id = :id", map[string]interface{}{
		"id":              args.ID,
		"wait_at":         waitAt,
		"good_type":       args.GoodType,
		"take_at":         args.TakeAt,
		"finish_at":       finishAt,
		"take_code":       args.TakeCode,
		"run_type":        args.RunType,
		"org_id":          args.OrgID,
		"user_id":         args.UserID,
		"order_id":        args.OrderID,
		"role_id":         args.RoleID,
		"run_pay_at":      runPayAt,
		"run_pay_id":      args.RunPayID,
		"run_price":       args.RunPrice,
		"run_wait_price":  args.RunWaitPrice,
		"run_pay_after":   args.RunPayAfter,
		"order_pay_after": args.OrderPayAfter,
		"order_price":     args.OrderPrice,
		"order_pay_at":    orderPayAt,
		"order_pay_id":    args.OrderPayID,
		"des":             args.Des,
		"order_des_files": args.OrderDesFiles,
		"order_des":       args.OrderDes,
		"good_widget":     args.GoodWidget,
		"from_address":    args.FromAddress,
		"to_address":      args.ToAddress,
		"logs":            newLog,
		"params":          args.Params,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//反馈
	return
}

// ArgsTakeMission 确认完成取货参数
type ArgsTakeMission struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
}

// TakeMission 确认完成取货
func TakeMission(args *ArgsTakeMission) (err error) {
	var newLog string
	newLog, err = getLogData("完成跑腿任务取货", []int64{})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET update_at = NOW(), take_at = NOW(), logs = logs || :logs WHERE id = :id AND (:role_id < 1 OR role_id = :role_id)", map[string]interface{}{
		"id":      args.ID,
		"role_id": args.RoleID,
		"logs":    newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//通知成功
	pushNatsStatusUpdate("send", args.ID, "跑腿员拿到货物")
	//反馈
	return
}

// ArgsFinishMission 完成任务参数
type ArgsFinishMission struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//履约跑腿员
	// 用户角色ID
	RoleID int64 `db:"role_id" json:"roleID" check:"id" empty:"true"`
	//是否需要取件码
	NeedTakeCode bool `json:"needTakeCode" check:"mark"`
	//取货码
	TakeCode string `json:"takeCode" check:"mark"`
}

// FinishMission 完成任务
func FinishMission(args *ArgsFinishMission) (err error) {
	//构建日志
	var newLog string
	newLog, err = getLogData("完成跑腿任务，送达客户", []int64{})
	if err != nil {
		return
	}
	//更新数据
	if args.NeedTakeCode {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET update_at = NOW(), finish_at = NOW(), logs = logs || :logs WHERE id = :id AND take_code = :take_code AND (:role_id < 1 OR role_id = :role_id)", map[string]interface{}{
			"id":        args.ID,
			"take_code": args.TakeCode,
			"logs":      newLog,
			"role_id":   args.RoleID,
		})
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET update_at = NOW(), finish_at = NOW(), logs = logs || :logs WHERE id = :id AND (:role_id < 1 OR role_id = :role_id)", map[string]interface{}{
			"id":      args.ID,
			"logs":    newLog,
			"role_id": args.RoleID,
		})
	}
	if err != nil {
		err = errors.New(fmt.Sprint("id: ", args.ID, ", role id: ", args.RoleID, ", update fail:", err))
		return
	}
	//获取数据
	var data FieldsMission
	data = getMissionID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//检查跑腿费，将跑腿费付款给跑腿员
	if data.RunPrice > 0 && data.RunPrice-data.ServicePrice > 0 {
		//获取储蓄配置
		var tmsRunningDepositMark string
		tmsRunningDepositMark, err = BaseConfig.GetDataString("TMSRunningDepositMark")
		if err != nil {
			return
		}
		//减去平台手续费，剩下的资金给与跑腿人员个人账户
		err = UserRole.PayToRole(&UserRole.ArgsPayToRole{
			RoleID:      data.RoleID,
			OrgID:       0,
			DepositMark: tmsRunningDepositMark,
			Currency:    86,
			Price:       data.RunPrice - data.ServicePrice,
			SystemFrom:  "tms_running",
			FromID:      data.ID,
			Des:         fmt.Sprint("跑腿人员发放佣金"),
		})
		if err != nil {
			return
		}
	}
	//清理缓冲
	deleteMissionCache(args.ID)
	//通知成功
	pushNatsStatusUpdate("finish", data.ID, "完成跑腿服务")
	//反馈
	return
}
