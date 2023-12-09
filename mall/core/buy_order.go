package MallCore

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	MallLogMod "gitee.com/weeekj/weeekj_core/v5/mall/log/mod"
	ServiceOrderWait "gitee.com/weeekj/weeekj_core/v5/service/order/wait"
	ServiceOrderWaitFields "gitee.com/weeekj/weeekj_core/v5/service/order/wait_fields"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateOrder 开始下单参数
type ArgsCreateOrder struct {
	//商户ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//收货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//期望送货时间
	TransportTaskAt string `db:"transport_task_at" json:"transportTaskAt" check:"isoTime" empty:"true"`
	//备注信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//商品ID列
	Products []ArgsGetProductPriceProduct `db:"products" json:"products"`
	//会员配置ID
	// 只能指定一个
	UserSubID int64 `db:"user_sub_id" json:"userSubID" check:"id" empty:"true"`
	//票据
	// 可以使用的票据列，具体的配置在票据配置内进行设置
	// 票据分平台和商户，平台票据需参与活动才能使用，否则将自动禁止设置和后期使用
	UserTicket pq.Int64Array `db:"user_ticket" json:"userTicket" check:"ids" empty:"true"`
	//是否使用积分
	UseIntegral bool `db:"use_integral" json:"useIntegral" check:"bool"`
	//订单总费用
	// 总费用是否支付
	PricePay bool `db:"price_pay" json:"pricePay" check:"bool"`
	//是否允许货到付款？
	TransportPayAfter bool `db:"transport_pay_after" json:"transportPayAfter" check:"bool" empty:"true"`
	//是否绕过对库存限制
	SkipProductCountLimit bool `json:"skipProductCountLimit" check:"bool" empty:"true"`
	//推荐人手机号
	ReferrerNationCode string `db:"referrer_nation_code" json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
	//强制插入其他费用
	OtherPriceList ServiceOrderWaitFields.FieldsPrices `json:"otherPriceList"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//配送方式
	TransportType int `db:"transport_type" json:"transportType"`
}

// CreateOrder 开始下单
// 该方法包含了BuyProduct的处理，请勿同时使用。
// 可以使用BuyProduct作为释放库存处理
// 注意，订单不根据商品地址拆分，拆分配送将依赖于额外的物流处理模块进行调配货物
func CreateOrder(args *ArgsCreateOrder) (orderWaitData ServiceOrderWaitFields.FieldsWait, errCode string, err error) {
	//计算送货时间
	var transportTaskAt time.Time
	transportTaskAt, err = CoreFilter.GetTimeByISO(args.TransportTaskAt)
	if err != nil {
		transportTaskAt = CoreFilter.GetNowTime()
		//errCode = "transport_task_at"
		//return
	}
	//预检查，计算费用内容
	var buyWaitData DataProductPrice
	buyWaitData, errCode, err = GetProductPrice(&ArgsGetProductPrice{
		Products:              args.Products,
		OrgID:                 args.OrgID,
		UserID:                args.UserID,
		UserSubID:             args.UserSubID,
		UserTicket:            args.UserTicket,
		UseIntegral:           args.UseIntegral,
		Address:               args.Address,
		SkipProductCountLimit: args.SkipProductCountLimit,
		TransportType:         args.TransportType,
	})
	if err != nil {
		return
	}
	if len(buyWaitData.Goods) < 1 {
		errCode = "err_buy_empty"
		err = errors.New("goods not exist")
		return
	}
	//费用构成
	priceList := ServiceOrderWaitFields.FieldsPrices{}
	/**
	//购物车最终价格已经包含配送费，所以需要减去配送费、增加已经减去的优惠价格
	buyWaitData.LastPrice = 0
	for _, v := range buyWaitData.Goods {
		buyWaitData.LastPrice += v.Price * v.Count
	}
	*/
	//构建价格列
	priceList = append(priceList, ServiceOrderWaitFields.FieldsPrice{
		PriceType: 0,
		IsPay:     buyWaitData.LastPriceProduct < 1,
		Price:     buyWaitData.LastPriceProduct,
	})
	priceList = append(priceList, ServiceOrderWaitFields.FieldsPrice{
		PriceType: 1,
		IsPay:     buyWaitData.LastPriceTransport < 1,
		Price:     buyWaitData.LastPriceTransport,
	})
	if len(args.OtherPriceList) > 0 {
		for _, v := range args.OtherPriceList {
			isFind := false
			for k2, v2 := range priceList {
				if v.PriceType == v2.PriceType {
					priceList[k2].IsPay = v.IsPay
					priceList[k2].Price = v.Price
					isFind = true
					break
				}
			}
			if !isFind {
				priceList = append(priceList, v)
			}
		}
	}
	//增加商品销量
	newGoods := ServiceOrderWaitFields.FieldsGoods{}
	for _, v := range buyWaitData.Goods {
		//跳过非商品
		if v.From.System != "mall" {
			continue
		}
		//直接减少该商品库存
		if !args.SkipProductCountLimit {
			err = UpdateProductAddCount(v.From.ID, v.OptionKey, 0-int(v.Count))
			if err != nil {
				CoreLog.Error("update product need count, mall id: ", v.From.ID, ", need count: ", v.Count, " , err: ", err)
				errCode = "err_mall_product_count"
				err = errors.New(fmt.Sprint("product not have count, ", err))
				return
			}
		}
		newGoods = append(newGoods, v)
	}
	buyWaitData.Goods = newGoods
	//创建订单
	orderWaitData, errCode, err = ServiceOrderWait.CreateOrder(&ServiceOrderWait.ArgsCreateOrder{
		SystemMark:         "mall",
		OrgID:              args.OrgID,
		UserID:             args.UserID,
		CreateFrom:         args.CreateFrom,
		AddressFrom:        buyWaitData.ProductList[0].Address,
		AddressTo:          args.Address,
		Goods:              buyWaitData.Goods,
		Exemptions:         buyWaitData.Exemptions,
		NeedAllowAutoAudit: false,
		AllowAutoAudit:     false,
		TransportAllowAuto: true,
		TransportTaskAt:    transportTaskAt,
		TransportPayAfter:  args.TransportPayAfter,
		TransportSystem:    getTransportType(args.TransportType),
		PriceList:          priceList,
		PricePay:           args.PricePay,
		NeedExPrice:        false,
		Currency:           86,
		Des:                args.Des,
		Logs:               []ServiceOrderWaitFields.FieldsLog{},
		ReferrerNationCode: args.ReferrerNationCode,
		ReferrerPhone:      args.ReferrerPhone,
		Params:             args.Params,
	})
	if err == nil {
		//添加统计
		for _, v := range buyWaitData.Goods {
			_ = appendAnalysisBuy(args.OrgID, v.From.ID, args.UserID, int(v.Count))
		}
		//增加用户日志统计
		//CoreLog.Info("user buy data: ", buyWaitData)
		//添加操作日志
		for _, v := range buyWaitData.Goods {
			if v.From.System != "mall" {
				continue
			}
			MallLogMod.AppendLog(args.UserID, "", args.OrgID, v.From.ID, 3)
		}
	}
	//反馈
	return
}

// 识别分类
// 0 self 其他配送; 1 take 自提; 2 transport 自运营配送; 3 running 跑腿服务; 4 housekeeping 家政服务
func getTransportType(transportType int) string {
	switch transportType {
	case 0:
		//第三方配送
		return "self"
	case 1:
		//自提
		return "take"
	case 2:
		//自运营配送
		return "transport"
	case 3:
		//跑腿
		return "running"
	case 4:
		//家政服务
		return "housekeeping"
	default:
		//第三方配送
		return "self"
	}
}
