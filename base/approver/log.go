package BaseApprover

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//关联的模块标识码
	// erp_project
	ModuleCode string `db:"module_code" json:"moduleCode" check:"des" min:"1" max:"50" empty:"true"`
	//审批配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	dataCount, err = configItemDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "flow_order"}).SetPages(args.Pages).SetDeleteQuery("delete_at", false).SetIDQuery("org_id", args.OrgID).SetIDQuery("org_bind_id", args.OrgBindID).SetIDQuery("user_id", args.UserID).SetIntQuery("status", args.Status).SetStringQuery("module_code", args.ModuleCode).SetIDQuery("config_id", args.ConfigID).SetSearchQuery([]string{"submitter_name", "approver_remark"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getLogByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetLogByID 获取日志数据包参数
type ArgsGetLogByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetLogByID 获取日志数据包
func GetLogByID(args *ArgsGetLogByID) (data DataLog, err error) {
	rawData := getLogByID(args.ID)
	if rawData.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, rawData.OrgID) {
		err = errors.New("no data")
		return
	}
	rawFlows, _, _ := GetLogFlows(&ArgsGetLogFlows{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  999,
			Sort: "flow_order",
			Desc: false,
		},
		LogID:     rawData.ID,
		Status:    -1,
		OrgBindID: -1,
		UserID:    -1,
		Search:    "",
	})
	var flows DataLogFlows
	if len(rawFlows) > 0 {
		for _, v := range rawFlows {
			flows = append(flows, DataLogFlow{
				FlowOrder:      v.FlowOrder,
				Status:         v.Status,
				CreateAt:       v.CreateAt,
				ApproveAt:      v.ApproveAt,
				OrgBindID:      v.OrgBindID,
				UserID:         v.UserID,
				ApproverName:   v.ApproverName,
				ApproverRemark: v.ApproverRemark,
				RejectRemark:   v.RejectRemark,
			})
		}
	}
	data = DataLog{
		ID:            rawData.ID,
		CreateAt:      rawData.CreateAt,
		UpdateAt:      rawData.UpdateAt,
		DeleteAt:      rawData.DeleteAt,
		OrgID:         rawData.OrgID,
		OrgBindID:     rawData.OrgBindID,
		UserID:        rawData.UserID,
		SubmitterName: rawData.SubmitterName,
		ModuleCode:    rawData.ModuleCode,
		ApproverID:    rawData.ApproverID,
		ConfigID:      rawData.ConfigID,
		Flows:         flows,
	}
	//反馈
	return
}

// ArgsGetLogByModuleAndID 通过模块和ID获取审批参数
type ArgsGetLogByModuleAndID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//关联的模块标识码
	// erp_project
	ModuleCode string `db:"module_code" json:"moduleCode" check:"des" min:"1" max:"50"`
	//审批ID
	ApproverID int64 `db:"approver_id" json:"approverID" check:"id"`
}

// GetLogByModuleAndID 通过模块和ID获取审批
func GetLogByModuleAndID(args *ArgsGetLogByModuleAndID) (data DataLog, err error) {
	//获取日志数据
	var rawLogData FieldsLog
	err = logDB.Get().SetFieldsOne([]string{"id"}).SetDeleteQuery("delete_at", false).SetIDQuery("org_id", args.OrgID).SetIDQuery("org_bind_id", args.OrgBindID).SetIDQuery("user_id", args.UserID).SetStringQuery("module_code", args.ModuleCode).SetIDQuery("approver_id", args.ApproverID).Result(&rawLogData)
	if err != nil {
		return
	}
	if rawLogData.ID > 0 {
		rawLogData = getLogByID(rawLogData.ID)
	}
	//通过ID获取数据包
	data, err = GetLogByID(&ArgsGetLogByID{
		ID:        rawLogData.ID,
		OrgID:     -1,
		OrgBindID: -1,
		UserID:    -1,
	})
	if err != nil {
		rawLogData = FieldsLog{}
		return
	}
	//反馈
	return
}

// ArgsCreateLog 发起新的审批参数
type ArgsCreateLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//关联的模块标识码
	// erp_project
	ModuleCode string `db:"module_code" json:"moduleCode" check:"des" min:"1" max:"50"`
	//审批分叉标识码
	// 用于识别模块内，不同的审批流程
	ForkCode string `db:"fork_code" json:"forkCode" check:"des" min:"1" max:"50"`
	//审批ID
	ApproverID int64 `db:"approver_id" json:"approverID" check:"id"`
	//审批备注
	ApproverRemark string `db:"approver_remark" json:"approverRemark" check:"des" min:"1" max:"300"`
}

// CreateLog 发起新的审批
func CreateLog(args *ArgsCreateLog) (errCode string, err error) {
	//检查是否存在审批的数据
	rawData, _ := GetLogByModuleAndID(&ArgsGetLogByModuleAndID{
		OrgID:      args.OrgID,
		OrgBindID:  -1,
		UserID:     -1,
		ModuleCode: args.ModuleCode,
		ApproverID: args.ApproverID,
	})
	if rawData.ID > 0 {
		errCode = "err_approver_log_repeat"
		err = errors.New("data exists")
		return
	}
	//获取配置数据包
	configData := getConfigByModule(args.OrgID, args.ModuleCode, args.ForkCode)
	if configData.ID < 1 {
		errCode = "err_config"
		err = errors.New("no data")
		return
	}
	//构建数据
	var newLogID int64
	newLogID, err = logDB.Insert().SetFields([]string{"org_id", "org_bind_id", "user_id", "submitter_name", "approver_remark", "status", "module_code", "approver_id", "config_id"}).Add(map[string]any{
		"org_id":          args.OrgID,
		"org_bind_id":     args.OrgBindID,
		"user_id":         args.UserID,
		"submitter_name":  getApproverName(args.OrgBindID, args.UserID),
		"approver_remark": args.ApproverRemark,
		"status":          0,
		"module_code":     configData.ModuleCode,
		"approver_id":     args.ApproverID,
		"config_id":       configData.ID,
	}).ExecAndResultID()
	if err != nil {
		errCode = "err_insert"
		return
	}
	newLogID = newLogID
	//反馈
	return
}

// ArgsDeleteLog 删除审批流参数
type ArgsDeleteLog struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteLog 删除审批流
func DeleteLog(args *ArgsDeleteLog) (err error) {
	//获取日志行数据
	flows, _, _ := GetLogFlows(&ArgsGetLogFlows{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  999,
			Sort: "flow_order",
			Desc: false,
		},
		LogID:     args.ID,
		Status:    -1,
		OrgBindID: -1,
		UserID:    -1,
		Search:    "",
	})
	//删除数据
	err = logDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteLogCache(args.ID)
	//删除行数据
	if len(flows) > 0 {
		for _, v := range flows {
			deleteLogFlowCache(v.ID)
		}
	}
	_ = clearLogFlow(args.ID)
	//反馈
	return
}

// getLogByID 获取审批流程信息
func getLogByID(id int64) (data FieldsLog) {
	cacheMark := getLogCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := logDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "org_bind_id", "user_id", "submitter_name", "approver_remark", "status", "module_code", "approver_id", "config_id"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheLogTime)
	return
}

// 缓冲
func getLogCacheMark(id int64) string {
	return fmt.Sprint("base:approver:log:id.", id)
}

func deleteLogCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(id))
}
