package CoreNextTime

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"github.com/golang-module/carbon"
)

// CheckTimeByTimeN2 规范的时间周期长度检查
// 检查当前时间是否符合条件
// prevTime 上次检查时间，一般存在配置内，方便计算下一次执行时间差异
// addTime 要计算的时间，默认为当前时间
// timeType 时间类型
// 0 每天重复 day / 1 每周重复 week / 2 每月重复 month /
// 3 临时1次 once /
// 4 每隔N天重复 day_n / 5 每隔N周重复 week_n / 6 每隔N月重复 month_n /
// 7 每个星期N重复 week_n / 8 每隔N小时重复 hour_n
// timeN 扩展N
//
//	重复时间内，数组的第一个值作为相隔N；
//	重复周内，数组代表指定的星期1-7
func CheckTimeByTimeN2(prevAt carbon.Carbon, checkAt carbon.Carbon, timeType int, timeN []int64) bool {
	//根据时间比对该设定是否满足？
	switch timeType {
	case 0:
		//每天
		if prevAt.AddDays(1).Time.Unix() > checkAt.Time.Unix() {
			return true
		}
		return false
	case 1:
		//每周
		if prevAt.AddWeek().Time.Unix() > checkAt.Time.Unix() {
			return true
		}
		return false
	case 2:
		//每月
		if prevAt.AddMonth().Time.Unix() > checkAt.Time.Unix() {
			return true
		}
		return false
	case 3:
		//临时一次，如果存在上一次，则反馈失败
		if prevAt.Time.Unix() > 100000 {
			return false
		}
		return true
	case 4:
		//每隔N天
		if len(timeN) < 1 {
			return false
		}
		if prevAt.AddDays(int(timeN[0])).Time.Unix() > checkAt.Time.Unix() {
			return true
		}
		return false
	case 5:
		//每隔N周
		if len(timeN) < 1 {
			return false
		}
		if prevAt.AddWeeks(int(timeN[0])).Time.Unix() > checkAt.Time.Unix() {
			return true
		}
		return false
	case 6:
		//每月的N日
		// 检查检查时间，是否在该月份的时间段内
		if len(timeN) < 1 {
			return false
		}
		for _, v := range timeN {
			if int(v) == checkAt.DayOfMonth() {
				return true
			}
		}
		return false
	case 7:
		//每周的星期几
		// 检查检查时间，是否在该周的时间内
		if len(timeN) < 1 {
			return false
		}
		for _, v := range timeN {
			if int(v) == checkAt.DayOfWeek() {
				return true
			}
		}
		return false
	case 8:
		//每隔N小时
		if len(timeN) < 1 {
			return false
		}
		if prevAt.AddHours(int(timeN[0])).Time.Unix() > checkAt.Time.Unix() {
			return true
		}
		return false
	}
	return false
}

// GetTimeByTimeN2 规范的时间周期长度
// 制定高级别的时间周期设计，计算出下一次检查创建的时间
func GetTimeByTimeN2(addTime carbon.Carbon, timeType int, timeN []int64) (carbon.Carbon, bool) {
	switch timeType {
	case 0:
		//每天
		return addTime.AddDay(), true
	case 1:
		//每周
		return addTime.AddWeek(), true
	case 2:
		//每月
		return addTime.AddMonth(), true
	case 3:
		//临时一次，如果存在上一次，则反馈失败
		return addTime, true
	case 4:
		//每相N天
		if len(timeN) < 1 {
			return carbon.Carbon{}, false
		}
		return addTime.AddDays(int(timeN[0])), true
	case 5:
		//每N个星期
		// 检查检查时间，是否在该周的时间内
		if len(timeN) < 1 {
			return carbon.Carbon{}, false
		}
		return addTime.AddWeeks(int(timeN[0])), true
	case 6:
		//每月的N日
		// 检查检查时间，是否在该月份的时间段内
		if len(timeN) < 1 {
			return carbon.Carbon{}, false
		}
		step := 0
		for {
			addTime = addTime.AddDay()
			for _, v := range timeN {
				if addTime.DayOfMonth() == int(v) {
					return addTime, true
				}
			}
			step += 1
			if step > 90 {
				break
			}
		}
		return carbon.Carbon{}, false
	case 7:
		//每星期N重复
		if len(timeN) < 1 {
			return carbon.Carbon{}, false
		}
		step := 0
		for {
			addTime = addTime.AddDay()
			for _, v := range timeN {
				if addTime.DayOfWeek() == int(v) {
					return addTime, true
				}
			}
			step += 1
			if step > 30 {
				break
			}
		}
		return carbon.Carbon{}, false
	case 8:
		//每N小时重复
		if len(timeN) < 1 {
			return carbon.Carbon{}, false
		}
		addTime = addTime.AddHours(int(timeN[0]))
		return addTime, true
	}
	return carbon.Carbon{}, false
}

// GetTimeByTimeN2Now 递归处理直到时间抵达今日
func GetTimeByTimeN2Now(addTime carbon.Carbon, timeType int, timeN []int64) (carbon.Carbon, bool) {
	nextAt, b := GetTimeByTimeN2(addTime, timeType, timeN)
	if !b {
		return nextAt, b
	}
	step := 0
	for {
		//超出9000次跳出
		if step > 9000 {
			return nextAt, false
		}
		//如果满足时间，跳出
		if nextAt.Time.Unix() >= CoreFilter.GetNowTime().Unix() {
			return nextAt, b
		}
		//继续获取
		nextAt, b = GetTimeByTimeN2(nextAt, timeType, timeN)
		//下一步
		step += 1
	}
}
