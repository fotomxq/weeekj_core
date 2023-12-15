package ServiceUserInfo

import (
	"fmt"
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
	"time"
)

// 更新组织的信息档案统计
func updateAnalysisOrg(orgID int64) {
	//统计机构下总人数
	var allCount int64
	_ = Router2SystemConfig.MainDB.Get(&allCount, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", orgID)
	AnalysisAny2.AppendData("re", "service_user_info_all_count", time.Time{}, orgID, 0, 0, 0, 0, allCount)
	var count int64
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND die_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", orgID)
	AnalysisAny2.AppendData("re", "service_user_info_count", time.Time{}, orgID, 0, 0, 0, 0, count)
	//男性人数
	var countMan int64
	_ = Router2SystemConfig.MainDB.Get(&countMan, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND gender = 0 AND die_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", orgID)
	AnalysisAny2.AppendData("re", "service_user_info_gender_count", time.Time{}, orgID, 0, 0, 0, 0, countMan)
	//女性人数
	var countWoman int64
	_ = Router2SystemConfig.MainDB.Get(&countWoman, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND gender = 1 AND die_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", orgID)
	AnalysisAny2.AppendData("re", "service_user_info_gender_count", time.Time{}, orgID, 0, 1, 0, 0, countWoman)
	//75岁以上
	var count75More int64
	_ = Router2SystemConfig.MainDB.Get(&count75More, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND date_of_birth < $2 AND die_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", orgID, CoreFilter.GetNowTimeCarbon().SubYears(75).Time)
	AnalysisAny2.AppendData("re", "service_user_info_old_75_count", time.Time{}, orgID, 0, 0, 0, 0, count75More)
	//75岁以下
	var count75Less int64
	_ = Router2SystemConfig.MainDB.Get(&count75Less, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND date_of_birth >= $2 AND die_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", orgID, CoreFilter.GetNowTimeCarbon().SubYears(75).Time)
	AnalysisAny2.AppendData("re", "service_user_info_old_75_count", time.Time{}, orgID, 0, 1, 0, 0, count75Less)
	//看护级别
	updateAnalysisInfoLevel(orgID, 0)
	updateAnalysisInfoLevel(orgID, 1)
	updateAnalysisInfoLevel(orgID, 2)
	updateAnalysisInfoLevel(orgID, 3)
	//统计去世人员
	var countDie int64
	_ = Router2SystemConfig.MainDB.Get(&countDie, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND die_at > to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", orgID)
	AnalysisAny2.AppendData("re", "service_user_info_die_count", time.Time{}, orgID, 0, 0, 0, 0, countDie)
	//统计离开人员
	var countOut int64
	_ = Router2SystemConfig.MainDB.Get(&countOut, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND out_at > to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", orgID)
	AnalysisAny2.AppendData("re", "service_user_info_out_count", time.Time{}, orgID, 0, 0, 0, 0, countOut)
	//按照月份统计的数据
	updateAnalysisTime(orgID, CoreFilter.GetNowTimeCarbon())
	//TODO：补救性修正
	updateAnalysisTime(orgID, CoreFilter.GetNowTimeCarbon().SubMonth())
	updateAnalysisTime(orgID, CoreFilter.GetNowTimeCarbon().SubMonths(2))
	updateAnalysisTime(orgID, CoreFilter.GetNowTimeCarbon().SubMonths(3))

}

// 带有时间范围的统计
func updateAnalysisTime(orgID int64, createAt carbon.Carbon) {
	startAt := createAt.StartOfMonth()
	endAt := createAt.EndOfMonth()
	//总人数变动，当月新入住人数
	var countMonth int64
	_ = Router2SystemConfig.MainDB.Get(&countMonth, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND die_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000) AND create_at >= $2 AND create_at <= $3", orgID, startAt.Time, endAt.Time)
	AnalysisAny2.AppendData("re", "service_user_info_month_count", createAt.Time, orgID, 0, 0, 0, 0, countMonth)
	//统计当月死亡
	var countMonthDie int64
	_ = Router2SystemConfig.MainDB.Get(&countMonthDie, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND die_at > to_timestamp(1000000) AND delete_at < to_timestamp(1000000) AND die_at >= $2 AND die_at <= $3", orgID, startAt.Time, endAt.Time)
	AnalysisAny2.AppendData("re", "service_user_info_die_month_count", createAt.Time, orgID, 0, 0, 0, 0, countMonthDie)
	//统计当月离开
	var countMonthOut int64
	_ = Router2SystemConfig.MainDB.Get(&countMonthOut, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND out_at > to_timestamp(1000000) AND delete_at < to_timestamp(1000000) AND out_at >= $2 AND out_at <= $3", orgID, startAt.Time, endAt.Time)
	AnalysisAny2.AppendData("re", "service_user_info_out_month_count", createAt.Time, orgID, 0, 0, 0, 0, countMonthOut)
}

// 更新老人统计
func updateAnalysisInfoLevel(orgID int64, level int) {
	var countLevel int64
	_ = Router2SystemConfig.MainDB.Get(&countLevel, "SELECT COUNT(id) FROM service_user_info WHERE org_id = $1 AND level = $2 AND die_at < to_timestamp(1000000) AND out_at < to_timestamp(1000000) AND delete_at < to_timestamp(1000000)", orgID, level)
	AnalysisAny2.AppendData("re", fmt.Sprint("service_user_info_level_count"), time.Time{}, orgID, 0, int64(level), 0, 0, countLevel)
}
