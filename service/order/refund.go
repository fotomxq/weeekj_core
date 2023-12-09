package ServiceOrder

import (
	"encoding/json"
	"errors"
	"fmt"
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseExpireTip "gitee.com/weeekj/weeekj_core/v5/base/expire_tip"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	MallCoreMod "gitee.com/weeekj/weeekj_core/v5/mall/core/mod"
	OrgWorkTipMod "gitee.com/weeekj/weeekj_core/v5/org/work_tip/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	TMSTransport "gitee.com/weeekj/weeekj_core/v5/tms/transport"
	UserTicket "gitee.com/weeekj/weeekj_core/v5/user/ticket"
	"github.com/lib/pq"
	"strings"
	"time"
)

// ArgsRefundPost 申请退货参数
type ArgsRefundPost struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//退货原因
	RefundWay string `db:"refund_way" json:"refundWay" check:"des" min:"1" max:"600" empty:"true"`
	//退货备注
	RefundDes string `db:"refund_des" json:"refundDes" check:"des" min:"1" max:"1000" empty:"true"`
	//退货图片列
	RefundFileIDs pq.Int64Array `db:"refund_file_ids" json:"refundFileIDs" check:"ids" empty:"true"`
	//退货是否收到货物
	RefundHaveGood bool `db:"refund_have_good" json:"refundHaveGood" check:"bool"`
}

// RefundPost 申请退货
func RefundPost(args *ArgsRefundPost) (errCode string, err error) {
	//获取订单
	orderData := getByID(args.ID)
	if orderData.ID < 1 || !CoreFilter.EqID2(args.OrgID, orderData.OrgID) || !CoreFilter.EqID2(args.UserID, orderData.UserID) || orderData.RefundStatus != 0 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//订单状态不符合基本退货条件
	if orderData.Status == 0 || orderData.Status == 1 || orderData.Status == 5 || orderData.Status == 6 {
		errCode = "err_refund_status"
		err = errors.New("no data")
		return
	}
	//获取配置
	serviceOrderRefundExpire := BaseConfig.GetDataStringNoErr("ServiceOrderRefundExpire")
	if serviceOrderRefundExpire == "" {
		serviceOrderRefundExpire = "168h"
	}
	var serviceOrderRefundExpireAt time.Time
	serviceOrderRefundExpireAt, err = CoreFilter.GetTimeByAdd(serviceOrderRefundExpire)
	if err != nil || serviceOrderRefundExpireAt.Unix() <= CoreFilter.GetNowTime().Unix()+3600 {
		serviceOrderRefundExpireAt = CoreFilter.GetNowTimeCarbon().AddHour().Time
	}
	//修正参数
	if args.RefundFileIDs == nil {
		args.RefundFileIDs = pq.Int64Array{}
	} else {
		if len(args.RefundFileIDs) > 7 {
			errCode = "err_file_max"
			err = errors.New("file too many")
			return
		}
	}
	//生成日志
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "refund_post", args.RefundWay)
	if err != nil {
		errCode = "err_log"
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), refund_status = 1, refund_way = :refund_way, refund_des = :refund_des, refund_file_ids = :refund_file_ids, refund_have_good = :refund_have_good, refund_expire_at = :refund_expire_at, logs = logs || :log WHERE id = :id", map[string]interface{}{
		"id":               args.ID,
		"refund_way":       args.RefundWay,
		"refund_des":       args.RefundDes,
		"refund_file_ids":  args.RefundFileIDs,
		"refund_have_good": args.RefundHaveGood,
		"refund_expire_at": serviceOrderRefundExpireAt,
		"log":              newLog,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//重新获取订单详情
	orderData = getByID(args.ID)
	//统计
	orderSystemMarkKey := getOrderSystemMarkKey(orderData.SystemMark)
	AnalysisAny2.AppendData("add", "service_order_refund_create_count", time.Time{}, orderData.OrgID, orderData.UserID, 0, orderSystemMarkKey, 0, 1)
	//退货过期处理
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      orderData.OrgID,
		UserID:     0,
		SystemMark: "service_order_refund",
		BindID:     orderData.ID,
		Hash:       "",
		ExpireAt:   orderData.RefundExpireAt,
	})
	//反馈
	return
}

// ArgsRefundAudit 审核退货进入退货中参数
type ArgsRefundAudit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
	//是否需要配送?
	NeedTransport bool `db:"need_transport" json:"needTransport" check:"bool"`
	//同时申请退款
	NeedRefundPay bool `json:"needRefundPay" check:"bool"`
	//退款金额
	RefundPrice int64 `json:"refundPrice" check:"int64Than0"`
}

// RefundAudit 审核退货进入退货中
// 同时将根据订单信息，发起退单配送单处理
func RefundAudit(args *ArgsRefundAudit) (errCode string, err error) {
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		errCode = "order_not_exist"
		return
	}
	if orderData.RefundStatus != 1 {
		errCode = "order_have_refund"
		err = errors.New("order have refund")
		return
	}
	if args.NeedRefundPay {
		if orderData.PayStatus != 1 {
			errCode = "order_no_pay"
			err = errors.New("order no have pay")
			return
		}
		if args.RefundPrice > orderData.Price {
			errCode = "too_much_refund_price"
			err = errors.New("too much refund price")
			return
		}
	}
	//退货日志
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "refund_audit", args.Des)
	if err != nil {
		errCode = "add_log"
		return
	}
	//如果是自运营配送
	if orderData.TransportSystem == "transport" && args.NeedTransport {
		//创建配送单
		var transportData TMSTransport.FieldsTransport
		var tmsGoods []TMSTransport.FieldsTransportGood
		for _, v := range orderData.Goods {
			tmsGoods = append(tmsGoods, TMSTransport.FieldsTransportGood{
				System: v.From.System,
				ID:     v.From.ID,
				Mark:   "",
				Name:   v.From.Name,
				Count:  int(v.Count),
			})
		}
		transportData, _, err = TMSTransport.CreateTransport(&TMSTransport.ArgsCreateTransport{
			OrgID:       orderData.ID,
			BindID:      0,
			InfoID:      0,
			UserID:      orderData.UserID,
			FromAddress: orderData.AddressTo,
			ToAddress:   orderData.AddressFrom,
			OrderID:     orderData.ID,
			Goods:       tmsGoods,
			Weight:      0,
			Length:      0,
			Width:       0,
			Currency:    86,
			Price:       0,
			PayFinish:   true,
			Params: CoreSQLConfig.FieldsConfigsType{
				{
					Mark: "isRefund",
					Val:  "true",
				},
			},
		})
		if err != nil {
			errCode = "create_tms"
			return
		}
		//更新信息
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET refund_transport_system = :refund_transport_system, refund_transport_sn = :refund_transport_sn, refund_transport_info = :refund_transport_info, transport_id = :transport_id, transport_ids = array_append(transport_ids, :transport_id) WHERE id = :id", map[string]interface{}{
			"id":                      args.ID,
			"refund_transport_system": "transport",
			"refund_transport_sn":     transportData.SN,
			"refund_transport_info":   "",
			"transport_id":            transportData.ID,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
	}
	//修改订单信息
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), refund_status = 2, logs = logs || :log WHERE id = :id", map[string]interface{}{
		"id":  args.ID,
		"log": newLog,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//发起退款请求
	if args.NeedRefundPay {
		errCode, err = RefundPay(&ArgsRefundPay{
			ID:          args.ID,
			OrgID:       args.OrgID,
			UserID:      args.UserID,
			OrgBindID:   args.OrgBindID,
			RefundPrice: args.RefundPrice,
			Des:         args.Des,
		})
		if err != nil {
			return
		}
	}
	//反馈
	return
}

// ArgsRefundUpdateTMS 更新退货的配送信息参数
type ArgsRefundUpdateTMS struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//退货快递类型
	// 0 self 其他配送; 1 take 自提; 2 transport 自运营配送; 3 running 跑腿服务; 4 housekeeping 家政服务
	RefundTransportSystem string `db:"refund_transport_system" json:"refundTransportSystem"`
	//退货快递单号
	RefundTransportSN string `db:"refund_transport_sn" json:"refundTransportSN"`
	//配送服务的状态信息
	RefundTransportInfo string `db:"refund_transport_info" json:"refundTransportInfo"`
}

// RefundUpdateTMS 更新退货的配送信息
// 只有其他配送，用户可以修改
func RefundUpdateTMS(args *ArgsRefundUpdateTMS) (errCode string, err error) {
	//获取订单
	orderData := getByID(args.ID)
	if orderData.ID < 1 || !CoreFilter.EqID2(args.OrgID, orderData.OrgID) || !CoreFilter.EqID2(args.UserID, orderData.UserID) || orderData.RefundStatus != 2 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//用户只能修改第三方配送单
	if orderData.RefundTransportSystem != "" {
		if orderData.RefundTransportSystem != "self" && args.UserID > 0 {
			errCode = "err_refund_other_tms"
			err = errors.New("order refund tms not self")
			return
		}
	}
	//日志
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "refund_tms", "更新退货配送单信息")
	if err != nil {
		errCode = "err_update"
		return
	}
	//更新信息
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), refund_transport_system = :refund_transport_system, refund_transport_sn = :refund_transport_sn, refund_transport_info = :refund_transport_info, logs = logs || :log WHERE id = :id", map[string]interface{}{
		"id":                      args.ID,
		"refund_transport_system": args.RefundTransportSystem,
		"refund_transport_sn":     args.RefundTransportSN,
		"refund_transport_info":   args.RefundTransportInfo,
		"log":                     newLog,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

// ArgsRefundTip 提醒商家尽快处理退货请求参数
type ArgsRefundTip struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// RefundTip 提醒商家尽快处理退货请求
func RefundTip(args *ArgsRefundTip) (errCode string, err error) {
	//获取订单
	orderData := getByID(args.ID)
	if orderData.ID < 1 || !CoreFilter.EqID2(args.UserID, orderData.UserID) || orderData.RefundStatus != 1 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//更新信息
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), refund_tip_at = NOW() WHERE id = :id", map[string]interface{}{
		"id": args.ID,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//发送组织工作提醒
	OrgWorkTipMod.AppendTip(&OrgWorkTipMod.ArgsAppendTip{
		OrgID:     orderData.OrgID,
		OrgBindID: 0,
		Msg:       "客户催促办理订单退货申请",
		System:    "service_order",
		BindID:    orderData.ID,
	})
	//反馈
	return
}

// ArgsRefundPay 发起退款请求处理参数
type ArgsRefundPay struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//退款金额
	RefundPrice int64 `json:"refundPrice" check:"price"`
	//调整说明描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// RefundPay 发起退款请求处理
func RefundPay(args *ArgsRefundPay) (errCode string, err error) {
	//获取订单
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		errCode = "order_not_exist"
		err = errors.New(fmt.Sprint("get order data, ", err))
		return
	}
	//检查是否已经退款
	refundPay, b := orderData.Params.GetValBool("refundPay")
	if b && refundPay {
		errCode = "repeat_refund_pay"
		err = errors.New("have refund pay")
		return
	}
	//记录日志
	var appendLog []FieldsLog
	//检查订单是否已经支付
	var payData FinancePay.FieldsPayType
	//遍历商品
	isHaveVirtual := false
	virtualCanRefund := true
	for _, v := range orderData.Goods {
		//退还虚拟商品
		// 检查是否为虚拟商品
		if MallCoreMod.CheckProductIsVirtual(v.From.ID, orderData.OrgID) {
			isHaveVirtual = true
			continue
		}
		if isHaveVirtual && v.From.System == "mall" {
			//获取商品
			var productData MallCoreMod.FieldsCore
			productData, err = MallCoreMod.GetProduct(&MallCoreMod.ArgsGetProduct{
				ID:    v.From.ID,
				OrgID: orderData.OrgID,
			})
			if err != nil {
				err = nil
				continue
			}
			//检查三天内，是否存在使用情况
			var oldOrderList []FieldsOrder
			err = Router2SystemConfig.MainDB.Select(&oldOrderList, "SELECT id, goods FROM service_order WHERE create_at >= $1 AND create_at >= $2 AND user_id = $3 AND org_id = $4 AND delete_at < to_timestamp(1000000)", CoreFilter.GetNowTimeCarbon().SubDays(3).Time, orderData.CreateAt, orderData.UserID, orderData.OrgID)
			if err == nil && len(oldOrderList) > 0 {
				//遍历旧的订单，检查是否存在使用该票据情况
				for _, vProductGivingTicket := range productData.GivingTickets {
					for _, vOld := range oldOrderList {
						for _, vGood := range vOld.Goods {
							for _, vEx := range vGood.Exemptions {
								if vEx.System == "user_ticket" && vEx.ConfigID == vProductGivingTicket.TicketConfigID {
									virtualCanRefund = false
									break
								}
							}
							if !virtualCanRefund {
								break
							}
						}
						if !virtualCanRefund {
							break
						}
					}
					if !virtualCanRefund {
						break
					}
				}
			}
		}
		if !virtualCanRefund {
			break
		}
	}
	if !virtualCanRefund {
		errCode = "user_ticket_is_used"
		err = errors.New("order refund user ticket, buy user ticket have use")
		return
	}
	if isHaveVirtual {
		//获取订单扩展参数
		mallVirtualUserTicketsStr, b := orderData.Params.GetVal("mall_virtual_user_tickets_can_refund")
		if b {
			//解析数据
			var mallVirtualUserTicketsIDs []int64
			mallVirtualUserTicketsIDsStrs := strings.Split(mallVirtualUserTicketsStr, ",")
			for _, vS := range mallVirtualUserTicketsIDsStrs {
				var vSID int64
				vSID, err = CoreFilter.GetInt64ByString(vS)
				if err != nil {
					CoreLog.Error("order params is error, mall_virtual_user_tickets have error format, order id: ", orderData.ID)
					continue
				}
				mallVirtualUserTicketsIDs = append(mallVirtualUserTicketsIDs, vSID)
			}
			if len(mallVirtualUserTicketsIDs) > 0 {
				//发起退票请求
				err = UserTicket.RefundUseTicket(&UserTicket.ArgsRefundUseTicket{
					OrgID:  orderData.OrgID,
					UserID: orderData.UserID,
					IDs:    mallVirtualUserTicketsIDs,
					Des:    fmt.Sprint("订单ID[", orderData.ID, "]SN[", orderData.SerialNumber, "]SNDay[", orderData.SerialNumberDay, "]发生退款退货"),
				})
				if err != nil {
					errCode = "refund_use_user_ticket"
					err = errors.New(fmt.Sprint("refund use user ticket, ", err))
					return
				}
			}
		}
	}
	//记录退款金额
	var refundPayID int64 = 0
	var refundPrice int64 = 0
	//检查金额
	if orderData.PayID > 0 && orderData.Price > 0 && orderData.PricePay {
		//获取支付请求
		payData, err = FinancePay.GetOne(&FinancePay.ArgsGetOne{
			ID:  orderData.PayID,
			Key: "",
		})
		if err != nil {
			errCode = "pay_not_exist"
			err = errors.New(fmt.Sprint("get pay data, ", err))
			return
		}
		//检查已经已经支付，或发生了退款
		if payData.Status != 3 {
			errCode = "pay_not_finish"
			err = errors.New("pay status not finish")
			return
		}
		//检查退款金额
		if args.RefundPrice == -1 {
			args.RefundPrice = orderData.Price
		}
		if args.RefundPrice > orderData.Price {
			errCode = "price_too_many"
			err = errors.New(fmt.Sprint("refund price more than order price, args price: ", args.RefundPrice, ", data price: ", orderData.Price))
			return
		}
		//发起退单处理
		err = FinancePay.CheckTakeFrom(&FinancePay.ArgsCheckTakeFrom{
			ID: orderData.PayID,
			TakeFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     orderData.OrgID,
				Mark:   "",
				Name:   "",
			},
		})
		if err != nil {
			errCode = "no_operate_pay"
			err = errors.New(fmt.Sprint("check finance pay take from, ", err))
			return
		}
		errCode, err = FinancePay.UpdateStatusRefund(&FinancePay.ArgsUpdateStatusRefund{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "org_bind",
				ID:     args.OrgBindID,
				Mark:   "",
				Name:   "",
			},
			ID:          orderData.PayID,
			Key:         "",
			Params:      CoreSQLConfig.FieldsConfigsType{},
			RefundPrice: 0,
			Des:         args.Des,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update status refund, pay id: ", orderData.PayID, ", refund price: ", args.RefundPrice, ", err: ", err))
			return
		}
		//自动审核退款请求
		errCode, err = FinancePay.UpdateStatusRefundAudit(&FinancePay.ArgsUpdateStatusRefundAudit{
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "org_bind",
				ID:     args.OrgBindID,
				Mark:   "",
				Name:   "",
			},
			ID:          orderData.PayID,
			Key:         "",
			Params:      nil,
			RefundPrice: args.RefundPrice,
			Des:         args.Des,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update finance pay status to refund, ", err))
			return
		}
		//记录关键信息
		refundPayID = orderData.PayID
		refundPrice = args.RefundPrice
		//记录日志
		appendLog = append(appendLog, FieldsLog{
			CreateAt:  CoreFilter.GetNowTime(),
			UserID:    args.UserID,
			OrgBindID: args.OrgBindID,
			Mark:      "refund_finish",
			Des:       fmt.Sprint("向财务中心发起退还用户缴纳的费用 ¥", float64(args.RefundPrice)/100, "，退款渠道[", payData.PaymentChannel.System, "~", payData.PaymentChannel.Mark, "]"),
		})
	}
	//开始检查票据是否存在
	for _, v := range orderData.Exemptions {
		if v.System != "user_ticket" {
			continue
		}
		var vTicketConfig UserTicket.FieldsConfig
		vTicketConfig, err = UserTicket.GetConfigByID(&UserTicket.ArgsGetConfigByID{
			ID:    v.ConfigID,
			OrgID: 0,
		})
		if err != nil {
			err = nil
			continue
		}
		canRefund, b := vTicketConfig.Params.GetValBool("canRefund")
		if !b || !canRefund {
			continue
		}
		err = UserTicket.AddTicket(&UserTicket.ArgsAddTicket{
			OrgID:       orderData.OrgID,
			ConfigID:    vTicketConfig.ID,
			UserID:      orderData.UserID,
			Count:       v.Count,
			UseFromName: "退还",
		})
		if err != nil {
			CoreLog.Error("update order cancel, ticket refund, order id: ", orderData.ID, ", ticket id: ", v.ConfigID, ", count: ", v.Count, ", err: ", err)
			err = nil
			continue
		}
		//记录日志
		appendLog = append(appendLog, FieldsLog{
			CreateAt:  CoreFilter.GetNowTime(),
			UserID:    args.UserID,
			OrgBindID: args.OrgBindID,
			Mark:      "refund_finish",
			Des:       fmt.Sprint("向用户票据发起退还用户的票据[", vTicketConfig.Title, "][", v.Count, "]张"),
		})
	}
	//遍历商品
	for _, v := range orderData.Goods {
		//退票用户票据数据
		for _, v2 := range v.Exemptions {
			if v2.System != "user_ticket" {
				continue
			}
			var vTicketConfig UserTicket.FieldsConfig
			vTicketConfig, err = UserTicket.GetConfigByID(&UserTicket.ArgsGetConfigByID{
				ID:    v2.ConfigID,
				OrgID: 0,
			})
			if err != nil {
				err = nil
				continue
			}
			canRefund, b := vTicketConfig.Params.GetValBool("canRefund")
			if !b || !canRefund {
				continue
			}
			err = UserTicket.AddTicket(&UserTicket.ArgsAddTicket{
				OrgID:       orderData.OrgID,
				ConfigID:    vTicketConfig.ID,
				UserID:      orderData.UserID,
				Count:       v2.Count,
				UseFromName: "退还",
			})
			if err != nil {
				CoreLog.Error("update order cancel, ticket refund, order id: ", orderData.ID, ", ticket id: ", v2.ConfigID, ", count: ", v2.Count, ", err: ", err)
				err = nil
				continue
			}
			//记录日志
			appendLog = append(appendLog, FieldsLog{
				CreateAt:  CoreFilter.GetNowTime(),
				UserID:    args.UserID,
				OrgBindID: args.OrgBindID,
				Mark:      "refund_finish",
				Des:       fmt.Sprint("向用户票据发起退还用户的票据[", vTicketConfig.Title, "][", v2.Count, "]张"),
			})
		}
	}
	//记录退款
	orderData.Params = append(orderData.Params, CoreSQLConfig.FieldsConfigType{
		Mark: "refundPay",
		Val:  "true",
	})
	//申请完成后，更新订单状态
	appendLog = append(appendLog, FieldsLog{
		CreateAt:  CoreFilter.GetNowTime(),
		UserID:    args.UserID,
		OrgBindID: args.OrgBindID,
		Mark:      "refund_finish",
		Des:       args.Des,
	})
	var newLogByte []byte
	newLogByte, err = json.Marshal(appendLog)
	if err != nil {
		errCode = "order_log"
		err = errors.New(fmt.Sprint("add order log data, ", err))
		return
	}
	newLog := string(newLogByte)
	if newLog == "" {
		errCode = "order_log"
		err = errors.New(fmt.Sprint("add order log data, ", err))
		return
	}
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), pay_status = 2, refund_pay_id = :refund_pay_id, refund_price = :refund_price, refund_pay_finish = :refund_pay_finish, logs = logs || :log, params = :params WHERE id = :id", map[string]interface{}{
		"id":                args.ID,
		"refund_pay_id":     refundPayID,
		"refund_price":      refundPrice,
		"refund_pay_finish": CoreFilter.GetNowTime(),
		"log":               newLog,
		"params":            orderData.Params,
	}); err != nil {
		errCode = "order_update"
		err = errors.New(fmt.Sprint("update order data, ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

// ArgsRefundCancel 取消退货参数
type ArgsRefundCancel struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}

// RefundCancel 取消退货
func RefundCancel(args *ArgsRefundCancel) (err error) {
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "refund_cancel", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), refund_status = 0, logs = logs || :log WHERE id = :id AND (refund_status = 1 OR refund_status = 2) AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"user_id": args.UserID,
		"log":     newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//获取订单的配送单
	orderData := getByID(args.ID)
	if orderData.TransportID > 0 {
		var tmsData TMSTransport.FieldsTransport
		tmsData, err = TMSTransport.GetTransport(&TMSTransport.ArgsGetTransport{
			ID:     orderData.TransportID,
			OrgID:  args.OrgID,
			InfoID: -1,
			UserID: args.UserID,
		})
		if err != nil {
			return
		}
		if tmsData.Status != 3 {
			err = TMSTransport.DeleteTransport(&TMSTransport.ArgsDeleteTransport{
				ID:     tmsData.ID,
				OrgID:  0,
				BindID: 0,
			})
			if err != nil {
				return
			}
		}
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

// ArgsRefundFailed 支付失败参数
type ArgsRefundFailed struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}

// RefundFailed 支付失败
func RefundFailed(args *ArgsRefundFailed) (err error) {
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "refund_failed", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), refund_status = 4, logs = logs || :log WHERE id = :id AND pay_status != 5 AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"user_id": args.UserID,
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

// ArgsRefundFinish 完成退货参数
type ArgsRefundFinish struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	// 可选，作为验证
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}

// RefundFinish 完成退货
func RefundFinish(args *ArgsRefundFinish) (err error) {
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "refund_finish", args.Des)
	if err != nil {
		err = errors.New(fmt.Sprint("add order log, ", err))
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), status = 6, refund_status = 3, logs = logs || :log WHERE id = :id AND refund_status != 3 AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"user_id": args.UserID,
		"log":     newLog,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update order data, ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//订单收尾工作
	orderCancelLast(args.ID, args.OrgBindID)
	//通知取消订单
	CoreNats.PushDataNoErr("/service/order/update", "refund", args.ID, "", nil)
	//反馈
	return
}
