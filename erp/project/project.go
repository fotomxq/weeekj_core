package ERPProject

import (
	"fmt"
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
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetProjectList 获取Project列表
func GetProjectList(args *ArgsGetProjectList) (dataList []FieldsProject, dataCount int64, err error) {
	dataCount, err = projectDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIntQuery("status", args.Status).SetSearchQuery([]string{"name", "desc"}, args.Search).SelectList("").ResultAndCount(&dataList)
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

// ArgsCreateProject 创建Project参数
type ArgsCreateProject struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//提交人姓名
	SubmitterName string `db:"submitter_name" json:"submitterName" check:"des" min:"1" max:"300"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Desc string `db:"desc" json:"desc" check:"des" min:"1" max:"3000" empty:"true"`
	//预估预算总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
}

// CreateProject 创建Project
func CreateProject(args *ArgsCreateProject) (id int64, err error) {
	//创建数据
	id, err = projectDB.Insert().SetFields([]string{"status", "org_id", "org_bind_id", "user_id", "submitter_name", "approver_id", "approver_name", "name", "desc", "total"}).Add(map[string]any{
		"status":         0,
		"org_id":         args.OrgID,
		"org_bind_id":    args.OrgBindID,
		"user_id":        args.UserID,
		"submitter_name": args.SubmitterName,
		"approver_id":    0,
		"approver_name":  "",
		"name":           args.Name,
		"desc":           args.Desc,
		"total":          args.Total,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateProject 修改Project参数
type ArgsUpdateProject struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Desc string `db:"desc" json:"desc" check:"des" min:"1" max:"3000" empty:"true"`
	//预估预算总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
}

// UpdateProject 修改Project
func UpdateProject(args *ArgsUpdateProject) (err error) {
	//更新数据
	err = projectDB.Update().SetFields([]string{"name", "desc", "total"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).AddWhereUserID(args.UserID).NamedExec(map[string]any{
		"name":  args.Name,
		"desc":  args.Desc,
		"total": args.Total,
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
	err := projectDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "status", "org_id", "org_bind_id", "user_id", "submitter_name", "approver_id", "approver_name", "plan_verifier_id", "plan_verifier_name", "name", "desc", "total"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheProjectTime)
	return
}

// 缓冲
func getProjectCacheMark(id int64) string {
	return fmt.Sprint("erp:Project:id.", id)
}

func deleteProjectCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getProjectCacheMark(id))
}
