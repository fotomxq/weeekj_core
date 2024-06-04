package RestaurantWeeklyRecipeMarge

import (
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetWeeklyRecipeDay 获取日数据参数
type ArgsGetWeeklyRecipeDay struct {
	//每周菜谱ID
	WeeklyRecipeID int64 `db:"weekly_recipe_id" json:"weeklyRecipeID" check:"id" index:"true"`
}

// GetWeeklyRecipeDay 获取日数据
func GetWeeklyRecipeDay(args *ArgsGetWeeklyRecipeDay) (dataList []FieldsWeeklyRecipeDay, err error) {
	cacheMark := getWeeklyRecipeDayCacheMark(args.WeeklyRecipeID)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	err = weeklyRecipeDayDB.Select().SetFieldsList([]string{"id", "create_at", "update_at", "delete_at", "weekly_recipe_id", "dining_date"}).SetIDQuery("weekly_recipe_id", args.WeeklyRecipeID).SetDeleteQuery("delete_at", false).SelectList("").Result(&dataList)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, CoreCache.CacheTime1Hour)
	return
}

// SetWeeklyRecipeDay 创建日数据
func SetWeeklyRecipeDay(weeklyRecipeID int64, newData []DataGetWeeklyRecipeMargeDay) (dataList []DataGetWeeklyRecipeMargeDay, err error) {
	//检查是否存在数据
	var rawList []FieldsWeeklyRecipeDay
	rawList, err = GetWeeklyRecipeDay(&ArgsGetWeeklyRecipeDay{
		WeeklyRecipeID: weeklyRecipeID,
	})
	if err == nil && len(rawList) > 0 {
		//清理并重建数据
		err = weeklyRecipeDayDB.Delete().NeedSoft(true).SetWhereOrThan("weekly_recipe_id", weeklyRecipeID).ExecNamed(nil)
		if err != nil {
			return
		}
	} else {
		err = nil
	}
	//创建数据
	for k := 0; k < len(newData); k++ {
		v := newData[k]
		var newID int64
		newID, err = weeklyRecipeDayDB.Insert().SetFields([]string{"weekly_recipe_id", "recipe_type_id", "recipe_type_name", "recipe_type_name", "dining_date"}).Add(map[string]any{
			"weekly_recipe_id": weeklyRecipeID,
			"recipe_type_id":   v.RecipeTypeID,
			"recipe_type_name": RecipeType.GetNameNoErr(v.RecipeTypeID),
			"dining_date":      v.DiningDate,
		}).ExecAndResultID()
		if err != nil {
			return
		}
		var breakfast, lunch, dinner []DataGetWeeklyRecipeMargeDayItem
		breakfast, err = SetWeeklyRecipeChild(weeklyRecipeID, newID, 1, v.Breakfast)
		if err != nil {
			return
		}
		lunch, err = SetWeeklyRecipeChild(weeklyRecipeID, newID, 2, v.Lunch)
		if err != nil {
			return
		}
		dinner, err = SetWeeklyRecipeChild(weeklyRecipeID, newID, 3, v.Dinner)
		if err != nil {
			return
		}
		dataList = append(dataList, DataGetWeeklyRecipeMargeDay{
			DiningDate:     v.DiningDate,
			RecipeTypeID:   v.RecipeTypeID,
			RecipeTypeName: RecipeType.GetNameNoErr(v.RecipeTypeID),
			Breakfast:      breakfast,
			Lunch:          lunch,
			Dinner:         dinner,
		})
	}
	//删除缓冲
	deleteWeeklyRecipeDayCache(weeklyRecipeID)
	//反馈
	return
}

// 缓冲
func getWeeklyRecipeDayCacheMark(weeklyRecipeID int64) string {
	return fmt.Sprint("restaurant:weekly.recipe:day:id.", weeklyRecipeID)
}

func deleteWeeklyRecipeDayCache(weeklyRecipeID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWeeklyRecipeDayCacheMark(weeklyRecipeID))
}
