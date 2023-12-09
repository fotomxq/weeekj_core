package FinancePay2

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//对应支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//退款ID
	RefundID int64 `db:"refund_id" json:"refundID"`
	//行为特征
	// pay_create 发起支付; pay_confirm 确认支付; pay_finish 完成支付; pay_failed 支付失败; pay_cancel 关闭支付
	// refund_create 发起退款; refund_confirm 确认退款; refund_finish 完成退款; refund_failed 退款失败; refund_cancel 关闭退款
	Action string `db:"action" json:"action"`
	//编码
	// 大部分情况下属于支付失败后反馈的编码信息
	Code string `db:"code" json:"code"`
	//文本描述
	Des string `db:"des" json:"des"`
	//附加参数结构
	// 该参数会混合使用，请注意区分存取的不同类型
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
