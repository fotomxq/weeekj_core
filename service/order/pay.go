package ServiceOrder

import (
	"errors"
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	FinancePayCreate "github.com/fotomxq/weeekj_core/v5/finance/pay_create"
	FinancePhysicalPay "github.com/fotomxq/weeekj_core/v5/finance/physical_pay"
	FinanceTakeCut "github.com/fotomxq/weeekj_core/v5/finance/take_cut"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsCreatePay 发起支付请求参数
type ArgsCreatePay struct {
	//订单ID组
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
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
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	PaymentChannel CoreSQLFrom.FieldsFrom `db:"payment_channel" json:"paymentChannel"`
	//支付备注
	// 用户环节可根据实际业务需求开放此项
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// CreatePay 发起支付请求
// 允许合并多个订单支付
// 1\不能包含已经支付完成的订单请求
// 2\不能是多种渠道，必须是一个渠道下单
func CreatePay(args *ArgsCreatePay) (payData FinancePay.FieldsPayType, errCode string, err error) {
	//订单数量不能超出10个
	if len(args.IDs) < 1 || len(args.IDs) > 10 {
		errCode = "order_num_limit"
		err = errors.New("too many ids")
		return
	}
	//遍历获取订单待支付金额
	var expireAt time.Time
	createFrom := -1
	currency := -1
	var priceTotal int64 = 0
	var rawOrderList []FieldsOrder
	rawOrderList, err = getByIDs(args.IDs)
	if err != nil || len(rawOrderList) < 1 {
		errCode = "order_not_exist"
		err = errors.New(fmt.Sprint("order not exist, ids: ", args.IDs, ", err: ", err))
		return
	}
	//剔除订单
	var orderList []FieldsOrder
	for _, v := range rawOrderList {
		if v.DeleteAt.Unix() > 1000000 || v.PricePay || v.PayStatus != 0 {
			continue
		}
		orderList = append(orderList, v)
	}
	//支付用户
	var payUserID int64 = 0
	var payOrgID int64 = 0
	for _, orderData := range orderList {
		if orderData.ID < 1 {
			continue
		}
		if orderData.PricePay {
			errCode = "order_is_pay"
			err = errors.New("have order is pay")
			return
		}
		if createFrom < 0 {
			createFrom = orderData.CreateFrom
			currency = orderData.Currency
		} else {
			if createFrom != orderData.CreateFrom {
				errCode = "order_system_mark"
				err = errors.New("system mark not one")
				return
			}
			if currency != orderData.Currency {
				errCode = "order_currency"
				err = errors.New("currency not one")
				return
			}
		}
		priceTotal += orderData.Price
		//最小过期时间，作为最终过期时间
		if expireAt.Unix() < orderData.ExpireAt.Unix() {
			expireAt = orderData.ExpireAt
		}
		//构建支付创建人信息
		payUserID = orderData.UserID
		payOrgID = orderData.OrgID
	}
	//构建支付请求
	payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
		UserID:         payUserID,
		OrgID:          payOrgID,
		IsRefund:       false,
		Currency:       currency,
		Price:          priceTotal,
		PaymentChannel: args.PaymentChannel,
		ExpireAt:       expireAt,
		Des:            args.Des,
	})
	if err != nil {
		return
	}
	//修改上述所有订单的支付ID
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "pay_create", args.Des)
	if err != nil {
		errCode = "order_update_pay_id"
		return
	}
	//计算支付方式
	payFromSystem := fmt.Sprint(payData.PaymentChannel.System)
	if payData.PaymentChannel.Mark != "" {
		payFromSystem = payFromSystem + "_" + payData.PaymentChannel.Mark
	}
	var orderCompanyID int64 = 0
	switch payData.PaymentChannel.System {
	case "company_returned":
		orderCompanyID = payData.PaymentChannel.ID
	}
	//更新数据
	if _, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET pay_id = :pay_id, company_id = :company_id, pay_from = :pay_from, pay_list = array_append(pay_list, :pay_id), logs = logs || :log WHERE id = ANY(:ids)", map[string]interface{}{
		"ids":        args.IDs,
		"pay_id":     payData.ID,
		"company_id": orderCompanyID,
		"pay_from":   payFromSystem,
		"log":        newLog,
	}); err != nil {
		errCode = "order_update_pay_id"
		err = errors.New(fmt.Sprint("order update pay id, ", err))
		return
	}
	//分批更新扩展参数
	for _, v := range orderList {
		v.Params = CoreSQLConfig.Set(v.Params, "paySystem", payFromSystem)
		if _, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET params = :params WHERE id = :id", map[string]interface{}{
			"id":     v.ID,
			"params": v.Params,
		}); err != nil {
			CoreLog.Warn("service order create pay, update params, order id: ", v.ID, ", err: ", err)
			err = nil
		}
		//清理缓冲
		deleteOrderCache(v.ID)
	}
	//反馈
	return
}

// ArgsCheckPay 请求检查支付状态参数
type ArgsCheckPay struct {
	//订单ID组
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//组织ID
	// 可选，检测
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可选，检测
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// CheckPay 请求检查支付状态参数
func CheckPay(args *ArgsCheckPay) (b bool) {
	var dataList []FieldsOrder
	if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT price_pay, pay_status FROM service_order WHERE id = ANY($1) AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR user_id = $3)", args.IDs, args.OrgID, args.UserID); err != nil {
		return
	}
	for _, v := range dataList {
		if !v.PricePay {
			return
		}
		if v.PayStatus == 0 {
			return
		}
	}
	b = true
	return
}

// ArgsUpdatePayID 更新payID参数
type ArgsUpdatePayID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID" check:"id"`
}

// UpdatePayID 更新payID
func UpdatePayID(args *ArgsUpdatePayID) (err error) {
	//必须尚未支付
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: -1,
	})
	if err != nil {
		return
	}
	if orderData.PricePay {
		err = errors.New("order is pay")
		return
	}
	//获取支付ID
	payData, _ := FinancePay.GetID(&FinancePay.ArgsGetID{
		ID:       args.PayID,
		IsSecret: false,
	})
	//计算支付方式
	payFromSystem := fmt.Sprint(payData.PaymentChannel.System)
	if payData.PaymentChannel.Mark != "" {
		payFromSystem = payFromSystem + "_" + payData.PaymentChannel.Mark
	}
	var orderCompanyID int64 = 0
	switch payData.PaymentChannel.System {
	case "company_returned":
		orderCompanyID = payData.PaymentChannel.ID
	}
	//修正扩展参数
	orderData.Params = CoreSQLConfig.Set(orderData.Params, "paySystem", payFromSystem)
	//修改订单支付ID
	// 其他内容将检测支付状态
	var newLog string
	newLog, err = getLogData(0, args.OrgBindID, "update_pay_id", fmt.Sprint("修改支付ID，系统准备检测支付完成情况"))
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), pay_id = :pay_id, company_id = :company_id, pay_from = :pay_from, logs = logs || :log, params = :params WHERE id = :id", map[string]interface{}{
		"id":         orderData.ID,
		"pay_id":     args.PayID,
		"company_id": orderCompanyID,
		"pay_from":   payFromSystem,
		"log":        newLog,
		"params":     orderData.Params,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

// ArgsUpdatePrice 修改订单金额参数
type ArgsUpdatePrice struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//费用组成
	PriceList FieldsPrices `db:"price_list" json:"priceList"`
}

// UpdatePrice 修改订单金额参数
// 必须在付款之前修改
func UpdatePrice(args *ArgsUpdatePrice) (err error) {
	//获取订单信息
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: -1,
	})
	if err != nil {
		return
	}
	//订单全局必须尚未支付
	if orderData.PricePay {
		err = errors.New("order is pay")
		return
	}
	//补全费用构成
	if args.PriceList == nil {
		args.PriceList = orderData.PriceList
	} else {
		for _, v := range orderData.PriceList {
			isFind := false
			for _, v2 := range args.PriceList {
				if v.PriceType == v2.PriceType {
					isFind = true
					break
				}
			}
			if isFind {
				continue
			}
			args.PriceList = append(args.PriceList, v)
		}
	}
	//计算新的价格
	var newPrice int64
	for k, v := range args.PriceList {
		//计算新的费用
		newPrice += v.Price
		//修正支付状态
		if v.Price < 1 {
			args.PriceList[k].IsPay = true
		}
	}
	//更新金额信息
	var newLog string
	newLog, err = getLogData(0, args.OrgBindID, "update_price", fmt.Sprint("修改订单金额[", orderData.Price, "]，新的价格为[", newPrice, "]"))
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), price_list = :price_list, price = :price, logs = logs || :log WHERE id = :id", map[string]interface{}{
		"id":         orderData.ID,
		"price_list": args.PriceList,
		"price":      newPrice,
		"log":        newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//如果金额小于1，则标记支付完成
	if newPrice < 1 {
		err = PayFinish(&ArgsPayFinish{
			ID:        orderData.ID,
			OrgID:     args.OrgID,
			UserID:    0,
			OrgBindID: args.OrgBindID,
			Des:       "由于价格为0，支付自动完成",
		})
		if err != nil {
			return
		}
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	return
}

// ArgsPayFinish 支付成功参数
type ArgsPayFinish struct {
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
	//调整说明描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// PayFinish 支付成功
func PayFinish(args *ArgsPayFinish) (err error) {
	//获取订单
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		return
	}
	//调用内部函数
	err = payFinishByOrderID(orderData.ID, args.OrgBindID, 0, "pay_finish", args.Des)
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

func payFinishByOrderID(orderID int64, orgBindID int64, payID int64, mark string, logDes string) (err error) {
	//获取订单数据
	var orderData FieldsOrder
	orderData = getByID(orderID)
	if orderData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//标记payList全部完成
	for k, v := range orderData.PriceList {
		if !v.IsPay {
			orderData.PriceList[k].IsPay = true
		}
	}
	//修改订单
	var newLog string
	newLog, err = getLogData(0, orgBindID, mark, logDes)
	if err != nil {
		err = errors.New(fmt.Sprint("get log, ", err))
		return
	}
	//标记订单完成支付
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), pay_status = 1, price_pay = true, price_list = :price_list, logs = logs || :log WHERE id = :id AND pay_status != 1", map[string]interface{}{
		"id":         orderData.ID,
		"price_list": orderData.PriceList,
		"log":        newLog,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update order status, ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//触发订单抽取佣金设计
	if orderData.OrgID > 0 {
		orderSystem := ""
		for _, v := range orderData.Goods {
			orderSystem = v.From.System
			break
		}
		var cutPrice int64
		cutPrice, err = FinanceTakeCut.AddLog(&FinanceTakeCut.ArgsAddLog{
			OrgID:       orderData.OrgID,
			OrderSystem: orderSystem,
			OrderPrice:  orderData.Price,
			OrderID:     orderData.ID,
		})
		if err != nil {
			//CoreLog.Warn("order finish and finance take cut price failed, order id: ", orderData.ID, err)
			err = nil
		} else {
			if cutPrice > 0 {
				newLog, err = getLogData(0, orgBindID, "take_cut", fmt.Sprint("平台扣除佣金费用: ￥", float64(cutPrice)/100))
				if err != nil {
					return
				}
				_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), logs = logs || :log WHERE id = :id", map[string]interface{}{
					"id":  orderData.ID,
					"log": newLog,
				})
				if err != nil {
					err = errors.New(fmt.Sprint("update order by finance take cut price, ", err))
					return
				}
				//清理缓冲
				deleteOrderCache(orderData.ID)
			}
		}
	}
	//推送完成支付
	pushNatsOrderPay(orderData.ID, payID, orderData.Status >= 2)
	//自动审核处理
	if orderData.Status == 1 && orderData.AllowAutoAudit {
		//先审核订单
		err = UpdateAudit(&ArgsUpdateAudit{
			ID:        orderData.ID,
			OrgID:     -1,
			UserID:    -1,
			OrgBindID: orgBindID,
			Des:       "支付完成，自动审核订单",
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update audit, ", err))
			return
		}
	}
	//反馈
	return
}

// 通知nats支付完成
func pushNatsOrderPay(orderID int64, payID int64, haveAudit bool) {
	//通知完成支付
	CoreNats.PushDataNoErr("/service/order/pay", "finish", orderID, "", map[string]interface{}{
		"payID": payID,
	})
	//尝试通知审核并完成支付
	if haveAudit {
		pushOrderAuditAndPay(orderID)
	}
}

func payFinishByPayID(payID int64) (err error) {
	//获取符合条件的所有订单
	var orderList []FieldsOrder
	orderList, err = getListByPayID(payID)
	if err != nil {
		err = errors.New("no data")
		return
	}
	for _, v := range orderList {
		err = payFinishByOrderID(v.ID, 0, payID, "pay_finish", "支付完成")
		if err != nil {
			return
		}
		//清理缓冲
		deleteOrderCache(v.ID)
	}
	return
}

// ArgsPayFinancePhysical 采用财务实物支付订单参数
type ArgsPayFinancePhysical struct {
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
	//调整说明描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//抵扣物品数量集合
	Data []ArgsPayFinancePhysicalData `json:"data"`
}

// ArgsPayFinancePhysicalData 抵扣物品
type ArgsPayFinancePhysicalData struct {
	//获取来源
	// 如果商品mark带有virtual标记，且订单商品全部带有该标记，订单将在付款后直接完成
	From CoreSQLFrom.FieldsFrom `db:"from" json:"from"`
	//给予标的物数量
	PhysicalCount int64 `db:"physical_count" json:"physicalCount" check:"int64Than0"`
}

// PayFinancePhysical 采用财务实物支付订单
func PayFinancePhysical(args *ArgsPayFinancePhysical) (err error) {
	//获取订单信息
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		err = errors.New("no find order")
		return
	}
	//采用实物抵扣
	var logData []FinancePhysicalPay.ArgsCreateLogData
	for _, v := range orderData.Goods {
		err = FinancePhysicalPay.CheckPhysicalByFrom(&FinancePhysicalPay.ArgsGetPhysicalByFrom{
			OrgID:    args.OrgID,
			BindFrom: v.From,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("have no support physical data, org id: ", args.OrgID, ", mall product: ", v.From, ", err: ", err))
			return
		}
		isFind := false
		var physicalCount int64 = 0
		for _, v2 := range args.Data {
			if v2.From.CheckEg(v.From) {
				physicalCount = v2.PhysicalCount
				isFind = true
				break
			}
		}
		if !isFind {
			err = errors.New("need more physical count")
			return
		}
		logData = append(logData, FinancePhysicalPay.ArgsCreateLogData{
			PhysicalCount: physicalCount,
			BindFrom:      v.From,
			BindCount:     v.Count,
		})
	}
	var newIDs pq.Int64Array
	params := CoreSQLConfig.FieldsConfigsType{
		{
			Mark: "order_id",
			Val:  fmt.Sprint(orderData.ID),
		},
	}
	newIDs, err = FinancePhysicalPay.CreateLog(&FinancePhysicalPay.ArgsCreateLog{
		OrgID:  orderData.OrgID,
		BindID: args.OrgBindID,
		UserID: orderData.UserID,
		System: "order",
		Data:   logData,
		Params: params,
	})
	if err != nil {
		return
	}
	//标记完成订单支付
	if !orderData.PricePay {
		err = PayFinish(&ArgsPayFinish{
			ID:        args.ID,
			OrgID:     args.OrgID,
			UserID:    args.UserID,
			OrgBindID: args.OrgBindID,
			Des:       args.Des,
		})
		if err != nil {
			return
		}
	}
	//订单增加扩展参数
	var newIDsStr string
	for _, v := range newIDs {
		if newIDsStr == "" {
			newIDsStr = fmt.Sprint(v)
		} else {
			newIDsStr = fmt.Sprint(newIDsStr, ",", v)
		}
	}
	orderData.Params = CoreSQLConfig.Set(orderData.Params, "physical_log_ids", newIDsStr)
	orderData.Params = CoreSQLConfig.Set(orderData.Params, "paySystem", "physical_pay")
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), params = :params WHERE id = :id", map[string]interface{}{
		"id":     orderData.ID,
		"params": orderData.Params,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

// ArgsPayFailed 支付失败参数
type ArgsPayFailed struct {
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
	//调整说明描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// PayFailed 支付失败
func PayFailed(args *ArgsPayFailed) (err error) {
	//获取订单数据
	var orderData FieldsOrder
	orderData, err = GetByID(&ArgsGetByID{
		ID:     args.ID,
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		return
	}
	err = payFailed(orderData.ID, args.OrgBindID, "pay_failed", args.Des)
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

func payFailed(orderID int64, orgBindID int64, mark string, logDes string) (err error) {
	//获取订单信息
	var orderData FieldsOrder
	orderData = getByID(orderID)
	if orderData.ID < 1 {
		err = errors.New("no data")
		return
	}
	var newLog string
	newLog, err = getLogData(0, orgBindID, mark, logDes)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), pay_status = 4, price_pay = false, logs = logs || :log WHERE id = :id AND pay_status != 4", map[string]interface{}{
		"id":  orderData.ID,
		"log": newLog,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}
