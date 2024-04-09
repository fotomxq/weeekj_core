package ServiceOrder

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	MallCoreMod "github.com/fotomxq/weeekj_core/v5/mall/core/mod"
	OrgSubscriptionMod "github.com/fotomxq/weeekj_core/v5/org/subscription/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceHousekeepingMod "github.com/fotomxq/weeekj_core/v5/service/housekeeping/mod"
	TMSTransportMod "github.com/fotomxq/weeekj_core/v5/tms/transport/mod"
	TMSUserRunningMod "github.com/fotomxq/weeekj_core/v5/tms/user_running/mod"
	UserSubscriptionMod "github.com/fotomxq/weeekj_core/v5/user/subscription/mod"
	"time"
)

// ArgsUpdatePost 提交审核参数
type ArgsUpdatePost struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}

// UpdatePost 提交审核
func UpdatePost(args *ArgsUpdatePost) (err error) {
	var newLog string
	newLog, err = getLogData(args.UserID, args.OrgBindID, "post", args.Des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), status = 1, logs = logs || :log WHERE id = :id AND status = 0 AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":      args.ID,
		"org_id":  args.OrgID,
		"user_id": args.UserID,
		"log":     newLog,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update order id: ", args.ID, ", err: ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(args.ID)
	//如果启动自动审核，自动完成审核
	// 必须同时是货到付款启动才能执行，否则付款后会自动检测第二次
	var data FieldsOrder
	data = getByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//如果是货到付款，则直接进入审核状态
	if data.TransportPayAfter {
		err = updateAudit(data.ID, 0, "订单货到付款，自动提交审核")
		if err != nil {
			return
		}
		return
	} else {
		//如果不是货到付款，则检查是否自动过审核？
		if data.AllowAutoAudit && (data.Price < 1 || data.PricePay) {
			err = updateAudit(data.ID, 0, "订单已经支付，自动提交审核")
			if err != nil {
				return
			}
			return
		}
	}
	//反馈
	return
}

// ArgsUpdateAudit 审核订单参数
type ArgsUpdateAudit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//调整说明描述
	Des string `db:"des" json:"des"`
}

// UpdateAudit 审核订单
func UpdateAudit(args *ArgsUpdateAudit) (err error) {
	//获取订单
	orderData := getByID(args.ID)
	if orderData.ID < 1 || CoreSQL.CheckTimeHaveData(orderData.DeleteAt) || !CoreFilter.EqID2(args.OrgID, orderData.OrgID) || !CoreFilter.EqID2(args.UserID, orderData.UserID) {
		err = errors.New("no data")
		return
	}
	//更新状态
	err = updateAudit(orderData.ID, args.OrgBindID, args.Des)
	if err != nil {
		return
	}
	//反馈
	return
}

// 提交订单审核
func updateAudit(orderID int64, orgBindID int64, des string) (err error) {
	//获取订单
	orderData := getByID(orderID)
	if orderData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//更新状态
	var newLog string
	newLog, err = getLogData(orderData.UserID, orgBindID, "audit", des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), status = 2, logs = logs || :log WHERE id = :id AND status = 1", map[string]interface{}{
		"id":  orderData.ID,
		"log": newLog,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update order id: ", orderData.ID, ", err: ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//检查并自动生成配送单
	updateAuditAuto(orderData.ID)
	//尝试通知审核并完成支付
	if orderData.PayStatus == 1 {
		pushOrderAuditAndPay(orderData.ID)
	}
	//反馈
	return
}

// updateAuditAuto 提交审核后的收尾工作
func updateAuditAuto(orderID int64) {
	var err error
	//日志
	logAppend := fmt.Sprint("service order update audit auto tms, order id: ", orderID, ", ")
	//获取订单数据
	orderData := getByID(orderID)
	if orderData.ID < 1 {
		CoreLog.Warn(logAppend, "order not exist, order id: ", orderID)
		return
	}
	//检查订单商品类型
	// 是否需要配送单
	needTransport := false
	//遍历商品，处理单一性商品处理规则
	// 该规则适用于每个商品对应一个服务体系
	for _, vGood := range orderData.Goods {
		switch vGood.From.System {
		case "mall":
			//家政服务系统或配送服务单
			switch vGood.From.Mark {
			case "housekeeping":
				//创建服务单
				// 注意禁止混合购买行为
				updateAuditAutoHousekeeping(&orderData)
				//禁止创建后续配送单，缴费家政服务系统处理
				needTransport = false
			case "virtual":
				//虚拟商品，丢给商品模块处理
				// 注意禁止混合购买行为
				MallCoreMod.SendProductVirtual(vGood.From.ID, vGood.Count, orderData.UserID, orderData.OrgID, orderData.ID)
				//禁止创建后续配送单，虚拟商品不需要发货
				needTransport = false
			default:
				//创建普通服务单
				needTransport = true
			}
		case "user_sub":
			//赠送用户指定会员
			// 注意禁止混合购买行为
			err = UserSubscriptionMod.SetSubAdd(&UserSubscriptionMod.ArgsSetSubAdd{
				ConfigID: vGood.From.ID,
				UserID:   orderData.UserID,
				Unit:     int(vGood.Count),
				OrderID:  orderData.ID,
			})
			if err != nil {
				CoreLog.Warn(logAppend, "set user sub add, ", err)
			}
			//禁止创建后续配送单，会员系统类似虚拟商品直接完成授权，不需要配送单
			needTransport = false
			continue
		case "org_sub":
			//赠送用户指定会员
			// 注意禁止混合购买行为
			orgID := orderData.Params.GetValInt64NoBool("org_id")
			if orgID > 0 {
				err = OrgSubscriptionMod.SetSubAdd(&OrgSubscriptionMod.ArgsSetSubAdd{
					ConfigID: vGood.From.ID,
					OrgID:    orgID,
					Unit:     int(vGood.Count),
					OrderID:  orderData.ID,
				})
			}
			if err != nil {
				CoreLog.Warn(logAppend, "set org sub add, ", err)
			}
			//禁止创建后续配送单，会员系统类似虚拟商品直接完成授权，不需要配送单
			needTransport = false
			continue
		case "service_info_exchange":
			//帖子交易
			// 注意禁止混合购买行为
			// 禁止创建后续配送单，让用户自行处理
			needTransport = false
			//case "core_api":
			//	//API服务增加次数
			//	addMaxCount := orderData.Params.GetValInt64NoBool("addMaxCount")
			//	if addMaxCount > 0 && len(orderData.Goods) > 0 {
			//		BaseAPIMod.PushUpdateUseAnalysisAndOrderID(orderData.ID, BaseAPIMod.ArgsUpdateAddUseAnalysis{
			//			OrgID:       orderData.OrgID,
			//			UserID:      orderData.UserID,
			//			Mark:        orderData.Goods[0].OptionKey,
			//			AddMaxCount: addMaxCount,
			//		})
			//	}
			//	//不需要配送
			//	needTransport = false
		}
	}
	//如果需要配送单
	if needTransport && orderData.TransportID < 1 {
		if err = updateAuditAutoTransport(&orderData); err != nil {
			CoreLog.Error(logAppend, "auto transport failed, ", err)
		}
	}
}

// 自动创建配送单
func updateAuditAutoTransport(orderData *FieldsOrder) (err error) {
	//是否启动自动生成配送单服务
	var serviceOrderAutoSelfTransport bool
	serviceOrderAutoSelfTransport, err = BaseConfig.GetDataBool("ServiceOrderAutoSelfTransport")
	if err != nil {
		serviceOrderAutoSelfTransport = false
	}
	if !serviceOrderAutoSelfTransport {
		return
	}
	//生成配送单
	switch orderData.TransportSystem {
	case "self":
		//不做任何处理，等待商家自己配送即可
	case "take":
		//不做任何处理，但交给客户和商家互动处理
		//该设计如果启动了相关配置，会在接口层面先使用order_take处理验证工作，然后才能标记订单的后续细节
	case "transport":
		//重组获取非虚拟商品，作为配送货物
		var goods []TMSTransportMod.FieldsTransportGood
		if orderData.Goods != nil {
			for _, vGood := range orderData.Goods {
				if vGood.From.System != "mall" {
					continue
				}
				if vGood.From.Mark == "housekeeping" {
					continue
				}
				if vGood.From.Mark == "virtual" {
					continue
				}
				goods = append(goods, TMSTransportMod.FieldsTransportGood{
					System: vGood.From.System,
					ID:     vGood.From.ID,
					Mark:   vGood.From.Mark,
					Name:   vGood.From.Name,
					Count:  int(vGood.Count),
				})
			}
		}
		if len(goods) < 1 {
			return
		}
		//组装扩展参数
		var params []CoreSQLConfig.FieldsConfigType
		if orderData.Des != "" {
			params = CoreSQLConfig.Set(params, "des", orderData.Des)
		}
		paySystem, b := orderData.Params.GetVal("paySystem")
		if !b || paySystem == "" {
			paySystem = "order"
		}
		params = CoreSQLConfig.Set(params, "paySystem", paySystem)
		TMSTransportMod.CreateTransport(TMSTransportMod.ArgsCreateTransport{
			OrgID:       orderData.OrgID,
			BindID:      0,
			InfoID:      0,
			UserID:      orderData.UserID,
			FromAddress: orderData.AddressFrom,
			ToAddress:   orderData.AddressTo,
			OrderID:     orderData.ID,
			Goods:       goods,
			Weight:      0,
			Length:      0,
			Width:       0,
			Currency:    orderData.Currency,
			Price:       orderData.Price,
			PayFinish:   orderData.PricePay,
			TaskAt:      CoreFilter.GetISOByTime(orderData.TransportTaskAt),
			Params:      params,
		})
		//请求服务单
		_ = addLog(orderData.ID, fmt.Sprint("请求创建配送单"))
	case "running":
		//构建跑腿单
		if orderData.Des == "" {
			orderData.Des = "订单商品"
		}
		orderPayAllPrice := false
		for _, v := range orderData.PriceList {
			if v.PriceType == 1 {
				if orderData.PricePay {
					orderPayAllPrice = true
				}
			}
		}
		TMSUserRunningMod.CreateMission(TMSUserRunningMod.ArgsCreateMission{
			RunType:          2,
			WaitAt:           CoreFilter.GetISOByTime(orderData.TransportTaskAt),
			GoodType:         "order",
			OrgID:            0,
			UserID:           orderData.UserID,
			OrderID:          orderData.ID,
			RunWaitPrice:     -1,
			RunPayAfter:      true,
			OrderPayAllPrice: orderPayAllPrice,
			Des:              orderData.Des,
			GoodWidget:       0,
			FromAddress:      orderData.AddressFrom,
			ToAddress:        orderData.AddressTo,
			Params:           []CoreSQLConfig.FieldsConfigType{},
		})
	case "housekeeping":
		//updateAuditAuto模块已经处理，根据商品类型构建了服务单
	}
	//反馈
	return
}

// 自动创建服务单
func updateAuditAutoHousekeeping(orderData *FieldsOrder) {
	//家政服务
	//是否启动自动生成配送单服务
	serviceOrderAutoSelfHousekeeping, err := BaseConfig.GetDataBool("ServiceOrderAutoSelfHousekeeping")
	if err != nil {
		serviceOrderAutoSelfHousekeeping = false
	}
	if !serviceOrderAutoSelfHousekeeping {
		return
	}
	//检查订单是否创建过服务
	if b := checkOrderHaveHousekeeping(orderData.ID); b {
		return
	}
	//遍历商品创建
	for _, vGood := range orderData.Goods {
		if vGood.From.Mark == "virtual" {
			continue
		}
		if vGood.From.System != "mall" {
			continue
		}
		var payAt time.Time
		isPay := "false"
		if orderData.PricePay {
			payAt = CoreFilter.GetNowTime()
		}
		if orderData.TransportTaskAt.Unix() < 100000 {
			orderData.TransportTaskAt = CoreFilter.GetNowTimeCarbon().AddHours(1).Time
		}
		ServiceHousekeepingMod.CreateLog(ServiceHousekeepingMod.ArgsCreateLog{
			UserID:        orderData.UserID,
			NeedAt:        orderData.TransportTaskAt,
			OrgID:         orderData.OrgID,
			BindID:        0,
			OtherBinds:    []int64{},
			MallProductID: vGood.From.ID,
			OrderID:       orderData.ID,
			Currency:      orderData.Currency,
			Price:         orderData.Price,
			PayAt:         payAt,
			Des:           orderData.Des,
			Address:       orderData.AddressTo,
			ConfigID:      0,
			Params: CoreSQLConfig.FieldsConfigsType{
				{
					Mark: "orderPrice",
					Val:  fmt.Sprint(orderData.Price),
				},
				{
					Mark: "orderPay",
					Val:  isPay,
				},
			},
		})
		//记录订单日志
		_ = addLog(orderData.ID, fmt.Sprint("请求创建服务单"))
	}
}

// 更新配送状态信息
func updateTransportStatus(id int64, status int) {
	_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET transport_status = :transport_status WHERE id = :id", map[string]interface{}{
		"id":               id,
		"transport_status": status,
	})
	//清理缓冲
	deleteOrderCache(id)
}

// 通知订单支付完成并审核通过
func pushOrderAuditAndPay(id int64) {
	data := getByID(id)
	if data.ID < 1 || data.Status < 2 || !data.PricePay {
		return
	}
	CoreNats.PushDataNoErr("service_order_next", "/service/order/next", "", id, "", nil)
}

// 记录新的日志
func addLog(orderID int64, des string) (err error) {
	var newLog string
	newLog, err = getLogData(orderID, 0, "log", des)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET logs = logs || :log WHERE id = :id", map[string]interface{}{
		"id":  orderID,
		"log": newLog,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update order id: ", orderID, ", err: ", err))
		return
	}
	//清理缓冲
	deleteOrderCache(orderID)
	//反馈
	return
}

// 修改订单扩展参数
func updateOrderParams(orderID int64, params []CoreSQLConfig.FieldsConfigType) (err error) {
	//获取订单
	orderData := getByID(orderID)
	if orderData.ID < 1 {
		err = errors.New("order not exist")
		return
	}
	//更新数据
	for _, v := range params {
		orderData.Params = CoreSQLConfig.Set(orderData.Params, v.Mark, v.Val)
	}
	_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET params = :params WHERE id = :id", map[string]interface{}{
		"id":     orderData.ID,
		"params": orderData.Params,
	})
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}

// FiledGoodsPrice 编辑商品价格
type FiledGoodsPrice struct {
	//商品ID
	ID int64 `json:"id"`
	//价格
	Price int64 `db:"price" json:"price"`
	// 数量
	Count int64 `db:"count" json:"count"`
}

// ArgsUpdateOrder 编辑订单
type ArgsUpdateOrder struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	// 订单ID
	ID int64 `db:"id" json:"id"`
	//货物清单
	Goods []FiledGoodsPrice `db:"goods" json:"goods"`
}

// UpdateOrderGoodsPrice 编辑订单
func UpdateOrderGoodsPrice(args *ArgsUpdateOrder) (err error) {
	//获取订单
	var orderData FieldsOrder
	//获取订单
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
		// 判断支付来源如果是公司支付则可以修改
		paySystem := orderData.Params.GetValNoErr("paySystem")
		if paySystem != "company_returned" {
			err = errors.New(fmt.Sprint("order is pay: ", orderData.ID))
			return
		}
	}
	// 判断订单配送状态状态
	if orderData.Status >= 3 {
		err = errors.New(fmt.Sprint("update order status >= 3 id: ", orderData.ID))
		return
	}
	//新的商品总价
	var newGoodsTotalPrice int64
	// 新的总价格
	var newTotalPrice int64
	// 获取清单, 遍历订单清单 更新价格
	for k, v := range orderData.Goods {
		for _, v2 := range args.Goods {
			if v2.ID == v.From.ID {
				// 计算新的总价
				newGoodsTotalPrice += v2.Price * v2.Count
				// 更新价格
				orderData.Goods[k].Price = v2.Price
				break
			}
		}
	}
	for k, v := range orderData.PriceList {
		if v.PriceType == 0 {
			orderData.PriceList[k].Price = newGoodsTotalPrice
		}
		// 计算包含配送费、保险费其他费用
		newTotalPrice += orderData.PriceList[k].Price
	}
	// 更新总价
	orderData.Price = newTotalPrice
	orderData.PriceTotal = newTotalPrice
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_order SET update_at = NOW(), goods = :goods, price = :price, price_total = :price_total, price_list = :price_list WHERE id = :id", map[string]interface{}{
		"id":          orderData.ID,
		"goods":       orderData.Goods,
		"price":       orderData.Price,
		"price_total": orderData.PriceTotal,
		"price_list":  orderData.PriceList,
	})
	//清理缓冲
	deleteOrderCache(orderData.ID)
	//反馈
	return
}
