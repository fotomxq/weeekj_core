package BlogStuRead

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 获取文章缓冲标识码
func getContentCacheMark(id int64) string {
	return fmt.Sprint("blog:stu:read:content:id:", id)
}

// url特殊的缓冲
func getContentURLListCacheMark(id int64) string {
	return fmt.Sprint("blog:stu:read:url:list:", id)
}

func getContentURLDataCacheMark(id int64) string {
	return fmt.Sprint("blog:stu:read:url:data:", id)
}

// 删除缓冲
func deleteContentCacheByID(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getContentCacheMark(id))
	Router2SystemConfig.MainCache.DeleteMark(getContentURLListCacheMark(id))
	Router2SystemConfig.MainCache.DeleteMark(getContentURLDataCacheMark(id))
}
