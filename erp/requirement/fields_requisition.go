package ERPRequirement

import "time"

// FieldsRequisition 采购申请单头
type FieldsRequisition struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//审批人ID
	ApproverID int64 `db:"approver_id" json:"approverID" check:"id"`
	//审批人姓名
	ApproverName string `db:"approver_name" json:"approverName" check:"des" min:"1" max:"300"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
	//关联的项目ID
	ProjectID int64 `db:"project_id" json:"projectID" check:"id" empty:"true"`
	//关联项目名称
	ProjectName string `db:"project_name" json:"projectName" check:"des" min:"1" max:"300" empty:"true"`
}
