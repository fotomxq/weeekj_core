package ServiceOrder

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsUpdateTransportID 修改配送ID参数
type ArgsUpdateTransportID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" emtpy:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//调整说明描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//配送服务系统
	// 0 self 自运营服务; 1 transport 自提; 2 running 跑腿服务; 3 housekeeping 家政服务
	TransportSystem string `db:"transport_system" json:"transportSystem"`
	//配送ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id"`
}

// UpdateTransportID 修改配送ID
func UpdateTransportID(args *ArgsUpdateTransportID) (err error) {
	var newLog string
	newLog, err = getLogData(0, args.OrgBindID, args.TransportSystem, args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), transport_system = :transport_system, transport_id = :transport_id, transport_ids = array_append(transport_ids, :transport_id), logs = logs || :log WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":               args.ID,
		"org_id":           args.OrgID,
		"log":              newLog,
		"transport_system": args.TransportSystem,
		"transport_id":     args.TransportID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update order failed, id: ", args.ID, ", org id: ", args.OrgID, ", transport id: ", args.TransportID, ", err: ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//获取订单数据
	var orderData FieldsOrder
	orderData = getByID(args.ID)
	if orderData.ID < 1 {
		err = errors.New("no data")
	}
	//检查订单是否完成支付？通知配送单
	if orderData.PayStatus == 1 {
		//推送完成支付
		pushNatsOrderPay(orderData.ID, orderData.PayID, orderData.Status >= 2)
	}
	//反馈
	return
}

// ArgsUpdateTransportAuto 修改配送安排时间和自动配送设置参数
type ArgsUpdateTransportAuto struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//调整说明描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//允许自动配送
	TransportAllowAuto bool `db:"transport_allow_auto" json:"transportAllowAuto" check:"bool"`
	//期望送货时间
	TransportTaskAt time.Time `db:"transport_task_at" json:"transportTaskAt" check:"isoTime"`
}

// UpdateTransportAuto 修改配送安排时间和自动配送设置
func UpdateTransportAuto(args *ArgsUpdateTransportAuto) (err error) {
	var newLog string
	newLog, err = getLogData(0, args.OrgBindID, "transport_auto", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), transport_allow_auto = :transport_allow_auto, transport_task_at = :transport_task_at, logs = logs || :log WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":                   args.ID,
		"org_id":               args.OrgID,
		"user_id":              0,
		"log":                  newLog,
		"transport_allow_auto": args.TransportAllowAuto,
		"transport_task_at":    args.TransportTaskAt,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//反馈
	return
}

// ArgsUpdateTransportInfo 修改第三方物流配送信息
type ArgsUpdateTransportInfo struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//配送服务系统
	// 0 self 其他配送; 1 transport 自运营配送; 2 running 跑腿服务; 3 housekeeping 家政服务
	TransportSystem string `db:"transport_system" json:"transportSystem" check:"mark"`
	//配送单号
	TransportSN string `db:"transport_sn" json:"transportSN"`
	//配送服务的状态信息
	TransportInfo string `db:"transport_info" json:"transportInfo"`
	//配送状态
	// 0 等待分配人员; 1 取货中; 2 送货中; 3 完成配送
	TransportStatus int `db:"transport_status" json:"transportStatus" check:"intThan0" empty:"true"`
}

func UpdateTransportInfo(args *ArgsUpdateTransportInfo) (err error) {
	var newLog string
	newLog, err = getLogData(0, args.OrgBindID, "transport_info", "修改订单的配送服务信息")
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), transport_system = :transport_system, transport_sn = :transport_sn, transport_info = :transport_info, transport_status = :transport_status, logs = logs || :log WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":               args.ID,
		"org_id":           args.OrgID,
		"user_id":          0,
		"log":              newLog,
		"transport_system": args.TransportSystem,
		"transport_sn":     args.TransportSN,
		"transport_info":   args.TransportInfo,
		"transport_status": args.TransportStatus,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//如果为完成配送，则检查订单状态并完成订单
	if args.TransportStatus == 3 {
		orderData := getByID(args.ID)
		if orderData.PricePay == true {
			//订单标记完成
			err = UpdateFinish(&ArgsUpdateFinish{
				ID:        args.ID,
				OrgID:     -1,
				UserID:    -1,
				OrgBindID: args.OrgBindID,
				Des:       "第三方配送完成自动完成订单",
			})
			if err != nil {
				err = errors.New(fmt.Sprint("update order tms self, and auto update order finish failed, ", err, ", order id: ", args.ID))
				return
			}
		} else {
			//完成配送标记
			var newLog2 string
			newLog2, err = getLogData(0, args.OrgBindID, "transport_info", "订单配送完成，等待完成订单")
			if err != nil {
				return
			}
			_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), status = 3, logs = logs || :log WHERE id = :id", map[string]interface{}{
				"id":  args.ID,
				"log": newLog2,
			})
			if err != nil {
				return
			}
		}
		//清理缓冲
		deleteOrderCache(args.ID)
	}
	//反馈
	return
}

// argsUpdateTransportFailed 配送失败，等待人工干预参数
type argsUpdateTransportFailed struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//调整说明描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// updateTransportFailed 配送失败，等待人工干预
func updateTransportFailed(args *argsUpdateTransportFailed) (err error) {
	var newLog string
	newLog, err = getLogData(0, args.OrgBindID, "transport_auto", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), transport_allow_auto = false, transport_id = 0, logs = logs || :log WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"user_id": 0,
		"log":     newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//反馈
	return
}

// 检查是否存在服务单
func checkOrderHaveHousekeeping(orderID int64) (b bool) {
	data := getByID(orderID)
	if data.ID < 1 {
		return
	}
	if data.TransportSystem == "housekeeping" {
		if data.TransportID > 0 {
			b = true
			return
		}
	}
	return
}
