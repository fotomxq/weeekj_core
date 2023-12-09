package ERPAudit

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	ERPCore "gitee.com/weeekj/weeekj_core/v5/erp/core"
	"github.com/lib/pq"
	"time"
)

// FieldsConfig 流程配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//发布状态
	PublishAt time.Time `db:"publish_at" json:"publishAt"`
	//hash
	// 如果hash和提交hash不同，服务端将自动拒绝更新，避免流处理异常
	Hash string `db:"hash" json:"hash" check:"sha1"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//流程名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//节点设置
	// 流程化的核心处理
	StepList FieldsConfigStepList `db:"step_list" json:"stepList"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// FieldsConfigStepList 流程节点列
type FieldsConfigStepList []FieldsConfigStep

// Value sql底层处理器
func (t FieldsConfigStepList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigStepList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// Len 排序支持
func (t FieldsConfigStepList) Len() int {
	return len(t)
}
func (t FieldsConfigStepList) Less(i, j int) bool {
	return t[i].Sort < t[j].Sort
}

func (t FieldsConfigStepList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// FieldsConfigStep 流程节点配置
type FieldsConfigStep struct {
	//节点顺序
	Sort int `db:"sort" json:"sort"`
	//节点key
	Key string `db:"key" json:"key" check:"mark"`
	//节点名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//过期时间
	// 单位: 秒
	ExpireSec int `db:"expire_sec" json:"expireSec"`
	//审核模式
	// none 无需审核(记录用); all 必须联合审核; only 只需其中一个审核完成; send 抄送模式
	AuditMode string `db:"audit_mode" json:"auditMode" check:"mark"`
	//审核成员组
	AuditOrgBindGroup int64 `db:"audit_org_bind_group" json:"auditOrgBindGroup" check:"id" empty:"true"`
	//审核指定成员
	AuditOrgBindIDs pq.Int64Array `db:"audit_org_bind_ids" json:"auditOrgBindIDs" check:"ids" empty:"true"`
	//审批角色
	AuditOrgRoleIDs pq.Int64Array `db:"audit_org_role_ids" json:"auditOrgRoleIDs" check:"ids" empty:"true"`
	//节点组件
	ComponentList ERPCore.FieldsComponentDefineList `db:"component_list" json:"componentList"`
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

// Value sql底层处理器
func (t FieldsConfigStep) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigStep) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
