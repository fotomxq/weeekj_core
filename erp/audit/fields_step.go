package ERPAudit

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsStep 流程控制
type FieldsStep struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//完成时间
	FinishAt time.Time `db:"finish_at" json:"finishAt"`
	//最终状态
	// 0 无状态; 1 审批通过; 2 拒绝审批
	FinishStatus int `db:"finish_status" json:"finishStatus"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//流程配置
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//流程SN
	// 商户下唯一
	SN int64 `db:"sn" json:"sn" check:"id"`
	//流程名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//创建成员
	CreateOrgBindID int64 `db:"create_org_bind_id" json:"createOrgBindID" check:"id" empty:"true"`
	//可预览的成员ID列
	CanViewOrgBindIDs pq.Int64Array `db:"can_view_org_bind_ids" json:"canViewOrgBindIDs" check:"ids" empty:"true"`
	//可编辑的成员ID列
	CanEditOrgBindIDs pq.Int64Array `db:"can_edit_org_bind_ids" json:"canEditOrgBindIDs" check:"ids" empty:"true"`
	//已经发生或下一步即将发生的成员ID列
	HaveOrgBindIDs pq.Int64Array `db:"have_org_bind_ids" json:"haveOrgBindIDs" check:"ids" empty:"true"`
	//当前节点
	// 对照配置的节点的key
	NowStepChildKey string `db:"now_step_child_key" json:"nowStepChildKey"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
