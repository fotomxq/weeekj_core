package FinancePay2

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsRefund struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
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
	// 0 refund 发起退款申请
	// 1 refundAudit 退款审核通过，等待处理中
	// 2 refundFailed 退款失败
	// 3 refundFinish 退款完成
	Status int `db:"status" json:"status"`
	//对应支付ID
	// 每一个支付只能发起10条退款请求
	PayID int64 `db:"pay_id" json:"payID"`
	//操作人
	// 发起交易的实际人员，可能是后台工作人员为客户发起的交易请求
	CreateUserID    int64 `db:"create_user_id" json:"createUserID"`
	CreateOrgBindID int64 `db:"create_org_id" json:"createOrgBindID"`
	//退款的金额
	// 该金额为实际记录、发起的请求金额，累计不能超出支付的金额
	RefundPrice int64 `db:"refund_price" json:"refundPrice"`
	//是否发送退款请求
	RefundSend bool `db:"refund_send" json:"refundSend"`
	//操作原因
	Des string `db:"des" json:"des"`
	//附加参数结构
	// 该参数会混合使用，请注意区分存取的不同类型
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
