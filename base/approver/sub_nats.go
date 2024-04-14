package BaseApprover

import (
	"fmt"
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//模块发出新的审批请求
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "基础审批服务发起新审批",
		Description:  "发起新的审批请求，接收后根据配置构建审批流程",
		EventSubType: "sub",
		Code:         "base_approver_request",
		EventType:    "nats",
		EventURL:     "/base/approver/request",
		EventParams:  "",
	})
	CoreNats.SubDataByteNoErr("base_approver_request", "/base/approver/request", subNatsRequest)
}

// dataSubNatsRequest 请求数据包
type dataSubNatsRequest struct {
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

// moduleCode 模块编码
// approverID 模块ID
// rawData 原始数据
func subNatsRequest(_ *nats.Msg, _ string, approverID int64, moduleCode string, rawData []byte) {
	//前置日志
	appendLog := fmt.Sprint("base approver request: ", moduleCode, " approverID: ", approverID, " data: ", string(rawData))
	//解析参数
	var params dataSubNatsRequest
	err := CoreNats.ReflectDataByte(rawData, &params)
	if err != nil {
		CoreLog.Error(appendLog, ", data parse error: ", err)
		return
	}
}
