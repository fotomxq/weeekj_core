package BaseApprover

import (
	"time"
)

// FieldsLog 审批日志
type FieldsLog struct {
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
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//提交人姓名
	SubmitterName string `db:"submitter_name" json:"submitterName" check:"des" min:"1" max:"300"`
	//审批备注
	ApproverRemark string `db:"approver_remark" json:"approverRemark" check:"des" min:"1" max:"300"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//关联的模块标识码
	// erp_project
	ModuleCode string `db:"module_code" json:"moduleCode" check:"des" min:"1" max:"50"`
	//审批ID
	ApproverID int64 `db:"approver_id" json:"approverID" check:"id"`
	//审批配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
}
