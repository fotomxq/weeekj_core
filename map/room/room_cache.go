package MapRoom

import (
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 获取模块前缀
func getRoomCacheMark(id int64) string {
	return fmt.Sprint("map:room:id:", id)
}

func getRoomInfoCacheMark(infoID int64) string {
	return fmt.Sprint("map:room:info:", infoID)
}

// 清理缓冲
func deleteRoomCache(id int64) {
	roomData := getRoomID(id)
	Router2SystemConfig.MainCache.DeleteMark(getRoomCacheMark(id))
	if roomData.ID > 0 {
		for _, v := range roomData.Infos {
			Router2SystemConfig.MainCache.DeleteMark(getRoomInfoCacheMark(v))
		}
	}
}
