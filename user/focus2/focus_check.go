package UserFocus2

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// CheckFocus 检查是否关注了数据
func CheckFocus(userID int64, mark string, system string, bindID int64) bool {
	if err := checkMark(mark); err != nil {
		return false
	}
	var data FieldsFocus
	cacheMark := getFocusUserCacheMark(userID, mark, system, bindID)
	if id, err := Router2SystemConfig.MainCache.GetInt64(cacheMark); err == nil && id > 0 {
		return true
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_focus2 WHERE user_id = $1 AND mark = $2 AND system = $3 AND bind_id = $4", userID, mark, system, bindID)
	if data.ID < 1 {
		return false
	}
	Router2SystemConfig.MainCache.SetInt64(cacheMark, data.ID, 1800)
	return true
}
