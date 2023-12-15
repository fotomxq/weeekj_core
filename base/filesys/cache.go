package BaseFileSys

import (
	"errors"
	"fmt"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 获取指定ID缓冲
func getClaimCacheData(claimID int64, orgID int64, userID int64) (data FieldsFileClaimType, err error) {
	cacheMark := fmt.Sprint(getClaimCacheMark(claimID))
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		if orgID > 0 && data.OrgID != orgID {
			err = errors.New("no data")
			return
		}
		if userID > 0 && data.UserID != userID {
			err = errors.New("no data")
			return
		}
		return
	}
	err = errors.New("no data")
	return
}

func getFileCacheData(fileID int64, fromInfo CoreSQLFrom.FieldsFrom) (data FieldsFileType, err error) {
	cacheMark := fmt.Sprint(getClaimCacheMark(fileID))
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		if data.FromInfo.System != "" && !fromInfo.CheckEg(data.FromInfo) {
			err = errors.New("no data")
			return
		}
		return
	}
	err = errors.New("no data")
	return
}

func getFileCacheDataByCreate(fileID int64, createInfo CoreSQLFrom.FieldsFrom) (data FieldsFileType, err error) {
	cacheMark := fmt.Sprint(getClaimCacheMark(fileID))
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		if data.CreateInfo.System == "" {
			return
		}
		if createInfo.CheckEg(data.CreateInfo) {
			return
		}
		err = errors.New("no data")
		return
	}
	err = errors.New("no data")
	return
}

// 写入缓冲
func setClaimCacheData(data FieldsFileClaimType) {
	cacheMark := fmt.Sprint(getClaimCacheMark(data.ID))
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTime)
}

func setFileCacheData(data FieldsFileType) {
	cacheMark := fmt.Sprint(getFileCacheMark(data.ID))
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheTime)
}

// 获取缓冲前缀
func getFileCacheMark(fileID int64) string {
	if fileID > 0 {
		return fmt.Sprint("base:filesys:file:", fileID)
	}
	return fmt.Sprint("base:filesys:file:")
}

func getClaimCacheMark(claimID int64) string {
	if claimID > 0 {
		return fmt.Sprint("base:filesys:claim:", claimID)
	}
	return fmt.Sprint("base:filesys:claim:")
}

// 清理缓冲
func clearFileCache(fileID int64) {
	if fileID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getFileCacheMark(fileID))
		Router2SystemConfig.MainCache.DeleteSearchMark(getFileCacheMark(fileID))
	} else {
		Router2SystemConfig.MainCache.DeleteSearchMark(getFileCacheMark(0))
	}
}

func clearClaimCache(claimID int64) {
	if claimID > 0 {
		Router2SystemConfig.MainCache.DeleteMark(getClaimCacheMark(claimID))
		Router2SystemConfig.MainCache.DeleteSearchMark(getClaimCacheMark(claimID))
	} else {
		Router2SystemConfig.MainCache.DeleteSearchMark(getClaimCacheMark(0))
	}
}
