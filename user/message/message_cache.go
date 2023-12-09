package UserMessage

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取缓冲名称
func getMessageCacheMark(id int64) string {
	return fmt.Sprint("user:message:id:", id)
}

func getReceiveMessageCountCacheMark(userID int64) string {
	return fmt.Sprint("user:message:count:receive:", userID)
}

// 删除缓冲
func deleteMessageCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getMessageCacheMark(id))
}

func deleteMessageReceiveCountCache(userID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getReceiveMessageCountCacheMark(userID))
}
