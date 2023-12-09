package FinanceAnalysis

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 写入新的统计数据
type ArgsAppendData struct {
	//付款人来源
	PaymentCreate CoreSQLFrom.FieldsFrom `db:"payment_create" json:"paymentCreate"`
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	PaymentChannel CoreSQLFrom.FieldsFrom `db:"payment_channel" json:"paymentChannel"`
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	PaymentFrom CoreSQLFrom.FieldsFrom `db:"payment_from" json:"paymentFrom"`
	//收款人来源
	TakeCreate CoreSQLFrom.FieldsFrom `db:"take_create" json:"takeCreate"`
	//收款渠道
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	TakeChannel CoreSQLFrom.FieldsFrom `db:"take_channel" json:"takeChannel"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
	//交易货币类型
	// 采用CoreCurrency匹配
	// 86 CNY
	Currency int `db:"currency" json:"currency"`
	//交易金额
	Price int64 `db:"price" json:"price"`
}

func AppendData(args *ArgsAppendData) (err error) {
	//获取数据集合
	var paymentCreate string
	paymentCreate, err = args.PaymentCreate.GetRaw()
	if err != nil {
		return
	}
	var paymentChannel string
	paymentChannel, err = args.PaymentChannel.GetRaw()
	if err != nil {
		return
	}
	var paymentFrom string
	paymentFrom, err = args.PaymentFrom.GetRaw()
	if err != nil {
		return
	}
	var takeCreate string
	takeCreate, err = args.TakeCreate.GetRaw()
	if err != nil {
		return
	}
	var takeChannel string
	takeChannel, err = args.TakeChannel.GetRaw()
	if err != nil {
		return
	}
	var takeFrom string
	takeFrom, err = args.TakeFrom.GetRaw()
	if err != nil {
		return
	}
	//检查最近1小时是否存在数据
	type fields struct {
		//ID
		ID int64 `db:"id" json:"id"`
	}
	var data fields
	err = Router2SystemConfig.MainDB.Get(
		&data,
		"SELECT id FROM finance_analysis WHERE day_time > $1 AND currency = $2 AND payment_create @> $3 AND payment_channel @> $4 AND payment_from @> $5 AND take_create @> $6 AND take_channel @> $7 AND take_from @> $8",
		CoreFilter.GetNowTimeCarbon().SubHour().Time,
		args.Currency,
		paymentCreate,
		paymentChannel,
		paymentFrom,
		takeCreate,
		takeChannel,
		takeFrom,
	)
	//如果不存在，则插入数据，否则修改数据
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_analysis SET price = price + :price WHERE id = :id", map[string]interface{}{
			"price": args.Price,
			"id":    data.ID,
		})
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_analysis (payment_create, payment_channel, payment_from, take_create, take_channel, take_from, currency, price) VALUES (:payment_create,:payment_channel,:payment_from,:take_create,:take_channel,:take_from,:currency,:price)", args)
		if err != nil {
			return
		}
	}
	return
}
