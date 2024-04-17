package ERPBudget

import (
	"errors"
	"fmt"
	BaseApproverMod "github.com/fotomxq/weeekj_core/v5/base/approver/mod"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetBudgetList 获取Budget列表参数
type ArgsGetBudgetList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//项目ID
	ProjectID int64 `db:"project_id" json:"projectID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetBudgetList 获取Budget列表
func GetBudgetList(args *ArgsGetBudgetList) (dataList []FieldsBudget, dataCount int64, err error) {
	dataCount, err = budgetDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIDQuery("project_id", args.ProjectID).SetIntQuery("status", args.Status).SetSearchQuery([]string{"name", "des"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getBudgetByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetBudgetByID 获取Budget数据包参数
type ArgsGetBudgetByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetBudgetByID 获取Budget数
func GetBudgetByID(args *ArgsGetBudgetByID) (data FieldsBudget, err error) {
	data = getBudgetByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsBudget{}
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateBudget 创建Budget参数
type ArgsCreateBudget struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"50"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//项目ID
	ProjectID int64 `db:"project_id" json:"projectID" check:"id" empty:"true"`
	//预算总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
	//已使用金额
	Used int64 `db:"used" json:"used" check:"int64Than0"`
	//占用金额
	// 正在使用中，但尚未归档
	Occupied int64 `db:"occupied" json:"occupied" check:"int64Than0"`
	//审批备注
	ApproverRemark string `db:"approver_remark" json:"approverRemark" check:"des" min:"1" max:"300"`
}

// CreateBudget 创建Budget
func CreateBudget(args *ArgsCreateBudget) (id int64, err error) {
	//创建数据
	id, err = budgetDB.Insert().SetFields([]string{"status", "org_id", "name", "des", "project_id", "total", "used", "occupied"}).Add(map[string]any{
		"status":     0,
		"org_id":     args.OrgID,
		"name":       args.Name,
		"des":        args.Des,
		"project_id": args.ProjectID,
		"total":      args.Total,
		"used":       args.Used,
		"occupied":   args.Occupied,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	BaseApproverMod.PushRequest("erp_budget", id, BaseApproverMod.ParamsPushRequest{
		OrgID:          args.OrgID,
		OrgBindID:      args.OrgBindID,
		UserID:         args.UserID,
		ForkCode:       "default",
		ApproverRemark: args.ApproverRemark,
	})
	//反馈
	return
}

// ArgsUpdateBudget 修改Budget参数
type ArgsUpdateBudget struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"50"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//项目ID
	ProjectID int64 `db:"project_id" json:"projectID" check:"id" empty:"true"`
	//预算总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
	//已使用金额
	Used int64 `db:"used" json:"used" check:"int64Than0"`
	//占用金额
	// 正在使用中，但尚未归档
	Occupied int64 `db:"occupied" json:"occupied" check:"int64Than0"`
}

// UpdateBudget 修改Budget
func UpdateBudget(args *ArgsUpdateBudget) (err error) {
	//更新数据
	err = budgetDB.Update().SetFields([]string{"name", "des", "project_id", "total", "used", "occupied"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]any{
		"name":       args.Name,
		"des":        args.Des,
		"project_id": args.ProjectID,
		"total":      args.Total,
		"used":       args.Used,
		"occupied":   args.Occupied,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteBudgetCache(args.ID)
	//反馈
	return
}

// ArgsAuditBudget 审批Budget参数
type ArgsAuditBudget struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
}

// AuditBudget 审批Budget
func AuditBudget(args *ArgsAuditBudget) (err error) {
	//更新数据
	err = budgetDB.Update().SetFields([]string{"status"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"status": args.Status,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteBudgetCache(args.ID)
	//反馈
	return
}

// ArgsDeleteBudget 删除Budget参数
type ArgsDeleteBudget struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteBudget 删除Budget
func DeleteBudget(args *ArgsDeleteBudget) (err error) {
	//删除数据
	err = budgetDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteBudgetCache(args.ID)
	//反馈
	return
}

// getBudgetByID 通过ID获取Budget数据包
func getBudgetByID(id int64) (data FieldsBudget) {
	cacheMark := getBudgetCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := budgetDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "status", "org_id", "name", "des", "project_id", "total", "used", "occupied"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheBudgetTime)
	return
}

// 缓冲
func getBudgetCacheMark(id int64) string {
	return fmt.Sprint("erp:budget:id.", id)
}

func deleteBudgetCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBudgetCacheMark(id))
}
