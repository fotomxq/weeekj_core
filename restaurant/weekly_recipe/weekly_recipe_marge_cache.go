package RestaurantWeeklyRecipeMarge

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 缓冲
func getWeeklyRecipeListCacheMark(id int64) string {
	return fmt.Sprint("restaurant:weekly_recipe:marge:list.id.", id)
}

func getWeeklyRecipeCacheMark(id int64) string {
	return fmt.Sprint("restaurant:weekly_recipe:marge:id.", id)
}

func deleteWeeklyRecipeCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getWeeklyRecipeListCacheMark(id))
	Router2SystemConfig.MainCache.DeleteMark(getWeeklyRecipeCacheMark(id))
}
