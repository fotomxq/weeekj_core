package Market2LogMod

import "time"

// CheckHaveLogByFrom 检查被奖励目标奖励是否存在
func CheckHaveLogByFrom(action string, bindID int64, orgID int64, orgBindID int64, userID int64, bindUserID int64) (b bool) {
	logData := GetLastLogByFrom(action, bindID, orgID, orgBindID, userID, bindUserID)
	return logData.ID > 0
}

// CheckLastTimeHaveLogByFrom 检查被奖励目标是否在指定时间段后存在奖励
func CheckLastTimeHaveLogByFrom(action string, bindID int64, orgID int64, orgBindID int64, userID int64, bindUserID int64, afterAt time.Time) (b bool) {
	logData := GetLastLogByFrom(action, bindID, orgID, orgBindID, userID, bindUserID)
	if logData.ID < 1 {
		return
	}
	b = logData.CreateAt.Unix() >= afterAt.Unix()
	return
}
