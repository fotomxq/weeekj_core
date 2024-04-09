package ServiceOrder

import (
	"errors"
	"fmt"
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	MallCoreMod "github.com/fotomxq/weeekj_core/v5/mall/core/mod"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
	UserSubscription "github.com/fotomxq/weeekj_core/v5/user/subscription"
	UserTicket "github.com/fotomxq/weeekj_core/v5/user/ticket"
	"github.com/lib/pq"
	"time"
)

// 创建订单
// 内部方法，该方法用于run维护中，自动检索列队数据并创建内容
// 该方法不对货物来源进行核对
type argsCreate struct {
	//来源系统
	// 该订单创建来源的系统
	// eg: user_sub / org_sub / mall/ user_integral
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//收取货物地址
	AddressFrom CoreSQLAddress.FieldsAddress `db:"address_from" json:"addressFrom"`
	//送货地址
	AddressTo CoreSQLAddress.FieldsAddress `db:"address_to" json:"addressTo"`
	//货物清单
	Goods FieldsGoods `db:"goods" json:"goods"`
	//订单总的抵扣
	// 例如满减活动，不局限于个别商品的活动
	Exemptions ServiceOrderWaitFields.FieldsExemptions `db:"exemptions" json:"exemptions"`
	//是否允许自动审核
	// 客户提交订单后，将自动审核该订单。订单如果存在至少一件未开启的商品，将禁止该操作
	AllowAutoAudit bool `db:"allow_auto_audit" json:"allowAutoAudit" check:"bool"`
	//允许自动配送
	TransportAllowAuto bool `db:"transport_allow_auto" json:"transportAllowAuto" check:"bool"`
	//是否允许货到付款？
	TransportPayAfter bool `db:"transport_pay_after" json:"transportPayAfter" check:"bool"`
	//期望送货时间
	TransportTaskAt time.Time `db:"transport_task_at" json:"transportTaskAt" check:"isoTime" empty:"true"`
	//配送服务系统
	// 0 self 自运营服务; 1 自提; 2 running 跑腿服务; 3 housekeeping 家政服务
	TransportSystem string `db:"transport_system" json:"transportSystem"`
	//费用组成
	PriceList FieldsPrices `db:"price_list" json:"priceList"`
	//订单总费用
	// 货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	// 总费用金额
	Price int64 `db:"price" json:"price" check:"price"`
	//折扣前费用
	PriceTotal int64 `db:"price_total" json:"priceTotal" check:"price"`
	//备注信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//日志
	Logs FieldsLogs `db:"logs" json:"logs"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func create(args *argsCreate) (data FieldsOrder, errCode, errMsg string, err error) {
	//初始化错误
	errCode = "unknow"
	errMsg = "未知错误"
	//生成SN
	var sn int64 = 1
	var snDay int64 = 1
	sn, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "service_order", "id", "org_id = $1", args.OrgID)
	if err != nil {
		sn = 1
		err = nil
	}
	if sn < 1 {
		sn = 1
	} else {
		sn += 1
	}
	snDay, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "service_order", "id", "org_id = $1 AND create_at >= $2", args.OrgID, CoreFilter.GetNowTimeCarbon().SetHour(0).SetMinute(0).SetSecond(0).Time)
	if err != nil {
		snDay = 1
		err = nil
	}
	if snDay < 1 {
		snDay = 1
	} else {
		snDay += 1
	}
	//追加日志
	if len(args.Logs) < 1 {
		args.Logs = append(args.Logs, FieldsLog{
			CreateAt:  CoreFilter.GetNowTime(),
			UserID:    args.UserID,
			OrgBindID: 0,
			Mark:      "create",
			Des:       "创建新的订单",
		})
	}
	//构建状态
	payStatus := 0
	if args.Price < 1 {
		args.Price = 0
		payStatus = 1
	}
	//占用资源
	var useUserSubID int64 = 0
	var paySystem string
	for _, v := range args.Goods {
		for _, v2 := range v.Exemptions {
			if v2.System == "user_sub" {
				useUserSubID = v2.ConfigID
			}
			if v2.System == "user_ticket" && v2.Count > 0 {
				if err = UserTicket.UseTicket(&UserTicket.ArgsUseTicket{
					ID:          0,
					OrgID:       args.OrgID,
					ConfigID:    v2.ConfigID,
					UserID:      args.UserID,
					Count:       v2.Count,
					UseFromName: "订单",
				}); err != nil {
					errCode = "user_ticket"
					errMsg = "票据限制或用户票据不足，无法使用用户票据"
					err = errors.New(fmt.Sprint("use user ticket failed, ", err))
					return
				}
				if args.Price < 1 {
					paySystem = "user_ticket"
				}
			}
		}
	}
	for _, v := range args.Exemptions {
		if v.System == "user_sub" {
			useUserSubID = v.ConfigID
		}
		if v.System == "user_ticket" && v.Count > 0 {
			if err = UserTicket.UseTicket(&UserTicket.ArgsUseTicket{
				ID:          0,
				OrgID:       args.OrgID,
				ConfigID:    v.ConfigID,
				UserID:      args.UserID,
				Count:       v.Count,
				UseFromName: "订单",
			}); err != nil {
				errCode = "user_ticket"
				errMsg = "票据限制或用户票据不足，无法使用用户票据"
				err = errors.New(fmt.Sprint("use user ticket failed, ", err))
				return
			}
			if args.Price < 1 {
				paySystem = "user_ticket"
			}
		}
	}
	if useUserSubID > 0 {
		//使用会员
		if err = UserSubscription.UseSub(&UserSubscription.ArgsUseSub{
			ConfigID:    useUserSubID,
			UserID:      args.UserID,
			UseFrom:     "order",
			UseFromName: "订单",
		}); err != nil {
			errCode = "user_sub"
			errMsg = "会员限制或其他原因，用户无法使用会员"
			err = errors.New(fmt.Sprint("use user sub failed, ", err))
			return
		}
	}
	//写入支付方式
	if paySystem != "" {
		args.Params = CoreSQLConfig.Set(args.Params, "paySystem", paySystem)
	}
	//生成数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_order", "INSERT INTO service_order (expire_at, system_mark, org_id, user_id, create_from, serial_number, serial_number_day, status, refund_status, address_from, address_to, goods, exemptions, allow_auto_audit, transport_id, transport_allow_auto, transport_task_at, transport_pay_after, transport_ids, transport_system, transport_sn, transport_info, transport_status, price_list, price_pay, currency, price, price_total, pay_status, pay_id, pay_list, des, logs, params) VALUES (:expire_at, :system_mark, :org_id, :user_id, :create_from, :serial_number, :serial_number_day, 0, 0, :address_from, :address_to, :goods, :exemptions, :allow_auto_audit, 0, :transport_allow_auto, :transport_task_at, :transport_pay_after, :transport_ids, :transport_system, '', '', 0, :price_list, :price < 1, :currency, :price, :price_total, :pay_status, 0, :pay_list, :des, :logs, :params)", map[string]interface{}{
		"expire_at":            CoreFilter.GetNowTimeCarbon().AddHours(24),
		"system_mark":          args.SystemMark,
		"org_id":               args.OrgID,
		"user_id":              args.UserID,
		"create_from":          args.CreateFrom,
		"serial_number":        sn,
		"serial_number_day":    snDay,
		"address_from":         args.AddressFrom,
		"address_to":           args.AddressTo,
		"goods":                args.Goods,
		"exemptions":           args.Exemptions,
		"allow_auto_audit":     args.AllowAutoAudit,
		"transport_allow_auto": args.TransportAllowAuto,
		"transport_task_at":    args.TransportTaskAt,
		"transport_pay_after":  args.TransportPayAfter,
		"transport_ids":        pq.Int64Array{},
		"transport_system":     args.TransportSystem,
		"price_list":           args.PriceList,
		"currency":             args.Currency,
		"price":                args.Price,
		"price_total":          args.PriceTotal,
		"pay_status":           payStatus,
		"pay_list":             pq.Int64Array{},
		"des":                  args.Des,
		"logs":                 args.Logs,
		"params":               args.Params,
	}, &data)
	if err != nil {
		err = errors.New(fmt.Sprint("insert order, ", err))
		errCode = "insert"
		errMsg = "创建数据失败，数据库异常"
		return
	}
	//检查订单商品的类型
	goodType := ""
	for _, v := range data.Goods {
		switch v.From.Mark {
		case "virtual":
			//虚拟商品标记
			if MallCoreMod.CheckProductIsVirtual(v.From.ID, data.OrgID) {
				if goodType == "" {
					goodType = "virtual"
				} else {
					goodType = "mixed"
				}
				continue
			}
		case "housekeeping":
			//家政商品标记
			if goodType == "" {
				goodType = "housekeeping"
			} else {
				goodType = "mixed"
			}
			continue
		default:
			if goodType == "" {
				goodType = "mall"
			} else {
				goodType = "mixed"
			}
			continue
		}
	}
	//发出列队请求
	CoreNats.PushDataNoErr("service_order_create", "/service/order/create", data.SystemMark, data.ID, goodType, data)
	//发送过期提醒模块
	if err = BaseExpireTip.AppendTip(&BaseExpireTip.ArgsAppendTip{
		OrgID:      data.OrgID,
		UserID:     data.UserID,
		SystemMark: "service_order",
		BindID:     data.ID,
		Hash:       "",
		ExpireAt:   data.ExpireAt,
	}); err != nil {
		CoreLog.Error("service order create order, append base expire tip failed, ", err)
		err = nil
	}
	//更新组织用户数据
	if data.OrgID > 0 && data.UserID > 0 {
		OrgUserMod.PushUpdateUserData(data.OrgID, data.UserID)
	}
	//反馈
	return
}

// 根据等待订单创建订单
func createByWait(waitData *ServiceOrderWaitFields.FieldsWait) error {
	//锁定
	createWaitLock.Lock()
	defer createWaitLock.Unlock()
	//检查订单是否被创建了
	if checkCreateWait(waitData.ID) {
		return nil
	}
	//构建订单结构
	var goods FieldsGoods
	for _, v2 := range waitData.Goods {
		var exemptions FieldsExemptions
		for _, v3 := range v2.Exemptions {
			exemptions = append(exemptions, FieldsExemption{
				System:   v3.System,
				ConfigID: v3.ConfigID,
				Name:     v3.Name,
				Des:      v3.Des,
				Count:    v3.Count,
				Price:    v3.Price,
			})
		}
		goods = append(goods, FieldsGood{
			From:            v2.From,
			OptionKey:       v2.OptionKey,
			Count:           v2.Count,
			Price:           v2.Price,
			Exemptions:      exemptions,
			CommentBuyer:    false,
			CommentBuyerID:  0,
			CommentSeller:   false,
			CommentSellerID: 0,
		})
	}
	var logs FieldsLogs
	for _, v2 := range waitData.Logs {
		logs = append(logs, FieldsLog{
			CreateAt:  v2.CreateAt,
			UserID:    v2.UserID,
			OrgBindID: v2.OrgBindID,
			Mark:      v2.Mark,
			Des:       v2.Des,
		})
	}
	priceList := FieldsPrices{}
	for _, v2 := range waitData.PriceList {
		priceList = append(priceList, FieldsPrice{
			PriceType: v2.PriceType,
			IsPay:     v2.IsPay,
			Price:     v2.Price,
		})
	}
	//创建订单
	newOrderData, errCode, errMsg, err := create(&argsCreate{
		SystemMark:         waitData.SystemMark,
		OrgID:              waitData.OrgID,
		UserID:             waitData.UserID,
		CreateFrom:         waitData.CreateFrom,
		AddressFrom:        waitData.AddressFrom,
		AddressTo:          waitData.AddressTo,
		Goods:              goods,
		Exemptions:         waitData.Exemptions,
		AllowAutoAudit:     waitData.AllowAutoAudit,
		TransportAllowAuto: waitData.TransportAllowAuto,
		TransportPayAfter:  waitData.TransportPayAfter,
		TransportTaskAt:    waitData.TransportTaskAt,
		TransportSystem:    waitData.TransportSystem,
		PriceList:          priceList,
		Currency:           waitData.Currency,
		Price:              waitData.Price,
		PriceTotal:         waitData.PriceTotal,
		Des:                waitData.Des,
		Logs:               logs,
		Params:             waitData.Params,
	})
	//如果失败，则标记订单等待失败
	if err != nil {
		_, err2 := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order_wait SET err_code = :err_code, err_msg = :err_msg WHERE id = :id AND order_id < 1", map[string]interface{}{
			"id":       waitData.ID,
			"err_code": errCode,
			"err_msg":  errMsg,
		})
		if err2 != nil {
			err = errors.New(fmt.Sprint("update order wait err, ", err, ", err2: ", err2, ", wait data: ", waitData))
		} else {
			err = errors.New(fmt.Sprint("update order wait err: ", err))
		}
		//清理缓冲
		deleteWaitCache(waitData.ID)
		//反馈失败
		return err
	}
	//更新记录的订单ID
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order_wait SET order_id = :order_id WHERE id = :id AND order_id < 1", map[string]interface{}{
		"id":       waitData.ID,
		"order_id": newOrderData.ID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update wait order id by wait id: ", waitData.ID, ", ", err))
		//清理缓冲
		deleteWaitCache(waitData.ID)
		//反馈失败
		return err
	}
	//清理缓冲
	deleteWaitCache(waitData.ID)
	//重新装在缓冲
	_, _ = getCreateWait(waitData.ID)
	//通知等待订单创建完成
	CoreNats.PushDataNoErr("service_order_create_wait_finish", "/service/order/create_wait_finish", "", waitData.ID, "", map[string]interface{}{
		"orderID": newOrderData.ID,
	})
	//统计
	orderSystemMarkKey := getOrderSystemMarkKey(newOrderData.SystemMark)
	AnalysisAny2.AppendData("add", "service_order_create_count", time.Time{}, newOrderData.OrgID, newOrderData.UserID, 0, orderSystemMarkKey, 0, 1)
	//反馈
	return nil
}

// 检查订单创建过
func checkCreateWait(id int64) (b bool) {
	data, err := getCreateWait(id)
	if err != nil {
		return
	}
	if data.OrderID > 0 {
		b = true
		return
	}
	return
}

// 获取等待订单
func getCreateWait(id int64) (data ServiceOrderWaitFields.FieldsWait, err error) {
	cacheMark := fmt.Sprint("service:order:wait:id:", id)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, order_id, system_mark, org_id, user_id, create_from, hash, address_from, address_to, goods, exemptions, allow_auto_audit, transport_allow_auto, transport_task_at, transport_pay_after, price_list, price_pay, currency, price, price_total, des, logs, params, err_code, err_msg, transport_system FROM service_order_wait WHERE id = $1", id)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 3600)
	return
}

// 删除等待订单ID
func deleteWaitCache(id int64) {
	cacheMark := fmt.Sprint("service:order:wait:id:", id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
