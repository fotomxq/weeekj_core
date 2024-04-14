package BaseApprover

import (
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetLogFlows 获取日志行列表参数
type ArgsGetLogFlows struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//日志ID
	LogID int64 `db:"log_id" json:"logID" check:"id" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//审批人用户ID
	// 用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogFlows 获取日志行列表
func GetLogFlows(args *ArgsGetLogFlows) (dataList []FieldsLogFlow, dataCount int64, err error) {
	dataCount, err = configItemDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "flow_order"}).SetPages(args.Pages).SetDeleteQuery("delete_at", false).SetIDQuery("log_id", args.LogID).SetIntQuery("status", args.Status).SetIDQuery("org_bind_id", args.OrgBindID).SetIDQuery("user_id", args.UserID).SetSearchQuery([]string{"approver_remark", "reject_remark"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getLogFlowByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsApproveLogFlow 审核某个节点参数
type ArgsApproveLogFlow struct {
}

// ApproveLogFlow 审核某个节点
func ApproveLogFlow(args *ArgsApproveLogFlow) (err error) {
	//反馈
	return
}

// argsCreateLogFlow 创建审批流参数
type argsCreateLogFlow struct {
	//日志ID
	LogID int64 `db:"log_id" json:"logID" check:"id"`
	//审批顺序
	FlowOrder int `db:"flow_order" json:"flowOrder" check:"intThan0" empty:"true"`
	//审批人ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//审批备注
	ApproverRemark string `db:"approver_remark" json:"approverRemark" check:"des" min:"1" max:"300" empty:"true"`
	//拒绝备注
	RejectRemark string `db:"reject_remark" json:"rejectRemark" check:"des" min:"1" max:"300" empty:"true"`
}

// createLogFlow 创建审批流
func createLogFlow(args *argsCreateLogFlow) (err error) {
	//插入数据
	_, err = logFlowDB.Insert().SetFields([]string{"log_id", "flow_order", "status", "approve_at", "org_bind_id", "user_id", "approver_name", "approver_remark", "reject_remark"}).Add(map[string]any{
		"log_id":          args.LogID,
		"flow_order":      args.FlowOrder,
		"status":          0,
		"approve_at":      time.Time{},
		"org_bind_id":     args.OrgBindID,
		"user_id":         args.UserID,
		"approver_name":   getApproverName(args.OrgBindID, args.UserID),
		"approver_remark": args.ApproverRemark,
		"reject_remark":   args.RejectRemark,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// clearLogItem 清理日志行
func clearLogItem(logID int64) (err error) {
	//获取配置列表
	dataList, _, _ := GetLogFlows(&ArgsGetLogFlows{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  999,
			Sort: "flow_order",
			Desc: false,
		},
		LogID:     logID,
		Status:    -1,
		OrgBindID: -1,
		UserID:    -1,
		Search:    "",
	})
	if len(dataList) < 1 {
		return
	}
	//删除数据
	err = logFlowDB.Delete().NeedSoft(true).SetWhereAnd("log_id", logID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	for _, v := range dataList {
		deleteLogFlowCache(v.ID)
	}
	//反馈
	return
}

// getLogFlowByID 获取审批流程信息
func getLogFlowByID(id int64) (data FieldsLogFlow) {
	cacheMark := getLogFlowCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := logFlowDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "log_id", "flow_order", "status", "approve_at", "org_bind_id", "user_id", "approver_name", "approver_remark", "reject_remark"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheLogFlowTime)
	return
}

// 缓冲
func getLogFlowCacheMark(id int64) string {
	return fmt.Sprint("base:approver:log:flow:id.", id)
}

func deleteLogFlowCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getLogFlowCacheMark(id))
}
