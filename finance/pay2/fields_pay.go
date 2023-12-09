package FinancePay2

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsPay struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//交易过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//成功时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//失败时间
	FailedAt time.Time `db:"failed_at" json:"failedAt"`
	//交易短hash
	// 用于微信、支付宝等接口对接时，采用的短Key处理机制
	Hash string `db:"hash" json:"hash"`
	//最终状态
	// 0 wait 客户端发起付款，并正在支付中
	// 1 finish 交易成功，服务器验证成功通过
	// 2 failed 交易失败，第三方支付或其他付款形式被关闭
	Status int `db:"status" json:"status"`
	//交易货币类型
	// 采用CoreCurrency匹配
	// 86 CNY
	Currency int `db:"currency" json:"currency"`
	//交易金额
	Price int64 `db:"price" json:"price"`
	//收支款所属商户ID
	// 最终归属权，主要分平台/商户ID
	// 用户所属将被列入平台级别
	PayAndFromOrgID int64 `db:"pay_and_from_org_id" json:"payAndFromOrgID"`
	//付款人来源
	PaymentOrgID  int64 `db:"payment_org_id" json:"paymentOrgID"`
	PaymentUserID int64 `db:"payment_user_id" json:"paymentUserID"`
	//支付渠道
	// cash 现金 / deposit_org_deposit 商户押金 / deposit_org_saving 商户储蓄 / deposit_user_free 用户免费储蓄 / deposit_user_saving 用户储蓄 / deposit_user_deposit 用户押金 / weixin_app 微信APP / weixin_wxx 微信小程序 / weixin_jsapi 微信js / weixin_native 微信二维码 / paypal 海外支付paypal / company_returned 公司预付款回款
	PaymentChannelSystem string `db:"payment_channel_system" json:"paymentChannelSystem"`
	//支付特殊标记
	// 例如用户的openID
	PaymentChannelMark string `db:"payment_channel_mark" json:"paymentChannelMark"`
	//收款人来源
	TakeOrgID  int64 `db:"take_org_id" json:"takeOrgID"`
	TakeUserID int64 `db:"take_user_id" json:"takeUserID"`
	//收款方渠道
	// cash 现金 / deposit_org_deposit 商户押金 / deposit_org_saving 商户储蓄 / deposit_user_free 用户免费储蓄 / deposit_user_saving 用户储蓄 / deposit_user_deposit 用户押金 / weixin_merchant 微信商户转账 / weixin_wxx_user 微信小程序用户
	TakeChannelSystem string `db:"take_channel_system" json:"takeChannelSystem"`
	//收款特殊标记
	// 例如商户ID或用户OpenID
	TakeChannelMark string `db:"take_channel_mark" json:"takeChannelMark"`
	//操作人
	// 发起交易的实际人员，可能是后台工作人员为客户发起的交易请求
	CreateUserID    int64 `db:"create_user_id" json:"createUserID"`
	CreateOrgBindID int64 `db:"create_org_id" json:"createOrgBindID"`
	//操作原因
	Des string `db:"des" json:"des"`
	//附加参数结构
	// 该参数会混合使用，请注意区分存取的不同类型
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
