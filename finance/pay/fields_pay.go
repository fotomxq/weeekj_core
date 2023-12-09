package FinancePay

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

type FieldsPayType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//交易过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
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
	//退款的金额
	// 该金额为实际记录、发起的请求金额，累计不能超出总金额
	RefundPrice int64 `db:"refund_price" json:"refundPrice"`
	//是否发送退款请求
	RefundSend bool `db:"refund_send" json:"refundSend"`
	//付款人来源
	// system: user / org
	// id: 用户ID或组织ID
	// mark: 用户OpenID数据
	PaymentCreate CoreSQLFrom.FieldsFrom `db:"payment_create" json:"paymentCreate"`
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝 ; paypal 国际信用卡支付 ;  company_returned 公司赊账付款
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	// company_returned.id 对应公司ID
	PaymentChannel CoreSQLFrom.FieldsFrom `db:"payment_channel" json:"paymentChannel"`
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	// system: 留空则为平台；org
	// id: 组织ID
	PaymentFrom CoreSQLFrom.FieldsFrom `db:"payment_from" json:"paymentFrom"`
	//收款人来源
	// system: user / org
	// id: 用户ID或组织ID
	// mark: 用户OpenID数据
	TakeCreate CoreSQLFrom.FieldsFrom `db:"take_create" json:"takeCreate"`
	//收款渠道
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	TakeChannel CoreSQLFrom.FieldsFrom `db:"take_channel" json:"takeChannel"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	// system: 留空则为平台；org
	// id: 组织ID
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
	//操作人
	// 发起交易的实际人员，可能是后台工作人员为客户发起的交易请求
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//操作原因
	Des string `db:"des" json:"des"`
	//支付失败后的代码
	// 用于系统识别错误类型
	FailedCode string `db:"failed_code" json:"failedCode"`
	//支付失败后的消息
	// 用于用户查看具体的失败原因，可以是财务人员指定的内容
	FailedMessage string `db:"failed_message" json:"failedMessage"`
	//附加参数结构
	// 该参数会混合使用，请注意区分存取的不同类型
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
