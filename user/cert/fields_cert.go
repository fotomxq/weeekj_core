package UserCert

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// 用户证件申请记录
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	// 审核前该时间为请求过期时间；审核后该时间为
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//绑定组织
	// 该组织根据资源来源设定
	// 如果是平台资源，则为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//审核用户
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID"`
	//证件ID
	CertID int64 `db:"cert_id" json:"cerID"`
	//申请状态
	// 0 等待审核 / 1 审核通过 / 2 审核失败
	Status int `db:"status" json:"status"`
	//审核备注
	Des string `db:"des" json:"des"`
	//证件内容
	// 该内容和步骤中设置的变量内容一一对应
	Contents CoreSQLConfig.FieldsConfigsType `db:"contents" json:"contents"`
}
