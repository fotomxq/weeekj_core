package BaseApprover

import "time"

// FieldsConfig 审批配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//关联的模块标识码
	// erp_project
	ModuleCode string `db:"module_code" json:"moduleCode" check:"des" min:"1" max:"50"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//审批分叉标识码
	// 用于识别模块内，不同的审批流程
	ForkCode string `db:"fork_code" json:"forkCode" check:"des" min:"1" max:"50"`
}
