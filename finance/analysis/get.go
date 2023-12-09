package FinanceAnalysis

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysis 获取平台总统计数据参数
type ArgsGetAnalysis struct {
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType"`
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
	//是否为历史数据
	IsHistory bool `json:"isHistory"`
}

// DataGetAnalysis 获取平台总统计数据结构
type DataGetAnalysis struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//价格合计
	Price int64 `db:"price_count" json:"price"`
}

// GetAnalysis 获取平台总统计数据
func GetAnalysis(args *ArgsGetAnalysis) (dataList []DataGetAnalysis, err error) {
	where := "currency = :currency"
	maps := map[string]interface{}{
		"currency": args.Currency,
	}
	var newWhere string
	newWhere, maps = CoreSQLTime.GetBetweenByTime("day_time", args.TimeBetween, maps)
	if newWhere != "" {
		where = where + " AND " + newWhere
	}
	newWhere, maps, err = args.PaymentCreate.GetList("payment_create", "payment_create", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	newWhere, maps, err = args.PaymentChannel.GetList("payment_channel", "payment_channel", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	newWhere, maps, err = args.PaymentFrom.GetList("payment_from", "payment_from", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	newWhere, maps, err = args.TakeCreate.GetList("take_create", "take_create", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	newWhere, maps, err = args.TakeChannel.GetList("take_channel", "take_channel", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	newWhere, maps, err = args.TakeFrom.GetList("take_from", "take_from", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	tableName := "finance_analysis"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("day_time", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(price) as price_count FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}
