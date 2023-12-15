package BlogUserRead

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 获取缓冲名称
func getLogCacheMark(contentID int64, userID int64) string {
	return fmt.Sprint("blog:user:read:user:", contentID, ".", userID)
}
func getLogContentCacheMark(contentID int64) string {
	return fmt.Sprint("blog:user:read:user:", contentID)
}

// 获取缓冲
func getLogCache(contentID int64, userID int64) (data FieldsLog) {
	err := Router2SystemConfig.MainCache.GetStruct(getLogCacheMark(contentID, userID), &data)
	if err != nil {
		return
	}
	if data.ID < 1 {
		return
	}
	return
}

// 写入缓冲
func setLogCache(data FieldsLog) {
	Router2SystemConfig.MainCache.SetStruct(getLogCacheMark(data.ContentID, data.UserID), data, 1800)
}

// 删除缓冲
func deleteLogTargetCache(contentID int64, userID int64) {
	data := getLogCache(contentID, userID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(data.ContentID, data.UserID))
}

func deleteLogContentCache(contentID int64) {
	Router2SystemConfig.MainCache.DeleteSearchMark(getLogContentCacheMark(contentID))
}
