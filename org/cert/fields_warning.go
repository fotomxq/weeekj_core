package OrgCert

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsWarning 异常消息
type FieldsWarning struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//处理时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID"`
	//异常证件ID
	CertID int64 `db:"cert_id" json:"certID"`
	//证件标识码
	ConfigMark string `db:"config_mark" json:"configMark"`
	//证件配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//消息
	Msg string `db:"msg" json:"msg"`
	//是否发送消息
	SendMsgAt time.Time `db:"send_msg_at" json:"sendMsgAt"`
	//是否发送短信
	SendSMSAt time.Time `db:"send_sms_at" json:"sendSMSAt"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
