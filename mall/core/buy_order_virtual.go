package MallCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	MallLogMod "github.com/fotomxq/weeekj_core/v5/mall/log/mod"
	ServiceOrderWait "github.com/fotomxq/weeekj_core/v5/service/order/wait"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
	"github.com/lib/pq"
)

// ArgsBuyOrderVirtual 虚拟订单快速构建参数
type ArgsBuyOrderVirtual struct {
	//订单创建时间
	CreateAt string `db:"create_at" json:"createAt" check:"defaultTime"`
	//商户ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//采购商
	BuyCompanyID int64 `json:"buyCompanyID" check:"id" empty:"true"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//收货地址
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address"`
	//备注信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//商品ID列
	Products []ArgsGetProductPriceProductVirtual `db:"products" json:"products"`
	//会员配置ID
	// 只能指定一个
	UserSubID int64 `db:"user_sub_id" json:"userSubID" check:"id" empty:"true"`
	//票据
	// 可以使用的票据列，具体的配置在票据配置内进行设置
	// 票据分平台和商户，平台票据需参与活动才能使用，否则将自动禁止设置和后期使用
	UserTicket pq.Int64Array `db:"user_ticket" json:"userTicket" check:"ids" empty:"true"`
	//是否使用积分
	UseIntegral bool `db:"use_integral" json:"useIntegral" check:"bool"`
	//强制插入其他费用
	PriceList []ArgsPriceVirtual `json:"priceList"`
	//订单总费用
	PriceReal int64 `db:"price_real" json:"priceReal" check:"price" empty:"true"`
	//推荐人手机号
	ReferrerNationCode string `db:"referrer_nation_code" json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//配送方式
	TransportType int `db:"transport_type" json:"transportType"`
}

type ArgsPriceVirtual struct {
	//费用类型
	// 0 货物费用；1 配送费用；2 保险费用; 3 跑腿费用
	PriceType int `db:"price_type" json:"priceType" check:"mark"`
	//总金额
	Price int64 `db:"price" json:"price" check:"price"`
}

type ArgsGetProductPriceProductVirtual struct {
	//商品ID
	ID int64 `db:"id" json:"id" check:"id"`
	//选项Key
	// 如果给空，则该商品必须也不包含选项
	OptionKey string `db:"option_key" json:"optionKey" check:"mark" empty:"true"`
	//购买数量
	// 如果为0，则只判断单价的价格
	BuyCount int `db:"buy_count" json:"buyCount" check:"int64Than0"`
	//进货价
	PriceIn int64 `db:"price_in" json:"priceIn" check:"price" empty:"true"`
	//销售价格
	PriceOut int64 `db:"price_out" json:"priceOut" check:"price" empty:"true"`
}

//BuyOrderVirtual 虚拟订单快速构建
/**
1. 仅适用于特定场景下
2. 用于快速构建订单，并完成订单
3. 订单不会产生真实收益，主要用于数据同步、手动填录等场景
*/
func BuyOrderVirtual(args *ArgsBuyOrderVirtual) (errCode string, err error) {
	//重组货物
	var productList []ArgsGetProductPriceProduct
	for _, v := range args.Products {
		productList = append(productList, ArgsGetProductPriceProduct{
			ID:        v.ID,
			OptionKey: v.OptionKey,
			BuyCount:  v.BuyCount,
		})
	}
	//预检查，计算费用内容
	var buyWaitData DataProductPrice
	buyWaitData, errCode, err = GetProductPrice(&ArgsGetProductPrice{
		Products:              productList,
		OrgID:                 args.OrgID,
		UserID:                args.UserID,
		UserSubID:             args.UserSubID,
		UserTicket:            args.UserTicket,
		UseIntegral:           args.UseIntegral,
		Address:               args.Address,
		SkipProductCountLimit: true,
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
	//构建价格列
	if len(args.PriceList) > 0 {
		for _, v := range args.PriceList {
			isFind := false
			for k2, v2 := range priceList {
				if v.PriceType == v2.PriceType {
					priceList[k2].IsPay = false
					priceList[k2].Price = v.Price
					isFind = true
					break
				}
			}
			if !isFind {
				priceList = append(priceList, ServiceOrderWaitFields.FieldsPrice{
					PriceType: v.PriceType,
					IsPay:     false,
					Price:     v.Price,
				})
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
		//修改销售价格
		for _, v2 := range args.Products {
			if v.From.ID == v2.ID {
				v.Price = v2.PriceOut
				break
			}
		}
		//加入数据
		newGoods = append(newGoods, v)
	}
	//追加虚拟订单标记
	args.Params = append(args.Params, CoreSQLConfig.FieldsConfigType{
		Mark: "virtual_sync",
		Val:  "true",
	})
	//创建订单
	var orderWaitData ServiceOrderWaitFields.FieldsWait
	orderWaitData, errCode, err = ServiceOrderWait.CreateOrder(&ServiceOrderWait.ArgsCreateOrder{
		SystemMark:         "mall",
		OrgID:              args.OrgID,
		UserID:             args.UserID,
		CreateFrom:         args.CreateFrom,
		AddressFrom:        buyWaitData.ProductList[0].Address,
		AddressTo:          args.Address,
		Goods:              newGoods,
		Exemptions:         buyWaitData.Exemptions,
		NeedAllowAutoAudit: false,
		AllowAutoAudit:     false,
		TransportAllowAuto: true,
		TransportTaskAt:    CoreFilter.GetNowTime(),
		TransportPayAfter:  false,
		TransportSystem:    getTransportType(args.TransportType),
		PriceList:          priceList,
		PricePay:           false,
		NeedExPrice:        false,
		Currency:           86,
		Des:                args.Des,
		Logs: []ServiceOrderWaitFields.FieldsLog{
			{
				CreateAt:  CoreFilter.GetNowTime(),
				UserID:    0,
				OrgBindID: 0,
				Mark:      "create",
				Des:       fmt.Sprint("创建虚拟订单，用于同步和手动填录数据"),
			},
		},
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
	//推送虚构订单情况
	CoreNats.PushDataNoErr("/service/order/create_wait_virtual", "finish", orderWaitData.ID, "", map[string]any{
		"products":  args.Products,
		"companyID": args.BuyCompanyID,
	})
	//反馈
	return
}
