package OrgTime

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	OrgCoreCoreMod "github.com/fotomxq/weeekj_core/v5/org/core/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ToolsHolidaySeason "github.com/fotomxq/weeekj_core/v5/tools/holiday_season"
	"github.com/golang-module/carbon"
)

// ArgsCheckIsWorkByID 检查某个ID，是否正在上班？参数
type ArgsCheckIsWorkByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	// 用于验证
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
}

// checkIsWorkByID 检查某个ID，是否正在上班？
func checkIsWorkByID(args *ArgsCheckIsWorkByID) bool {
	data := getConfigByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		return false
	}
	if !CoreSQL.CheckTimeThanNow(data.ExpireAt) {
		return false
	}
	return data.IsWork
}

// CheckIsWorkByOrgBindID 检查组织成语ID是否上班？
func CheckIsWorkByOrgBindID(orgBindID int64) bool {
	//获取组织成员信息
	bindData := OrgCoreCoreMod.GetBind(orgBindID, -1, -1)
	if bindData.ID < 1 {
		return false
	}
	//遍历检查是否上班？
	var data FieldsWorkTime
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT is_work FROM org_work_time WHERE (expire_at > NOW() OR expire_at < to_timestamp(1000000)) AND is_work = true AND ($1 && groups OR $2 = ANY(binds)) LIMIT 1", bindData.GroupIDs, bindData.ID); err != nil {
		return false
	}
	//检查是否请假
	isLeave := CheckLeaveByBindID(orgBindID)
	if isLeave {
		return false
	}
	//反馈
	return data.IsWork
}

// checkIsWorkByData 检查某个数据，是否正在上班？
func checkIsWorkByData(data *FieldsWorkTime) bool {
	//检查日周期
	b := checkIsWorkNowDayByData(data)
	if !b {
		return false
	}
	//检查24小时
	nowTime := CoreFilter.GetNowTimeCarbon()
	for _, v := range data.Configs.WorkTime {
		//检查是否为节假日？
		if !data.Configs.AllowHoliday && ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
			DateAt: nowTime.Time,
		}) {
			continue
		}
		//获取考勤的开始和结束时间
		var startAt, endAt carbon.Carbon
		//计算当前时间是否在该时间段内
		if checkTimeThan(v, nowTime) {
			startAt = nowTime.SetHour(v.StartHour).SetMinute(v.StartMinute).StartOfMinute()
			endAt = nowTime.SetHour(v.EndHour).SetMinute(v.EndMinute).StartOfMinute()
			//fmt.Println("startAt: ", startAt, ", endAt: ", endAt, ", nowTime: ", nowTime)
			if nowTime.Time.Unix() >= startAt.Time.Unix() && nowTime.Time.Unix() <= endAt.Time.Unix() {
				return true
			}
		} else {
			startAt = nowTime.SubDay().SetHour(v.StartHour).SetMinute(v.StartMinute).StartOfMinute()
			endAt = nowTime.SetHour(v.EndHour).SetMinute(v.EndMinute).StartOfMinute()
			if nowTime.Time.Unix() >= startAt.Time.Unix() && nowTime.Time.Unix() <= endAt.Time.Unix() {
				return true
			}
		}
	}
	if len(data.RotConfig.WorkTime) > 0 {
		if data.RotConfig.NowKey < 1 {
			data.RotConfig.NowKey = 0
		}
		if data.RotConfig.NowKey >= len(data.RotConfig.WorkTime) {
			data.RotConfig.NowKey = len(data.RotConfig.WorkTime) - 1
		}
		rotConfig := data.RotConfig.WorkTime[data.RotConfig.NowKey]
		//获取考勤的开始和结束时间
		var startAt, endAt carbon.Carbon
		//是否存在跨天，计算规则方向相反
		if checkTimeThan(rotConfig, nowTime) {
			startAt = nowTime.SetHour(rotConfig.StartHour).SetMinute(rotConfig.StartMinute).StartOfMinute()
			endAt = nowTime.SetHour(rotConfig.EndHour).SetMinute(rotConfig.EndMinute).StartOfMinute()
			//计算当前时间是否在该时间段内
			if nowTime.Time.Unix() >= startAt.Time.Unix() && nowTime.Time.Unix() <= endAt.Time.Unix() {
				return true
			}
		} else {
			//检查考勤是否满足？
			startAt = nowTime.SubDay().SetHour(rotConfig.StartHour).SetMinute(rotConfig.StartMinute).StartOfMinute()
			endAt = nowTime.SetHour(rotConfig.EndHour).SetMinute(rotConfig.EndMinute).StartOfMinute()
			//fmt.Println("startAt: ", startAt, ", endAt: ", endAt, ", nowTime: ", nowTime)
			//计算当前时间是否在该时间段内
			if nowTime.Time.Unix() >= startAt.Time.Unix() && nowTime.Time.Unix() <= endAt.Time.Unix() {
				return true
			}
		}
	}
	return false
}

// checkIsWorkNowDayByData 检查今天是否需要上班？
// 条件全部满足后继续
func checkIsWorkNowDayByData(data *FieldsWorkTime) bool {
	nowTime := CoreFilter.GetNowTimeCarbon()
	isFind := false
	//检查这个月是否存在满足条件
	if len(data.Configs.Month) > 0 {
		nowMonth := nowTime.Month()
		for _, v := range data.Configs.Month {
			if v == nowMonth {
				isFind = true
				break
			}
		}
		if !isFind {
			return false
		}
	}
	//检查本天是否需要上班
	if len(data.Configs.MonthDay) > 0 {
		isFind = false
		nowDay := nowTime.Day()
		for _, v := range data.Configs.MonthDay {
			if v == nowDay {
				isFind = true
				break
			}
		}
		if !isFind {
			return false
		}
	}
	//检查是每月第几周
	if len(data.Configs.MonthWeek) > 0 {
		isFind = false
		nowMonthWeek := nowTime.WeekOfMonth()
		for _, v := range data.Configs.MonthWeek {
			if v == nowMonthWeek {
				isFind = true
				break
			}
		}
		if !isFind {
			return false
		}
	}
	//检查本周是否需要上班
	if len(data.Configs.Week) > 0 {
		isFind = false
		nowWeekDay := nowTime.DayOfWeek()
		for _, v := range data.Configs.Week {
			if v == nowWeekDay {
				isFind = true
				break
			}
		}
		if !isFind {
			return false
		}
	}
	return true
}

// getWorkDayInMonth 检查本月需要上几天班？
func getWorkDayInMonth(data *FieldsWorkTime) int {
	if len(data.Configs.Month) > 0 {
		return len(data.Configs.Month)
	}
	return len(data.Configs.Week) * 4
}

// 检查第一个时间是否大于第二个时间
// 用于判断是否存在跨天考勤情况
func checkTimeThan(data FieldsWorkTimeTime, nowTime carbon.Carbon) bool {
	if nowTime.SetHour(data.StartHour).SetMinute(data.EndMinute).Time.Unix() <= nowTime.SetHour(data.EndHour).SetMinute(data.EndMinute).Time.Unix() {
		return true
	}
	return false
}
