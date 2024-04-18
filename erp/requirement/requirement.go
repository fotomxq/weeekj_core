package ERPRequirement

import (
	"errors"
	"fmt"
	BaseApproverMod "github.com/fotomxq/weeekj_core/v5/base/approver/mod"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetRequirementList 获取Requirement列表参数
type ArgsGetRequirementList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联的项目ID
	ProjectID int64 `db:"project_id" json:"projectID" check:"id" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRequirementList 获取Requirement列表
func GetRequirementList(args *ArgsGetRequirementList) (dataList []FieldsRequisition, dataCount int64, err error) {
	dataCount, err = requirementDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIDQuery("project_id", args.ProjectID).SetIntQuery("status", args.Status).SetSearchQuery([]string{"project_name", "remark"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getRequirementByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetRequirementByID 获取Requirement数据包参数
type ArgsGetRequirementByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetRequirementByID 获取Requirement数
func GetRequirementByID(args *ArgsGetRequirementByID) (data FieldsRequisition, err error) {
	data = getRequirementByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsRequisition{}
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateRequirement 创建Requirement参数
type ArgsCreateRequirement struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
	//关联的项目ID
	ProjectID int64 `db:"project_id" json:"projectID" check:"id" empty:"true"`
	//关联项目名称
	ProjectName string `db:"project_name" json:"projectName" check:"des" min:"1" max:"300" empty:"true"`
}

// CreateRequirement 创建Requirement
func CreateRequirement(args *ArgsCreateRequirement) (id int64, err error) {
	//创建数据
	id, err = requirementDB.Insert().SetFields([]string{"status", "org_id", "org_bind_id", "remark", "project_id", "project_name"}).Add(map[string]any{
		"status":       0,
		"org_id":       args.OrgID,
		"org_bind_id":  args.OrgBindID,
		"remark":       args.Remark,
		"project_id":   args.ProjectID,
		"project_name": args.ProjectName,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//nats 通知审批
	BaseApproverMod.PushRequest("erp_requirement", id, BaseApproverMod.ParamsPushRequest{
		OrgID:          args.OrgID,
		OrgBindID:      args.OrgBindID,
		UserID:         0,
		ForkCode:       "default",
		ApproverRemark: args.Remark,
	})
	//反馈
	return
}

// ArgsUpdateRequirement 修改Requirement参数
type ArgsUpdateRequirement struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
	//关联的项目ID
	ProjectID int64 `db:"project_id" json:"projectID" check:"id" empty:"true"`
	//关联项目名称
	ProjectName string `db:"project_name" json:"projectName" check:"des" min:"1" max:"300" empty:"true"`
}

// UpdateRequirement 修改Requirement
func UpdateRequirement(args *ArgsUpdateRequirement) (err error) {
	//更新数据
	err = requirementDB.Update().SetFields([]string{"org_bind_id", "remark", "project_id", "project_name"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]any{
		"org_bind_id":  args.OrgBindID,
		"remark":       args.Remark,
		"project_id":   args.ProjectID,
		"project_name": args.ProjectName,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteRequirementCache(args.ID)
	//反馈
	return
}

// ArgsAuditRequirement 审批Requirement参数
type ArgsAuditRequirement struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
}

// AuditRequirement 审批Requirement
func AuditRequirement(args *ArgsAuditRequirement) (err error) {
	//更新数据
	err = requirementDB.Update().SetFields([]string{"status"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"status": args.Status,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteRequirementCache(args.ID)
	//反馈
	return
}

// ArgsDeleteRequirement 删除Requirement参数
type ArgsDeleteRequirement struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteRequirement 删除Requirement
func DeleteRequirement(args *ArgsDeleteRequirement) (err error) {
	//删除数据
	err = requirementDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteRequirementCache(args.ID)
	//反馈
	return
}

// getRequirementByID 通过ID获取Requirement数据包
func getRequirementByID(id int64) (data FieldsRequisition) {
	cacheMark := getRequirementCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := requirementDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "status", "org_id", "org_bind_id", "remark", "project_id", "project_name"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheRequirementTime)
	return
}

// 缓冲
func getRequirementCacheMark(id int64) string {
	return fmt.Sprint("erp:requirement:id.", id)
}

func deleteRequirementCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getRequirementCacheMark(id))
}
