package FinanceSafe

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsSafeType 安全事件
type FieldsSafeType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//造成该事件的来源
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//交易发生双方
	//来源，支付方
	PaymentCreate CoreSQLFrom.FieldsFrom `db:"payment_create" json:"paymentCreate"`
	//支付方的来源
	// 留空则代表平台方，否则为商户或加盟商
	PaymentFrom CoreSQLFrom.FieldsFrom `db:"payment_from" json:"paymentFrom"`
	//目标，接收方
	TakeCreate CoreSQLFrom.FieldsFrom `db:"take_create" json:"takeCreate"`
	//收款方来源
	// 留空则代表平台方，否则为商户或加盟商
	TakeFrom CoreSQLFrom.FieldsFrom `db:"take_from" json:"takeFrom"`
	//支付ID
	PayID string `db:"pay_id" json:"payID"`
	//日志ID
	PayLogID string `db:"pay_log_id" json:"payLogID"`
	//安全事件详细描述信息
	Message string `db:"message" json:"message"`
	//安全标识码
	// 用于其他语言翻译或前端传输
	Code string `db:"code" json:"code"`
	//是否需要发出预警消息
	NeedEW bool `db:"need_ew" json:"needEW"`
	//是否已经发出了预警
	AllowEW bool `db:"allow_ew" json:"allowEW"`
	//预警消息模版Mark
	EWTemplateMark string `db:"ew_template_mark" json:"ewTemplateMark"`
	//是否打开
	// 处理完成后将标记为false
	AllowOpen bool `db:"allow_open" json:"allowOpen"`
}
