package RestaurantWeeklyRecipe

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetWeeklyRecipeItemList 获取WeeklyRecipeItem列表参数
type ArgsGetWeeklyRecipeItemList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id" empty:"true"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
	//菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id" empty:"true"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetWeeklyRecipeItemList 获取WeeklyRecipeItem列表
func GetWeeklyRecipeItemList(args *ArgsGetWeeklyRecipeItemList) (dataList []FieldsWeeklyRecipeItem, dataCount int64, err error) {
	dataCount, err = weeklyRecipeItemDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("raw_org_id", args.RawOrgID).SetIDQuery("org_id", args.OrgID).SetIDQuery("store_id", args.StoreID).SetIDQuery("weekly_recipe_id", args.WeeklyRecipeID).SetIDQuery("recipe_id", args.RecipeID).SetSearchQuery([]string{"name"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getWeeklyRecipeItemByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetWeeklyRecipeItemByID 获取WeeklyRecipeItem数据包参数
type ArgsGetWeeklyRecipeItemByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id" empty:"true"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
}

// GetWeeklyRecipeItemByID 获取WeeklyRecipeItem数
func GetWeeklyRecipeItemByID(args *ArgsGetWeeklyRecipeItemByID) (data FieldsWeeklyRecipeItem, err error) {
	data = getWeeklyRecipeItemByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.RawOrgID, data.RawOrgID) || !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.StoreID, data.StoreID) {
		err = errors.New("no data")
		return
	}
	return
}

// GetWeeklyRecipeItemNameByID 获取菜品名称
func GetWeeklyRecipeItemNameByID(id int64) (name string) {
	data := getWeeklyRecipeItemByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

// ArgsCreateWeeklyRecipeItem 创建WeeklyRecipeItem参数
type ArgsCreateWeeklyRecipeItem struct {
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
}

// CreateWeeklyRecipeItem 创建WeeklyRecipeItem
func CreateWeeklyRecipeItem(args *ArgsCreateWeeklyRecipeItem) (id int64, err error) {
	//创建数据
	id, err = weeklyRecipeItemDB.Insert().SetFields([]string{"raw_org_id", "org_id", "store_id", "weekly_recipe_id", "recipe_id", "name", "price"}).Add(map[string]any{
		"raw_org_id":       args.RawOrgID,
		"org_id":           args.OrgID,
		"store_id":         args.StoreID,
		"weekly_recipe_id": args.WeeklyRecipeID,
		"recipe_id":        args.RecipeID,
		"name":             args.Name,
		"price":            args.Price,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateWeeklyRecipeItem 修改WeeklyRecipeItem参数
type ArgsUpdateWeeklyRecipeItem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	//菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
}

// UpdateWeeklyRecipeItem 修改WeeklyRecipeItem
func UpdateWeeklyRecipeItem(args *ArgsUpdateWeeklyRecipeItem) (err error) {
	//更新数据
	err = weeklyRecipeItemDB.Update().SetFields([]string{"raw_org_id", "org_id", "store_id", "weekly_recipe_id", "recipe_id", "name", "price"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"raw_org_id":       args.RawOrgID,
		"org_id":           args.OrgID,
		"store_id":         args.StoreID,
		"weekly_recipe_id": args.WeeklyRecipeID,
		"recipe_id":        args.RecipeID,
		"name":             args.Name,
		"price":            args.Price,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteWeeklyRecipeItemCache(args.ID)
	//反馈
	return
}

// ArgsDeleteWeeklyRecipeItem 删除WeeklyRecipeItem参数
type ArgsDeleteWeeklyRecipeItem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteWeeklyRecipeItem 删除WeeklyRecipeItem
func DeleteWeeklyRecipeItem(args *ArgsDeleteWeeklyRecipeItem) (err error) {
	//删除数据
	err = weeklyRecipeItemDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteWeeklyRecipeItemCache(args.ID)
	//反馈
	return
}

// getWeeklyRecipeItemByID 通过ID获取WeeklyRecipeItem数据包
func getWeeklyRecipeItemByID(id int64) (data FieldsWeeklyRecipeItem) {
	cacheMark := getWeeklyRecipeItemCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := weeklyRecipeItemDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "raw_org_id", "org_id", "store_id", "weekly_recipe_id", "recipe_id", "name", "price"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheWeeklyRecipeItemTime)
	return
}

// 缓冲
func getWeeklyRecipeItemCacheMark(id int64) string {
	return fmt.Sprint("restaurant:weekly_recipe_item:id.", id)
}

func deleteWeeklyRecipeItemCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWeeklyRecipeItemCacheMark(id))
}
