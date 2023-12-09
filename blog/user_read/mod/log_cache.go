package BlogUserReadMod

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取缓冲名称
func getLogCacheMark(userID int64, contentID int64) string {
	return fmt.Sprint("blog:user:read:user:", userID, ".", contentID)
}

// 获取缓冲
func getLogCache(userID int64, contentID int64) (data FieldsLog) {
	err := Router2SystemConfig.MainCache.GetStruct(getLogCacheMark(userID, contentID), &data)
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
	Router2SystemConfig.MainCache.SetStruct(getLogCacheMark(data.UserID, data.ContentID), data, 86400)
	Router2SystemConfig.MainCache.SetStruct(getLogCacheByIDMark(data.ID), data, 86400)
}

// 根据ID获取缓冲
func getLogCacheByIDMark(id int64) string {
	return fmt.Sprint("blog:user:read:id:", id)
}

func getLogCacheByID(id int64) (data FieldsLog) {
	err := Router2SystemConfig.MainCache.GetStruct(getLogCacheByIDMark(id), &data)
	if err != nil {
		return
	}
	if data.ID < 1 {
		return
	}
	return
}

// 删除缓冲
func deleteLogCache(userID int64, contentID int64) {
	data := getLogCache(userID, contentID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheByIDMark(data.ID))
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(data.UserID, data.ContentID))
}
