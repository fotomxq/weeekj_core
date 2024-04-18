package ERPProject

import (
	"errors"
	"fmt"
	BaseApproverMod "github.com/fotomxq/weeekj_core/v5/base/approver/mod"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetProjectList 获取Project列表参数
type ArgsGetProjectList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//验收状态
	// 0: 未验收; 1: 验收中; 2: 验收通过; 3: 验收拒绝
	AcceptanceStatus int `db:"acceptance_status" json:"acceptanceStatus"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetProjectList 获取Project列表
func GetProjectList(args *ArgsGetProjectList) (dataList []FieldsProject, dataCount int64, err error) {
	dataCount, err = projectDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIntQuery("status", args.Status).SetIntQuery("acceptance_status", args.AcceptanceStatus).SetSearchQuery([]string{"name", "des"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getProjectByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetProjectByID 获取Project数据包参数
type ArgsGetProjectByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetProjectByID 获取Project数
func GetProjectByID(args *ArgsGetProjectByID) (data FieldsProject, err error) {
	data = getProjectByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsProject{}
		err = errors.New("no data")
		return
	}
	return
}

// GetProjectNameByID 获取项目名称
func GetProjectNameByID(id int64) string {
	data := getProjectByID(id)
	if data.ID < 1 {
		return ""
	}
	return data.Name
}

// ArgsCreateProject 创建Project参数
type ArgsCreateProject struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//计划验证人ID
	PlanVerifierID int64 `db:"plan_verifier_id" json:"planVerifierID" check:"id" empty:"true"`
	//计划验收人姓名
	PlanVerifierName string `db:"plan_verifier_name" json:"planVerifierName" check:"des" min:"1" max:"300" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//预估预算总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
}

// CreateProject 创建Project
func CreateProject(args *ArgsCreateProject) (id int64, err error) {
	//创建数据
	id, err = projectDB.Insert().SetFields([]string{"status", "acceptance_status", "org_id", "plan_verifier_id", "plan_verifier_name", "name", "des", "total"}).Add(map[string]any{
		"status":             0,
		"acceptance_status":  0,
		"org_id":             args.OrgID,
		"plan_verifier_id":   args.PlanVerifierID,
		"plan_verifier_name": args.PlanVerifierName,
		"name":               args.Name,
		"des":                args.Des,
		"total":              args.Total,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//nats 通知审批
	BaseApproverMod.PushRequest("erp_project", id, BaseApproverMod.ParamsPushRequest{
		OrgID:          args.OrgID,
		OrgBindID:      0,
		UserID:         0,
		ForkCode:       "default",
		ApproverRemark: args.Des,
	})
	//反馈
	return
}

// ArgsUpdateProject 修改Project参数
type ArgsUpdateProject struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//计划验证人ID
	PlanVerifierID int64 `db:"plan_verifier_id" json:"planVerifierID" check:"id" empty:"true"`
	//计划验收人姓名
	PlanVerifierName string `db:"plan_verifier_name" json:"planVerifierName" check:"des" min:"1" max:"300" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//预估预算总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
}

// UpdateProject 修改Project
func UpdateProject(args *ArgsUpdateProject) (err error) {
	//更新数据
	err = projectDB.Update().SetFields([]string{"plan_verifier_id", "plan_verifier_name", "name", "des", "total"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]any{
		"plan_verifier_id":   args.PlanVerifierID,
		"plan_verifier_name": args.PlanVerifierName,
		"name":               args.Name,
		"des":                args.Des,
		"total":              args.Total,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteProjectCache(args.ID)
	//反馈
	return
}

// ArgsAuditProject 审批Project参数
type ArgsAuditProject struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
}

// AuditProject 审批Project
func AuditProject(args *ArgsAuditProject) (err error) {
	//更新数据
	err = projectDB.Update().SetFields([]string{"status"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"status": args.Status,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteProjectCache(args.ID)
	//反馈
	return
}

// ArgsAcceptanceProject 验收Project参数
type ArgsAcceptanceProject struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//验收状态
	// 0: 未验收; 1: 验收中; 2: 验收通过; 3: 验收拒绝
	AcceptanceStatus int `db:"acceptance_status" json:"acceptanceStatus"`
}

// AcceptanceProject 验收Project
func AcceptanceProject(args *ArgsAcceptanceProject) (err error) {
	//更新数据
	err = projectDB.Update().SetFields([]string{"acceptance_status"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"acceptance_status": args.AcceptanceStatus,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteProjectCache(args.ID)
	//反馈
	return
}

// ArgsDeleteProject 删除Project参数
type ArgsDeleteProject struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteProject 删除Project
func DeleteProject(args *ArgsDeleteProject) (err error) {
	//删除数据
	err = projectDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteProjectCache(args.ID)
	//反馈
	return
}

// getProjectByID 通过ID获取Project数据包
func getProjectByID(id int64) (data FieldsProject) {
	cacheMark := getProjectCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := projectDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "status", "acceptance_status", "org_id", "plan_verifier_id", "plan_verifier_name", "name", "des", "total"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheProjectTime)
	return
}

// 缓冲
func getProjectCacheMark(id int64) string {
	return fmt.Sprint("erp:project:id.", id)
}

func deleteProjectCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getProjectCacheMark(id))
}
