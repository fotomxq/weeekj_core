package BlogUserReadMod

import Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

// CheckUserLogExist 检查用户是否存在访问记录
func CheckUserLogExist(userID int64, contentID int64) bool {
	if userID < 1 {
		return false
	}
	data := getLogCache(userID, contentID)
	if data.ID > 0 {
		return true
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, child_org_id, user_id, from_mark, from_name, name, ip, sort_id, content_id, leave_at, read_time FROM blog_user_read_log WHERE content_id = $1 AND user_id = $2", contentID, userID)
	if err == nil && data.ID > 0 {
		setLogCache(data)
		return true
	} else {
		return false
	}
}
