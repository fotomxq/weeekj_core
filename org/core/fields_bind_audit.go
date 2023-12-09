package OrgCoreCore

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsBindAudit 请求加入组织审核表
type FieldsBindAudit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//审核通过时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//审核拒绝时间
	BanAt time.Time `db:"ban_at" json:"banAt"`
	//拒绝审核原因
	BanDes string `db:"ban_des" json:"banDes"`
	//审核人员
	AuditBindID int64 `db:"audit_bind_id" json:"auditBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//名称
	Name string `db:"name" json:"name"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织分组ID列
	GroupIDs pq.Int64Array `db:"group_ids" json:"groupIDs"`
	//权利主张
	Manager pq.StringArray `db:"manager" json:"manager"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
