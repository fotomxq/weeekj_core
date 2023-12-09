package BlogUserRead

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 删除文章的所有记录
func deleteLogByContentID(contentID int64) {
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "blog_user_read_log", "content_id = :content_id", map[string]interface{}{
		"content_id": contentID,
	})
	deleteLogContentCache(contentID)
}
