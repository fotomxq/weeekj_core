package RestaurantWeeklyRecipeMarge

import (
	"errors"
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	RestaurantRawMaterials "github.com/fotomxq/weeekj_core/v5/restaurant/raw_materials"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsGetWeeklyRecipeRaw struct {
	//每周菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id" index:"true"`
}

// GetWeeklyRecipeRaw 获取指定周菜谱的原材料
func GetWeeklyRecipeRaw(args *ArgsGetWeeklyRecipeRaw) (dataList []FieldsWeeklyRecipeRaw, err error) {
	cacheMark := getWeeklyRecipeRawCacheMark(args.WeeklyRecipeID)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	err = weeklyRecipeRawDB.Select().SetFieldsList([]string{"id", "create_at", "update_at", "delete_at", "org_id", "store_id", "weekly_recipe_id", "dining_date", "day_type", "recipe_id", "recipe_name", "material_id", "material_name", "use_count"}).SetIDQuery("weekly_recipe_id", args.WeeklyRecipeID).SetDeleteQuery("delete_at", false).SelectList("").Result(&dataList)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, CoreCache.CacheTime1Hour)
	return
}

type ArgsSetWeeklyRecipeRaw struct {
	//每周菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id" index:"true"`
	//数据
	Items []ArgsSetWeeklyRecipeRawItem `db:"items" json:"items"`
}

type ArgsSetWeeklyRecipeRawItem struct {
	// 用餐日期
	// 例如：20210101
	DiningDate int `db:"dining_date" json:"diningDate" index:"true"`
	//每日类型
	// 1:早餐; 2:中餐; 3:晚餐
	DayType int `db:"day_type" json:"dayType" check:"intThan0" empty:"true" index:"true"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id" index:"true"`
	//原材料ID
	MaterialID int64 `db:"material_id" json:"materialID" check:"id" empty:"true" index:"true"`
	//用量
	UseCount float64 `db:"use_count" json:"useCount" check:"intThan0"`
}

func SetWeeklyRecipeRaw(args *ArgsSetWeeklyRecipeRaw) (err error) {
	//获取数据
	recipeData := getWeeklyRecipeByID(args.WeeklyRecipeID)
	if recipeData.ID < 1 {
		err = errors.New("no data")
		return
	}
	if recipeData.AuditStatus != 1 {
		err = errors.New("audit status error")
		return
	}
	//检查是否存在数据
	var rawList []FieldsWeeklyRecipeRaw
	rawList, err = GetWeeklyRecipeRaw(&ArgsGetWeeklyRecipeRaw{
		WeeklyRecipeID: args.WeeklyRecipeID,
	})
	if err == nil && len(rawList) > 0 {
		//清理并重建数据
		err = weeklyRecipeRawDB.Delete().NeedSoft(true).SetWhereOrThan("weekly_recipe_id", args.WeeklyRecipeID).ExecNamed(nil)
		if err != nil {
			return
		}
	} else {
		err = nil
	}
	//创建数据
	for k := 0; k < len(args.Items); k++ {
		v := args.Items[k]
		err = weeklyRecipeRawDB.Insert().SetFields([]string{"org_id", "store_id", "weekly_recipe_id", "dining_date", "day_type", "recipe_id", "recipe_name", "material_id", "material_name", "use_count"}).Add(map[string]any{
			"org_id":           recipeData.OrgID,
			"store_id":         recipeData.StoreID,
			"weekly_recipe_id": args.WeeklyRecipeID,
			"dining_date":      v.DiningDate,
			"day_type":         v.DayType,
			"recipe_id":        v.RecipeID,
			"recipe_name":      GetWeeklyRecipeNameByID(v.RecipeID),
			"material_id":      v.MaterialID,
			"material_name":    RestaurantRawMaterials.GetRawNameByID(v.MaterialID),
			"use_count":        v.UseCount,
		}).ExecAndCheckID()
		if err != nil {
			return
		}
	}
	//删除缓冲
	deleteWeeklyRecipeRawCache(args.WeeklyRecipeID)
	//反馈
	return
}

// 缓冲
func getWeeklyRecipeRawCacheMark(weeklyRecipeID int64) string {
	return fmt.Sprint("restaurant:weekly.recipe:raw:id.", weeklyRecipeID)
}

func deleteWeeklyRecipeRawCache(weeklyRecipeID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWeeklyRecipeRawCacheMark(weeklyRecipeID))
}
