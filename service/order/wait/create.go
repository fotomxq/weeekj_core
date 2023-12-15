package ServiceOrderWait

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"time"
)

//创建订单请求
// 请求将直接进入列队，需等待订单服务核心确认
// 自带排重处理，10秒内禁止重复提交相同来源和货品的内容

type ArgsCreateOrder struct {
	//来源系统
	// 该订单创建来源的系统
	// eg: user_sub / org_sub / mall
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
	Goods ServiceOrderWaitFields.FieldsGoods `db:"goods" json:"goods"`
	//订单总的抵扣
	// 例如满减活动，不局限于个别商品的活动
	Exemptions ServiceOrderWaitFields.FieldsExemptions `db:"exemptions" json:"exemptions"`
	//强制订单自动审核开关
	NeedAllowAutoAudit bool
	AllowAutoAudit     bool
	//允许自动配送
	TransportAllowAuto bool `db:"transport_allow_auto" json:"transportAllowAuto" check:"bool" empty:"true"`
	//期望送货时间
	TransportTaskAt time.Time `db:"transport_task_at" json:"transportTaskAt" check:"isoTime" empty:"true"`
	//是否允许货到付款？
	TransportPayAfter bool `db:"transport_pay_after" json:"transportPayAfter" check:"bool" empty:"true"`
	//配送服务系统
	// 0 self 自运营服务; 1 自提; 2 running 跑腿服务; 3 housekeeping 家政服务
	TransportSystem string `db:"transport_system" json:"transportSystem"`
	//费用组成
	PriceList ServiceOrderWaitFields.FieldsPrices `db:"price_list" json:"priceList"`
	//订单总费用
	// 总费用是否支付
	PricePay bool `db:"price_pay" json:"pricePay" check:"bool"`
	//是否需辅助计算抵扣费用
	NeedExPrice bool `json:"needExPrice"`
	// 货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//备注信息
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//日志
	Logs ServiceOrderWaitFields.FieldsLogs `db:"logs" json:"logs"`
	//推荐人手机号
	ReferrerNationCode string `db:"referrer_nation_code" json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func CreateOrder(args *ArgsCreateOrder) (data ServiceOrderWaitFields.FieldsWait, errCode string, err error) {
	//约束参数
	if args.Exemptions == nil {
		args.Exemptions = ServiceOrderWaitFields.FieldsExemptions{}
	}
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//检查TransportSystem
	if args.TransportSystem == "" {
		args.TransportSystem = "self"
	}
	//计算hash值
	hash := CoreFilter.GetSha1Str(fmt.Sprint(args))
	//必须存在货物
	if len(args.Goods) < 1 {
		errCode = "goods_is_empty"
		err = errors.New("goods is empty")
		return
	}
	//获取数据，检查hash是否存在
	if err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_order_wait WHERE hash = $1 AND create_at >= $2", hash, CoreFilter.GetNowTimeCarbon().SubSeconds(10).Time); err == nil {
		if data.ID > 0 {
			errCode = "repeat"
			err = errors.New("repeat hash")
			return
		}
	}
	//检查订单是否存在实体商品
	isFind := false
	for _, v := range args.Goods {
		if v.From.Mark != "virtual" {
			isFind = true
		}
	}
	if isFind {
		//检查收货地址
		if args.AddressTo.Address == "" || args.AddressTo.Phone == "" {
			errCode = "address_to_empty"
			err = errors.New("address to is empty")
			return
		}
	}
	//计算总金额
	var priceTotal int64 = 0
	var price int64 = 0
	//叠加费用组成
	for _, v := range args.PriceList {
		priceTotal += v.Price
		if v.IsPay {
			continue
		}
		price += v.Price
	}
	if args.NeedExPrice && price > 0 {
		for _, v := range args.Exemptions {
			price -= v.Price
		}
		for _, v := range args.Goods {
			for _, v2 := range v.Exemptions {
				price -= v2.Price
			}
		}
	}
	if price < 1 {
		price = 0
	}
	//获取订单审核开关
	var allowAutoAudit bool
	if args.NeedAllowAutoAudit {
		allowAutoAudit = args.AllowAutoAudit
	} else {
		serviceOrderForceAutoAudit := BaseConfig.GetDataBoolNoErr("ServiceOrderForceAutoAudit")
		if !serviceOrderForceAutoAudit {
			allowAutoAudit = OrgCore.Config.GetConfigValBoolNoErr(args.OrgID, "OrderAutoAudit")
		} else {
			allowAutoAudit = serviceOrderForceAutoAudit
		}
	}
	//获取推荐人信息
	var referrerUserID int64
	var referrerUserData UserCore.FieldsUserType
	referrerUserData, err = UserCore.GetUserByPhone(&UserCore.ArgsGetUserByPhone{
		OrgID:      args.OrgID,
		NationCode: args.ReferrerNationCode,
		Phone:      args.ReferrerPhone,
	})
	if err != nil {
		//不记录错误
		err = nil
	} else {
		referrerUserID = referrerUserData.ID
		args.Params = CoreSQLConfig.Set(args.Params, "referrerUserID", fmt.Sprint(referrerUserID))
	}
	//构建请求
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_order_wait", "INSERT INTO service_order_wait (system_mark, org_id, user_id, create_from, hash, address_from, address_to, goods, exemptions, allow_auto_audit, transport_allow_auto, transport_task_at, transport_pay_after, transport_system, price_list, price_pay, currency, price, price_total, des, logs, params, err_code, err_msg) VALUES (:system_mark, :org_id, :user_id, :create_from, :hash, :address_from, :address_to, :goods, :exemptions, :allow_auto_audit, :transport_allow_auto, :transport_task_at, :transport_pay_after, :transport_system, :price_list, :price_pay, :currency, :price, :price_total, :des, :logs, :params, '', '')", map[string]interface{}{
		"system_mark":          args.SystemMark,
		"org_id":               args.OrgID,
		"user_id":              args.UserID,
		"create_from":          args.CreateFrom,
		"hash":                 hash,
		"address_from":         args.AddressFrom,
		"address_to":           args.AddressTo,
		"goods":                args.Goods,
		"exemptions":           args.Exemptions,
		"allow_auto_audit":     allowAutoAudit,
		"transport_allow_auto": args.TransportAllowAuto,
		"transport_task_at":    args.TransportTaskAt,
		"transport_pay_after":  args.TransportPayAfter,
		"transport_system":     args.TransportSystem,
		"price_list":           args.PriceList,
		"price_pay":            args.PricePay,
		"currency":             args.Currency,
		"price":                price,
		"price_total":          priceTotal,
		"des":                  args.Des,
		"logs":                 args.Logs,
		"params":               args.Params,
	}, &data)
	if err != nil {
		errCode = "insert"
		err = errors.New(fmt.Sprint("create order, ", err))
		return
	}
	//推送创建订单请求
	CoreNats.PushDataNoErr("/service/order/create_wait", "", data.ID, "", nil)
	//创建过期请求
	if err = BaseExpireTip.AppendTip(&BaseExpireTip.ArgsAppendTip{
		OrgID:      data.OrgID,
		UserID:     data.UserID,
		SystemMark: "service_order_wait",
		BindID:     data.ID,
		Hash:       "",
		ExpireAt:   CoreFilter.GetNowTimeCarbon().AddDay().Time,
	}); err != nil {
		CoreLog.Error("service order wait, append expire tip, ", err)
		err = nil
	}
	//反馈
	return
}
