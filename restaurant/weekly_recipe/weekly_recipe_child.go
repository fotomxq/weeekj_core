package RestaurantWeeklyRecipeMarge

import (
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	RestaurantRecipe "github.com/fotomxq/weeekj_core/v5/restaurant/recipe"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetWeeklyRecipeChild 获取日明细数据参数
type ArgsGetWeeklyRecipeChild struct {
	//每周菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id" index:"true"`
	//每日菜谱ID
	WeeklyRecipeDayID int64 `db:"weekly_recipe_day_id" json:"weeklyRecipeDayID" check:"id" index:"true"`
	//每日类型
	// 1:早餐; 2:中餐; 3:晚餐
	DayType int `db:"day_type" json:"dayType" check:"intThan0" empty:"true" index:"true"`
}

// GetWeeklyRecipeChild 获取日明细数据
func GetWeeklyRecipeChild(args *ArgsGetWeeklyRecipeChild) (dataList []FieldsWeeklyRecipeChild, err error) {
	cacheMark := getWeeklyRecipeChildCacheMark(args.WeeklyRecipeID, args.WeeklyRecipeDayID, args.DayType)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	err = weeklyRecipeChildDB.Select().SetFieldsSortDefault().SetFieldsAll().SetIDQuery("weekly_recipe_id", args.WeeklyRecipeID).SetIDQuery("weekly_recipe_day_id", args.WeeklyRecipeDayID).SetIntQuery("day_type", args.DayType).SetDeleteQuery("delete_at", false).SelectList("").Result(&dataList)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, CoreCache.CacheTime1Hour)
	return
}

func getWeeklyRecipeChildByID(id int64) (data FieldsWeeklyRecipeChild) {
	err := weeklyRecipeChildDB.Get().SetDefaultFields().GetByID(id).Result(&data)
	if err != nil {
		return
	}
	return
}

// SetWeeklyRecipeChild 修改日明细数据
func SetWeeklyRecipeChild(weeklyRecipeID int64, weeklyRecipeDayID int64, dayType int, newData []DataGetWeeklyRecipeMargeDayItem) (dataList []DataGetWeeklyRecipeMargeDayItem, err error) {
	//检查是否存在数据
	var rawList []FieldsWeeklyRecipeChild
	rawList, err = GetWeeklyRecipeChild(&ArgsGetWeeklyRecipeChild{
		WeeklyRecipeID:    weeklyRecipeID,
		WeeklyRecipeDayID: weeklyRecipeDayID,
		DayType:           dayType,
	})
	if err == nil && len(rawList) > 0 {
		//清理并重建数据
		err = weeklyRecipeChildDB.Delete().NeedSoft(true).SetWhereAnd("weekly_recipe_id", weeklyRecipeID).SetWhereAnd("weekly_recipe_day_id", weeklyRecipeDayID).SetWhereAnd("day_type", dayType).ExecNamed(nil)
		if err != nil {
			return
		}
	} else {
		err = nil
	}
	//创建数据
	for k := 0; k < len(newData); k++ {
		v := newData[k]
		if v.Name == "" {
			v.Name = RestaurantRecipe.GetRecipeNameByID(v.RecipeID)
		}
		err = weeklyRecipeChildDB.Insert().SetFields([]string{"weekly_recipe_id", "weekly_recipe_day_id", "day_type", "recipe_id", "name", "price", "recipe_count", "unit", "unit_id"}).Add(map[string]any{
			"weekly_recipe_id":     weeklyRecipeID,
			"weekly_recipe_day_id": weeklyRecipeDayID,
			"day_type":             dayType,
			"recipe_id":            v.RecipeID,
			"name":                 v.Name,
			"price":                v.Price,
			"recipe_count":         v.RecipeCount,
			"unit":                 v.Unit,
			"unit_id":              v.UnitID,
		}).ExecAndCheckID()
		if err != nil {
			return
		}
		dataList = append(dataList, DataGetWeeklyRecipeMargeDayItem{
			ID:          v.ID,
			RecipeID:    v.RecipeID,
			Name:        v.Name,
			Price:       v.Price,
			RecipeCount: v.RecipeCount,
			Unit:        v.Unit,
			UnitID:      v.UnitID,
			IsRepeat:    false,
			IsRepeatAll: false,
		})
	}
	//删除缓冲
	deleteWeeklyRecipeChildCache(weeklyRecipeID, weeklyRecipeDayID, dayType)
	//反馈
	return
}

// 缓冲
func getWeeklyRecipeChildCacheMark(weeklyRecipeID, weeklyRecipeDayID int64, dayType int) string {
	return fmt.Sprint("restaurant:weekly.recipe:day:id.", weeklyRecipeID, ".", weeklyRecipeDayID, ".", dayType)
}

func deleteWeeklyRecipeChildCache(weeklyRecipeID, weeklyRecipeDayID int64, dayType int) {
	Router2SystemConfig.MainCache.DeleteMark(getWeeklyRecipeChildCacheMark(weeklyRecipeID, weeklyRecipeDayID, dayType))
}
