package BaseSQLTools

import (
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetCacheMarkByInfoID 获取Info缓冲标识码
func (c *Quick) GetCacheMarkByInfoID(id int64) string {
	return fmt.Sprint(c.prefixCacheMark, "id:", id)
}

// GetCacheMarkByListID 获取List缓冲标识码
func (c *Quick) GetCacheMarkByListID(id int64) string {
	return fmt.Sprint(c.prefixCacheMark, "list:", id)
}

// GetCacheInfoByID 获取缓冲Info数据
func (c *Quick) GetCacheInfoByID(id int64, data any) (err error) {
	return Router2SystemConfig.MainCache.GetStruct(c.GetCacheMarkByInfoID(id), data)
}

// GetCacheListByID 获取缓冲List数据
func (c *Quick) GetCacheListByID(id int64, data any) (err error) {
	return Router2SystemConfig.MainCache.GetStruct(c.GetCacheMarkByListID(id), data)
}

// SetCacheInfoByID 设置Info缓冲
func (c *Quick) SetCacheInfoByID(id int64, data any) {
	Router2SystemConfig.MainCache.SetStruct(c.GetCacheMarkByInfoID(id), data, CoreCache.CacheTime1Hour)
}

// SetCacheListByID 设置List缓冲
func (c *Quick) SetCacheListByID(id int64, data any) {
	Router2SystemConfig.MainCache.SetStruct(c.GetCacheMarkByListID(id), data, CoreCache.CacheTime1Hour)
}

// DeleteCacheByID 删除缓冲
func (c *Quick) DeleteCacheByID(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(c.GetCacheMarkByInfoID(id))
	Router2SystemConfig.MainCache.DeleteMark(c.GetCacheMarkByListID(id))
}

// DeleteCachePrefix 删除缓冲
func (c *Quick) DeleteCachePrefix() {
	Router2SystemConfig.MainCache.DeleteSearchMark(c.prefixCacheMark)
}
