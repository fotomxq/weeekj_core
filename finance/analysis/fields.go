package FinanceAnalysis

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsAnalysis 总的统计表
type FieldsAnalysis struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计周期
	// 完全相同的一个来源体系，1小时仅构建一条数据
	DayTime time.Time `db:"day_time" json:"dayTime"`
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
