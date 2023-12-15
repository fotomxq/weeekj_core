package BaseToken2

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func getTokenCacheMark(id int64) string {
	return fmt.Sprint("base:token2:id:", id)
}

func getTokenFromCacheMark(userID, deviceID int64, loginFrom string) string {
	return fmt.Sprint("base:token2:from:", userID, ".", deviceID, ".", loginFrom)
}

func deleteTokenCache(id int64) {
	data := getByID(id)
	Router2SystemConfig.MainCache.DeleteMark(getTokenCacheMark(data.ID))
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.DeleteMark(getTokenFromCacheMark(data.UserID, data.DeviceID, data.LoginFrom))
}
