package BaseApproverMod

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

// ParamsPushRequest 请求数据包
type ParamsPushRequest struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//审批分叉标识码
	// 用于识别模块内，不同的审批流程
	ForkCode string `db:"fork_code" json:"forkCode" check:"des" min:"1" max:"50"`
	//审批备注
	ApproverRemark string `db:"approver_remark" json:"approverRemark" check:"des" min:"1" max:"300"`
}

// PushRequest 请求创建新的审批
func PushRequest(moduleCode string, approverID int64, params ParamsPushRequest) {
	CoreNats.PushDataNoErr("base_approver_request", "/base/approver/request", "", approverID, moduleCode, params)
}
