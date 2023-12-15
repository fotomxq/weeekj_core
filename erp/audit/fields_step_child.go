package ERPAudit

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsStepChild 流程节点子步骤
type FieldsStepChild struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//所属流程
	StepID int64 `db:"step_id" json:"stepID"`
	//审批通过时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//拒绝时间
	BanAt time.Time `db:"ban_at" json:"banAt"`
	//过期时间
	// 过期将自动进入驳回程序
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//节点Key
	Key string `db:"key" json:"key" check:"mark"`
	//节点名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//审核模式
	// none 无需审核(记录用); all 必须联合审核; only 只需其中一个审核完成; send 抄送模式
	AuditMode string `db:"audit_mode" json:"auditMode" check:"mark"`
	//审核成员组
	AuditOrgBindGroup int64 `db:"audit_org_bind_group" json:"auditOrgBindGroup" check:"id" empty:"true"`
	//审核指定成员
	AuditOrgBindIDs pq.Int64Array `db:"audit_org_bind_ids" json:"auditOrgBindIDs" check:"ids" empty:"true"`
	//审批角色
	AuditOrgRoleIDs pq.Int64Array `db:"audit_org_role_ids" json:"auditOrgRoleIDs" check:"ids" empty:"true"`
	//等待审核成员
	WaitAuditOrgBindIDs pq.Int64Array `db:"wait_audit_org_bind_ids" json:"waitAuditOrgBindIDs" check:"ids" empty:"true"`
	//已经完成审批成员
	// 已经参与审核的人员，在创建该流程会自动和配置匹配，审核通过后将禁止写入新数据
	FinishAuditOrgBinds pq.Int64Array `db:"finish_audit_org_binds" json:"finishAuditOrgBindIDs"`
	//完成后下一个节点
	// 对照节点key
	// 如果为空，则为最后一个节点处理
	NextStepKey string `db:"next_step_key" json:"nextStepKey"`
	//驳回后下一个节点
	// 驳回必须同时设置完成后节点
	BanNextStepKey string `db:"ban_next_step_key" json:"banNextStepKey"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
