package FinanceLog

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

// 资金流水账单
type FieldsLogType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//支付渠道信息ID
	PayID int64 `db:"pay_id" json:"payID"`
	//混淆验证
	Hash string `db:"hash" json:"hash"`
	//交易短key
	// 在历史表中，该值可能发生重复，请勿以该值作为最终唯一判断
	// 用于微信、支付宝等接口对接时，采用的短Key处理机制
	Key string `db:"key" json:"key"`
	//最终状态
	// 0 wait 客户端发起付款，并正在支付中
	// 1 client 客户端完成支付，等待服务端验证
	// 2 failed 交易失败，服务端主动取消交易或其他原因取消交易
	// 3 finish 交易成功
	// 4 remove 交易销毁
	// 5 expire 交易过期
	// 6 refund 发起退款申请
	// 7 refundAudit 退款审核通过，等待处理中
	// 8 refundFailed 退款失败
	// 9 refundFinish 退款完成
	Status int `db:"status" json:"status"`
	//交易货币类型
	// 采用CoreCurrency匹配
	// 86 CNY
	Currency int `db:"currency" json:"currency"`
	//交易金额
	Price int64 `db:"price" json:"price"`
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
	//操作人
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//操作原因
	Des string `db:"des" json:"des"`
}
