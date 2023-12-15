package BlogExam

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// GetTopicGroupLogCount 获取用户学习天数
func GetTopicGroupLogCount(userID int64) (count int64) {
	cacheMark := getTopicGroupLogCountCacheMark(userID)
	var err error
	count, err = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err == nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM blog_exam_topic_group_log WHERE user_id = $1", userID)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetInt64(cacheMark, count, 1800)
	return
}

// GetTopicGroupLogSign 获取用户最近签到情况
func GetTopicGroupLogSign(userID int64, limit int64) (dataList []FieldsTopicGroupLog) {
	cacheMark := getTopicGroupLogSignCacheMark(userID, limit)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil {
		return
	}
	err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_day, run_time, org_id, user_id FROM blog_exam_topic_group_log WHERE user_id = $1 ORDER BY id DESC LIMIT $2", userID, limit)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, 1800)
	return
}

// GetTopicGroupLogRunTime 获取累计学习时间长度
func GetTopicGroupLogRunTime(userID int64, afterAt time.Time) (runTime int) {
	cacheMark := getTopicGroupLogRunTimeCacheMark(userID, afterAt.Unix())
	var err error
	runTime, err = Router2SystemConfig.MainCache.GetInt(cacheMark)
	if err == nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&runTime, "SELECT SUM(run_time) FROM blog_exam_topic_group_log WHERE user_id = $1 AND create_day >= $2", userID, afterAt)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetInt(cacheMark, runTime, 1800)
	return
}

// CreateTopicGroupLog 添加学习记录
func CreateTopicGroupLog(orgID, userID int64, runTime int) {
	//检查是否存在记录
	var id int64
	err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM blog_exam_topic_group_log WHERE user_id = $1 AND create_day >= $2", userID, CoreFilter.GetNowTimeCarbon().StartOfDay().Time)
	if err == nil && id > 0 {
		//更新记录
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE blog_exam_topic_group_log SET run_time = run_time + :run_time WHERE id = :id", map[string]interface{}{
			"id":       id,
			"run_time": runTime,
		})
		if err != nil {
			return
		}
	} else {
		//添加记录
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO blog_exam_topic_group_log (run_time, org_id, user_id) VALUES (:run_time,:org_id,:user_id)", map[string]interface{}{
			"run_time": runTime,
			"org_id":   orgID,
			"user_id":  userID,
		})
		if err != nil {
			return
		}
	}
	//清理缓冲
	deleteTopicGroupLogCache(userID)
	//反馈
	return
}

// 缓冲
func getTopicGroupLogCountCacheMark(userID int64) string {
	return fmt.Sprint("blog:exam:topic:group:log:count:user:", userID)
}
func getTopicGroupLogSignCacheMark(userID int64, limit int64) string {
	if limit < 1 {
		return fmt.Sprint("blog:exam:topic:group:log:sign:", userID)
	} else {
		return fmt.Sprint("blog:exam:topic:group:log:sign:", userID, ".", limit)
	}
}
func getTopicGroupLogRunTimeCacheMark(userID int64, afterAtUnix int64) string {
	if afterAtUnix < 1 {
		return fmt.Sprint("blog:exam:topic:group:log:runtime:user:", userID)
	} else {
		return fmt.Sprint("blog:exam:topic:group:log:runtime:user:", userID, ".", afterAtUnix)
	}
}

// 清理缓冲
func deleteTopicGroupLogCache(userID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getTopicGroupLogCountCacheMark(userID))
	Router2SystemConfig.MainCache.DeleteSearchMark(getTopicGroupLogSignCacheMark(userID, -1))
	Router2SystemConfig.MainCache.DeleteSearchMark(getTopicGroupLogRunTimeCacheMark(userID, -1))
}
