package RestaurantRecipe

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetRecipeList 获取Recipe列表参数
type ArgsGetRecipeList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRecipeList 获取Recipe列表
func GetRecipeList(args *ArgsGetRecipeList) (dataList []FieldsRecipe, dataCount int64, err error) {
	dataCount, err = recipeDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("category_id", args.CategoryID).SetIDQuery("org_id", args.OrgID).SetIDQuery("store_id", args.StoreID).SetSearchQuery([]string{"name"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getRecipeByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetRecipeByID 获取Recipe数据包参数
type ArgsGetRecipeByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id" empty:"true"`
}

// GetRecipeByID 获取Recipe数
func GetRecipeByID(args *ArgsGetRecipeByID) (data FieldsRecipe, err error) {
	data = getRecipeByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.StoreID, data.StoreID) {
		err = errors.New("no data")
		return
	}
	return
}

// GetRecipeNameByID 获取菜品名称
func GetRecipeNameByID(id int64) (name string) {
	data := getRecipeByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

// ArgsCreateRecipe 创建Recipe参数
type ArgsCreateRecipe struct {
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//单位
	Unit string `db:"unit" json:"unit" check:"des" min:"1" max:"60" empty:"true"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	// 建议售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"0" max:"3000" empty:"true"`
}

// CreateRecipe 创建Recipe
func CreateRecipe(args *ArgsCreateRecipe) (id int64, err error) {
	//创建数据
	id, err = recipeDB.Insert().SetFields([]string{"category_id", "name", "unit", "org_id", "store_id", "price", "remark"}).Add(map[string]any{
		"category_id": args.CategoryID,
		"name":        args.Name,
		"unit":        args.Unit,
		"org_id":      args.OrgID,
		"store_id":    args.StoreID,
		"price":       args.Price,
		"remark":      args.Remark,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateRecipe 修改Recipe参数
type ArgsUpdateRecipe struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//菜品名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//单位
	Unit string `db:"unit" json:"unit" check:"des" min:"1" max:"60" empty:"true"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID" check:"id"`
	// 建议售价
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"0" max:"3000" empty:"true"`
}

// UpdateRecipe 修改Recipe
func UpdateRecipe(args *ArgsUpdateRecipe) (err error) {
	//更新数据
	err = recipeDB.Update().SetFields([]string{"category_id", "name", "unit", "org_id", "store_id", "price", "remark"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"category_id": args.CategoryID,
		"name":        args.Name,
		"unit":        args.Unit,
		"org_id":      args.OrgID,
		"store_id":    args.StoreID,
		"price":       args.Price,
		"remark":      args.Remark,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteRecipeCache(args.ID)
	//反馈
	return
}

// ArgsDeleteRecipe 删除Recipe参数
type ArgsDeleteRecipe struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteRecipe 删除Recipe
func DeleteRecipe(args *ArgsDeleteRecipe) (err error) {
	//删除数据
	err = recipeDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteRecipeCache(args.ID)
	//反馈
	return
}

// getRecipeByID 通过ID获取Recipe数据包
func getRecipeByID(id int64) (data FieldsRecipe) {
	cacheMark := getRecipeCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := recipeDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "category_id", "name", "unit", "org_id", "store_id", "price", "remark"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheRecipeTime)
	return
}

// 缓冲
func getRecipeCacheMark(id int64) string {
	return fmt.Sprint("restaurant:recipe:recipe:id.", id)
}

func deleteRecipeCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getRecipeCacheMark(id))
}
