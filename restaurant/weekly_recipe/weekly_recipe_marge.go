package RestaurantWeeklyRecipeMarge

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetWeeklyRecipeList 获取WeeklyRecipe列表参数
type ArgsGetWeeklyRecipeList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分公司ID组
	OrgIDs []int64 `db:"org_ids" json:"orgIDs" check:"ids" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
	//门店ID列
	StoreIDs []int64 `db:"store_ids" json:"storeIDs" check:"ids" empty:"true"`
	//提交组织成员ID
	SubmitOrgBindID int64 `db:"submit_org_bind_id" json:"submitOrgBindID" check:"id" empty:"true"`
	//提交用户ID
	// 与组织ID二选一，如果组织成员ID为空，则使用用户ID；如果组织ID不为空，则使用组织成员ID+用户ID
	SubmitUserID int64 `db:"submit_user_id" json:"submitUserID" check:"id" empty:"true"`
	//审核状态
	// 0 未审核; 1 审核通过; 2 审核不通过
	AuditStatus int `db:"audit_status" json:"auditStatus" check:"intThan0" empty:"true"`
	//审核人ID
	AuditOrgBindID int64 `db:"audit_org_bind_id" json:"auditOrgBindID" check:"id" empty:"true"`
	//审核用户ID
	// 与组织ID二选一，如果组织成员ID为空，则使用用户ID；如果组织ID不为空，则使用组织成员ID+用户ID
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetWeeklyRecipeList 获取WeeklyRecipe列表
func GetWeeklyRecipeList(args *ArgsGetWeeklyRecipeList) (dataList []FieldsWeeklyRecipe, dataCount int64, err error) {
	dataCount, err = weeklyRecipeMargeDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIDsQuery("org_ids", args.OrgIDs).SetIDQuery("store_id", args.StoreID).SetIDsQuery("store_id", args.StoreIDs).SetIDQuery("submit_org_bind_id", args.SubmitOrgBindID).SetIDQuery("submit_user_id", args.SubmitUserID).SetIntQuery("audit_status", args.AuditStatus).SetIDQuery("audit_org_bind_id", args.AuditOrgBindID).SetIDQuery("audit_user_id", args.AuditUserID).SetSearchQuery([]string{"name", "remark"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getWeeklyRecipeSimpleByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetWeeklyRecipeByID 获取WeeklyRecipe数据包参数
type ArgsGetWeeklyRecipeByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
}

// GetWeeklyRecipeByID 获取WeeklyRecipe数
func GetWeeklyRecipeByID(args *ArgsGetWeeklyRecipeByID) (data FieldsWeeklyRecipe, err error) {
	data = getWeeklyRecipeByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.StoreID, data.StoreID) {
		err = errors.New("no data")
		return
	}
	return
}

// GetWeeklyRecipeNameByID 获取菜品名称
func GetWeeklyRecipeNameByID(id int64) (name string) {
	data := getWeeklyRecipeByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

// getWeeklyRecipeByID 通过ID获取WeeklyRecipe数据包
func getWeeklyRecipeByID(id int64) (data FieldsWeeklyRecipe) {
	cacheMark := getWeeklyRecipeCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := weeklyRecipeMargeDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "store_id", "submit_org_bind_id", "submit_user_id", "submit_user_name", "audit_at", "audit_status", "audit_org_bind_id", "audit_user_id", "audit_user_name", "name", "remark", "raw_data"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheWeeklyRecipeTime)
	return
}

// getWeeklyRecipeSimpleByID 通过ID获取WeeklyRecipe数据包
func getWeeklyRecipeSimpleByID(id int64) (data FieldsWeeklyRecipe) {
	cacheMark := getWeeklyRecipeListCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := weeklyRecipeMargeDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "store_id", "submit_org_bind_id", "submit_user_id", "submit_user_name", "audit_at", "audit_status", "audit_org_bind_id", "audit_user_id", "audit_user_name", "name", "remark"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheWeeklyRecipeTime)
	return
}
