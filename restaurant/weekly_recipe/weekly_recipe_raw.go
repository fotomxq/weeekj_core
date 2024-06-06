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
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id" empty:"true" index:"true"`
	//菜谱类型ID
	RecipeTypeID int64 `db:"recipe_type_id" json:"recipeTypeID" check:"id" empty:"true" index:"true"`
	// 用餐日期
	// 例如：20210101
	DiningDate int `db:"dining_date" json:"diningDate" empty:"true" index:"true"`
	//每日类型
	// 1:早餐; 2:中餐; 3:晚餐
	DayType int `db:"day_type" json:"dayType" check:"intThan0" empty:"true" index:"true"`
	//菜品ID
	RecipeID int64 `db:"recipe_id" json:"recipeID" check:"id" empty:"true" index:"true"`
}

// GetWeeklyRecipeRaw 获取指定周菜谱的原材料
func GetWeeklyRecipeRaw(args *ArgsGetWeeklyRecipeRaw) (dataList []FieldsWeeklyRecipeRaw, err error) {
	err = weeklyRecipeRawDB.Select().SetIDQuery("weekly_recipe_id", args.WeeklyRecipeID).SetIDQuery("recipe_type_id", args.RecipeTypeID).SetIntQuery("dining_date", args.DiningDate).SetIntQuery("day_type", args.DayType).SetIDQuery("recipe_id", args.RecipeID).SetDeleteQuery("delete_at", false).SelectList("").Result(&dataList)
	if err != nil {
		return
	}
	for k, v := range dataList {
		dataList[k] = GetWeeklyRecipeRawData(v.ID)
	}
	return
}

// ArgsGetWeeklyRecipeRawByChildID 获取指定周菜谱的原材料参数
type ArgsGetWeeklyRecipeRawByChildID struct {
	//周菜品关联行ID
	RecipeChildID int64 `db:"recipe_child_id" json:"recipeChildID" check:"id" index:"true"`
}

// GetWeeklyRecipeRawByChildID 获取指定周菜谱的原材料
func GetWeeklyRecipeRawByChildID(args *ArgsGetWeeklyRecipeRawByChildID) (dataList []FieldsWeeklyRecipeRaw, err error) {
	err = weeklyRecipeRawDB.Select().SetIDQuery("recipe_child_id", args.RecipeChildID).SelectList("").Result(&dataList)
	if err != nil {
		return
	}
	for k, v := range dataList {
		dataList[k] = GetWeeklyRecipeRawData(v.ID)
	}
	return
}

type ArgsSetWeeklyRecipeRaw struct {
	//周菜品关联行ID
	RecipeChildID int64 `db:"recipe_child_id" json:"recipeChildID" check:"id" index:"true"`
	//原材料组成
	RawList []ArgsSetWeeklyRecipeRawItem `db:"raw_list" json:"rawList"`
}

type ArgsSetWeeklyRecipeRawItem struct {
	//原材料ID
	MaterialID int64 `db:"material_id" json:"materialID" check:"id" empty:"true" index:"true"`
	//用量
	UseCount float64 `db:"use_count" json:"useCount" check:"intThan0"`
	//单价
	Price float64 `db:"price" json:"price" check:"intThan0" empty:"true"`
	//总价
	TotalPrice float64 `db:"total_price" json:"totalPrice" check:"intThan0" empty:"true"`
}

func SetWeeklyRecipeRaw(args *ArgsSetWeeklyRecipeRaw) (err error) {
	//查到关联行
	recipeRawData := getWeeklyRecipeChildByID(args.RecipeChildID)
	if recipeRawData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//查到关联日
	recipeDayData := getWeeklyRecipeDayByID(recipeRawData.WeeklyRecipeDayID)
	if recipeDayData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//获取数据
	recipeData := getWeeklyRecipeByID(recipeRawData.WeeklyRecipeID)
	if recipeData.ID < 1 {
		err = errors.New("no data")
		return
	}
	if recipeData.AuditStatus != 1 {
		err = errors.New("audit status error")
		return
	}
	//检查是否存在数据
	rawCount, _ := weeklyRecipeRawDB.Select().SetIDQuery("recipe_child_id", recipeRawData.ID).ResultCount()
	if rawCount > 0 {
		//清理并重建数据
		err = weeklyRecipeRawDB.Delete().NeedSoft(true).SetWhereOrThan("recipe_child_id", recipeRawData.ID).ExecNamed(nil)
		if err != nil {
			return
		}
	} else {
		err = nil
	}
	//创建数据
	for k := 0; k < len(args.RawList); k++ {
		v := args.RawList[k]
		err = weeklyRecipeRawDB.Insert().SetFields([]string{"org_id", "store_id", "weekly_recipe_id", "recipe_type_id", "recipe_type_name", "dining_date", "day_type", "recipe_id", "recipe_name", "recipe_child_id", "material_id", "material_name", "use_count"}).Add(map[string]any{
			"org_id":           recipeData.OrgID,
			"store_id":         recipeData.StoreID,
			"weekly_recipe_id": recipeDayData.WeeklyRecipeID,
			"recipe_type_id":   recipeData.RecipeTypeID,
			"recipe_type_name": recipeData.RecipeTypeName,
			"dining_date":      recipeDayData.DiningDate,
			"day_type":         recipeRawData.DayType,
			"recipe_id":        recipeRawData.RecipeID,
			"recipe_name":      recipeRawData.Name,
			"recipe_child_id":  recipeRawData.ID,
			"material_id":      v.MaterialID,
			"material_name":    RestaurantRawMaterials.GetRawNameByID(v.MaterialID),
			"use_count":        v.UseCount,
		}).ExecAndCheckID()
		if err != nil {
			return
		}
	}
	//反馈
	return
}

// 更新指定周菜谱的原材料参数
type ArgsUpdateWeeklyRecipeRaw struct {
	//周菜品关联行ID
	RecipeChildID int64 `db:"recipe_child_id" json:"recipeChildID" check:"id" index:"true"`
	//原材料ID
	MaterialID int64 `db:"material_id" json:"materialID" check:"id" empty:"true" index:"true"`
	//用量
	UseCount float64 `db:"use_count" json:"useCount" check:"intThan0"`
}

// 更新指定周菜谱的原材料
func UpdateWeeklyRecipeRaw(args *ArgsUpdateWeeklyRecipeRaw) (err error) {
	//更新数据
	err = weeklyRecipeRawDB.Update().SetFields([]string{"material_id", "material_name", "use_count"}).NeedUpdateTime().AddWhereID(args.RecipeChildID).NamedExec(map[string]any{
		"material_id":   args.MaterialID,
		"material_name": RestaurantRawMaterials.GetRawNameByID(args.MaterialID),
		"use_count":     args.UseCount,
	})
	//反馈
	return

}

// DeleteWeeklyRecipeRaw 删除指定周菜谱的原材料
func DeleteWeeklyRecipeRaw(id int64) (err error) {
	err = weeklyRecipeRawDB.Delete().NeedSoft(true).AddWhereID(id).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteWeeklyRecipeRawCache(id)
	return
}

func GetWeeklyRecipeRawData(id int64) (data FieldsWeeklyRecipeRaw) {
	cacheMark := getWeeklyRecipeRawCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := weeklyRecipeRawDB.Get().SetDefaultFields().GetByID(id).Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Hour)
	return
}

// 缓冲
func getWeeklyRecipeRawCacheMark(id int64) string {
	return fmt.Sprint("restaurant:weekly:recipe:raw:id.", id)
}

func deleteWeeklyRecipeRawCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWeeklyRecipeRawCacheMark(id))
}
