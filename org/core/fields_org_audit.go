package OrgCoreCore

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsOrgAudit 组织申请表
type FieldsOrgAudit struct {
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
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID"`
	//所属用户
	// 掌管该数据的用户，创建人和根管理员，不可删除只能更换
	UserID int64 `db:"user_id" json:"userID"`
	//企业唯一标识码
	// 用于特殊识别和登陆识别等操作
	Key string `db:"key" json:"key"`
	//构架名称，或组织名称
	Name string `db:"name" json:"name"`
	//组织描述
	Des string `db:"des" json:"des"`
	//上级ID
	ParentID int64 `db:"parent_id" json:"parentID"`
	//上级控制权限限制
	ParentFunc pq.StringArray `db:"parent_func" json:"parentFunc"`
	//开通业务
	// 该内容只有总管理员或订阅能进行控制
	OpenFunc pq.StringArray `db:"open_func" json:"openFunc"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
