package BlogStuRead

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// CheckLogByUserID 检查用户是否在指定时间段学习？
func CheckLogByUserID(orgID int64, userID int64, contentID int64, startAt, endAt time.Time) bool {
	var id int64
	err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM blog_stu_read_log WHERE ($1 < 0 OR org_id = $1) AND user_id = $2 AND ($3 < 1 OR content_id = $3) AND ($4 < to_timestamp(1000000) OR create_at >= $4) AND ($5 < to_timestamp(1000000) OR create_at <= $5)", orgID, userID, contentID, startAt, endAt)
	if err == nil && id > 0 {
		return true
	}
	return false
}

// GetLogByUserTime 获取最近N时间的时间长度
func GetLogByUserTime(userID int64, startAt time.Time) (runTime int64) {
	err := Router2SystemConfig.MainDB.Get(&runTime, "SELECT SUM(run_time) FROM blog_stu_read_log WHERE user_id = $1 AND end_at >= $2", userID, startAt)
	if err != nil {
		return
	}
	return
}
