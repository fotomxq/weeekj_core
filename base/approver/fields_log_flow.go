package BaseApprover

import "time"

// FieldsLogFlow 审批
type FieldsLogFlow struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//日志ID
	LogID int64 `db:"log_id" json:"logID" check:"id"`
	//审批顺序
	FlowOrder int `db:"flow_order" json:"flowOrder" check:"intThan0" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//审批时间
	// 审批完成或拒绝时间
	ApproveAt time.Time `db:"approve_at" json:"approveAt"`
	//审批人ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//审批人姓名
	ApproverName string `db:"approver_name" json:"approverName" check:"des" min:"1" max:"300"`
	//审批备注
	ApproverRemark string `db:"approver_remark" json:"approverRemark" check:"des" min:"1" max:"300" empty:"true"`
	//拒绝备注
	RejectRemark string `db:"reject_remark" json:"rejectRemark" check:"des" min:"1" max:"300" empty:"true"`
}
