package ERPBudget

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetBudgetFlowList 获取BudgetFlow列表参数
type ArgsGetBudgetFlowList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetBudgetFlowList 获取BudgetFlow列表
func GetBudgetFlowList(args *ArgsGetBudgetFlowList) (dataList []FieldsBudgetFlow, dataCount int64, err error) {
	dataCount, err = budgetFlowDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIntQuery("status", args.Status).SetSearchQuery([]string{"name", "desc"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getBudgetFlowByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetBudgetFlowByID 获取BudgetFlow数据包参数
type ArgsGetBudgetFlowByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetBudgetFlowByID 获取BudgetFlow数
func GetBudgetFlowByID(args *ArgsGetBudgetFlowByID) (data FieldsBudgetFlow, err error) {
	data = getBudgetFlowByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsBudgetFlow{}
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateBudgetFlow 创建BudgetFlow参数
type ArgsCreateBudgetFlow struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"50"`
	//描述
	Desc string `db:"desc" json:"desc" check:"des" min:"1" max:"300" empty:"true"`
	//预算总金额
	Total float64 `db:"total" json:"total" check:"float64Than0"`
	//已使用金额
	Used float64 `db:"used" json:"used" check:"float64Than0"`
	//占用金额
	// 正在使用中，但尚未归档
	Occupied float64 `db:"occupied" json:"occupied" check:"float64Than0"`
}

// CreateBudgetFlow 创建BudgetFlow
func CreateBudgetFlow(args *ArgsCreateBudgetFlow) (id int64, err error) {
	//创建数据
	id, err = budgetFlowDB.Insert().SetFields([]string{"status", "org_id", "name", "desc", "total", "used", "occupied"}).Add(map[string]any{
		"status":   0,
		"org_id":   args.OrgID,
		"name":     args.Name,
		"desc":     args.Desc,
		"total":    args.Total,
		"used":     args.Used,
		"occupied": args.Occupied,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateBudgetFlow 修改BudgetFlow参数
type ArgsUpdateBudgetFlow struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"50"`
	//描述
	Desc string `db:"desc" json:"desc" check:"des" min:"1" max:"300" empty:"true"`
	//预算总金额
	Total float64 `db:"total" json:"total" check:"float64Than0"`
	//已使用金额
	Used float64 `db:"used" json:"used" check:"float64Than0"`
	//占用金额
	// 正在使用中，但尚未归档
	Occupied float64 `db:"occupied" json:"occupied" check:"float64Than0"`
}

// UpdateBudgetFlow 修改BudgetFlow
func UpdateBudgetFlow(args *ArgsUpdateBudgetFlow) (err error) {
	//更新数据
	err = budgetFlowDB.Update().SetFields([]string{"name", "desc", "total", "used", "occupied"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]any{
		"name":     args.Name,
		"desc":     args.Desc,
		"total":    args.Total,
		"used":     args.Used,
		"occupied": args.Occupied,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteBudgetFlowCache(args.ID)
	//反馈
	return
}

// ArgsAuditBudgetFlow 审批BudgetFlow参数
type ArgsAuditBudgetFlow struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
}

// AuditBudgetFlow 审批BudgetFlow
func AuditBudgetFlow(args *ArgsAuditBudgetFlow) (err error) {
	//更新数据
	err = budgetFlowDB.Update().SetFields([]string{"status"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"status": args.Status,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteBudgetFlowCache(args.ID)
	//反馈
	return
}

// ArgsDeleteBudgetFlow 删除BudgetFlow参数
type ArgsDeleteBudgetFlow struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteBudgetFlow 删除BudgetFlow
func DeleteBudgetFlow(args *ArgsDeleteBudgetFlow) (err error) {
	//删除数据
	err = budgetFlowDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteBudgetFlowCache(args.ID)
	//反馈
	return
}

// getBudgetFlowByID 通过ID获取BudgetFlow数据包
func getBudgetFlowByID(id int64) (data FieldsBudgetFlow) {
	cacheMark := getBudgetFlowCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := budgetFlowDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "status", "org_id", "name", "desc", "total", "used", "occupied"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheBudgetFlowTime)
	return
}

// 缓冲
func getBudgetFlowCacheMark(id int64) string {
	return fmt.Sprint("erp:BudgetFlow:id.", id)
}

func deleteBudgetFlowCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getBudgetFlowCacheMark(id))
}
