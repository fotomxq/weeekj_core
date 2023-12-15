package ServiceInfoExchange

import (
	"errors"
	"fmt"
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	ServiceOrderWait "github.com/fotomxq/weeekj_core/v5/service/order/wait"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
	"time"
)

// ArgsCreateInfoOrder 创建信息订单参数
type ArgsCreateInfoOrder struct {
	//商户ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//购买用户
	BuyUserID int64 `db:"buy_user_id" json:"buyUserID" check:"id"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//收货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//备注信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id"`
	//订单总费用
	// 总费用是否支付
	PricePay bool `db:"price_pay" json:"pricePay" check:"bool"`
	//是否允许货到付款？
	TransportPayAfter bool `db:"transport_pay_after" json:"transportPayAfter" check:"bool" empty:"true"`
	//推荐人手机号
	ReferrerNationCode string `db:"referrer_nation_code" json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateInfoOrder 创建信息订单
func CreateInfoOrder(args *ArgsCreateInfoOrder) (orderWaitData ServiceOrderWaitFields.FieldsWait, errCode string, err error) {
	//信息数据
	var infoData FieldsInfo
	infoData, err = GetInfoPublishID(&ArgsGetInfoID{
		ID:     args.InfoID,
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		errCode = "info_not_exist"
		return
	}
	//创建订单
	orderWaitData, errCode, err = ServiceOrderWait.CreateOrder(&ServiceOrderWait.ArgsCreateOrder{
		SystemMark:  "mall",
		OrgID:       infoData.OrgID,
		UserID:      args.BuyUserID,
		CreateFrom:  args.CreateFrom,
		AddressFrom: infoData.Address,
		AddressTo:   args.Address,
		Goods: ServiceOrderWaitFields.FieldsGoods{
			{
				From: CoreSQLFrom.FieldsFrom{
					System: "service_info_exchange",
					ID:     infoData.ID,
					Mark:   infoData.InfoType,
					Name:   infoData.Title,
				},
				OptionKey:  "",
				Count:      1,
				Price:      infoData.Price,
				Exemptions: ServiceOrderWaitFields.FieldsExemptions{},
			},
		},
		Exemptions:         ServiceOrderWaitFields.FieldsExemptions{},
		NeedAllowAutoAudit: false,
		AllowAutoAudit:     false,
		TransportAllowAuto: true,
		TransportTaskAt:    CoreFilter.GetNowTime(),
		TransportPayAfter:  args.TransportPayAfter,
		PriceList: ServiceOrderWaitFields.FieldsPrices{
			{
				PriceType: 0,
				IsPay:     args.PricePay,
				Price:     infoData.Price,
			},
		},
		PricePay:           args.PricePay,
		NeedExPrice:        false,
		Currency:           infoData.Currency,
		Des:                args.Des,
		Logs:               []ServiceOrderWaitFields.FieldsLog{},
		ReferrerNationCode: args.ReferrerNationCode,
		ReferrerPhone:      args.ReferrerPhone,
		Params:             args.Params,
	})
	if err != nil {
		return
	}
	//更新wait orderID
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_info_exchange SET wait_order_id = :wait_order_id WHERE id = :id", map[string]interface{}{
		"id":            infoData.ID,
		"wait_order_id": orderWaitData.ID,
	})
	if err != nil {
		return
	}
	//清除缓冲
	deleteInfoCache(infoData.ID)
	//反馈
	return
}

// ArgsUpdateInfoOrderPrice 修改订单费用参数
type ArgsUpdateInfoOrderPrice struct {
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id"`
	//组织ID
	// 可选，作为验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//日志
	//操作组织人员ID
	// 如果留空则说明为系统自动调整或创建人产生
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//费用组成
	PriceList ServiceOrderMod.FieldsPrices `db:"price_list" json:"priceList"`
}

// UpdateInfoOrderPrice 修改订单费用
func UpdateInfoOrderPrice(args *ArgsUpdateInfoOrderPrice) (err error) {
	//信息数据
	var infoData FieldsInfo
	infoData, err = GetInfoPublishID(&ArgsGetInfoID{
		ID:     args.InfoID,
		OrgID:  args.OrgID,
		UserID: args.UserID,
	})
	if err != nil {
		return
	}
	if infoData.OrderID < 1 {
		err = errors.New("info no order")
		return
	}
	//清除缓冲
	deleteInfoCache(infoData.ID)
	//修改价格
	ServiceOrderMod.UpdatePrice(ServiceOrderMod.ArgsUpdatePrice{
		ID:        infoData.OrderID,
		OrgID:     infoData.OrgID,
		OrgBindID: 0,
		PriceList: args.PriceList,
	})
	//反馈
	return
}

// 修改等待订单
func updateInfoOrderID(infoID int64, orderID int64) (err error) {
	//更新wait orderID
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_info_exchange SET wait_order_id = 0, order_id = :order_id WHERE id = :id", map[string]interface{}{
		"id":       infoID,
		"order_id": orderID,
	})
	if err != nil {
		return
	}
	//清除缓冲
	deleteInfoCache(infoID)
	//反馈
	return
}

// 修改订单完成
func updateInfoOrderFinish(infoID int64, orderID int64) (err error) {
	//锁定机制
	orderFinishLock.Lock()
	defer orderFinishLock.Unlock()
	//获取订单信息
	orderData := ServiceOrderMod.GetByIDNoErr(orderID)
	if orderData.ID < 1 {
		err = errors.New("order not exist")
		return
	}
	//获取数据
	infoData := getInfoByID(infoID)
	//避免重复完成
	if infoData.OrderFinish {
		return
	}
	//更新
	if Router2SystemConfig.GlobConfig.Service.InfoExchangeOrderFinishAutoDown {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_info_exchange SET order_finish = true, publish_at = to_timestamp(0) WHERE id = :id", map[string]interface{}{
			"id": infoID,
		})
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_info_exchange SET order_finish = true WHERE id = :id", map[string]interface{}{
			"id": infoID,
		})
	}
	if err != nil {
		return
	}
	//清除缓冲
	deleteInfoCache(infoID)
	//获取数据
	infoData = getInfoByID(infoID)
	if infoData.ID > 0 {
		//统计行为
		AnalysisAny2.AppendData("add", "service_info_exchange_order_price", time.Time{}, infoData.OrgID, infoData.UserID, 0, 0, 0, infoData.Price)
		AnalysisAny2.AppendData("add", "service_info_exchange_order_count", time.Time{}, infoData.OrgID, infoData.UserID, 0, 0, 0, 1)
	}
	//划拨资金，用户储蓄已经划拨了资金到平台，需将平台资金划拨给发起人
	if infoData.UserID > 0 && infoData.Price > 0 {
		var errCode string
		_, errCode, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
			UpdateHash: "",
			CreateInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     infoData.UserID,
				Mark:   "",
				Name:   "",
			},
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     0,
				Mark:   "",
				Name:   "",
			},
			//TODO: 暂时写死，后续将改为第二代储蓄模块
			ConfigMark:      "savings",
			AppendSavePrice: infoData.Price,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("service info ex order finish, set user finance deposit failed, info id: ", infoData.ID, ", user id: ", infoData.UserID, ", add price: ", infoData.Price, ", err: ", errCode, ", ", err))
			return
		}
	}
	//反馈
	return
}
