package CoreNextTime

import (
	ToolsHolidaySeason "gitee.com/weeekj/weeekj_core/v5/tools/holiday_season"
	"github.com/golang-module/carbon"
)

// MakeBeforeAt 生成上一轮日期
func MakeBeforeAt(timeType int, timeN []int64, skipHoliday bool, nextAt carbon.Carbon) (newTime carbon.Carbon, needDeleteConfig bool, b bool) {
	//内部循环限制
	step := 0
	limit := 60
	//选择日期
	switch timeType {
	case 0:
		//每天重复
		for {
			if step > limit {
				break
			}
			step += 1
			nextAt = nextAt.SubDay()
			//如果不是节假日，则跳出
			if skipHoliday && !ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
				DateAt: nextAt.Time,
			}) {
				continue
			}
			break
		}
	case 1:
		//每周重复
		for {
			if step > limit {
				break
			}
			step += 1
			nextAt = nextAt.SubWeek()
			//如果不是节假日，则跳出
			if skipHoliday && !ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
				DateAt: nextAt.Time,
			}) {
				continue
			}
			break
		}
	case 2:
		//每月重复
		for {
			if step > limit {
				break
			}
			step += 1
			nextAt = nextAt.SubMonth()
			//如果不是节假日，则跳出
			if skipHoliday && !ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
				DateAt: nextAt.Time,
			}) {
				continue
			}
			break
		}
	case 3:
		//临时一次
		// 标记删除配置
		needDeleteConfig = true
		break
	case 4:
		//每隔N天重复
		for {
			if step > limit {
				break
			}
			step += 1
			if len(timeN) > 0 {
				nextAt = nextAt.SubDays(int(timeN[0]))
			} else {
				return
			}
			//如果不是节假日，则跳出
			if skipHoliday && !ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
				DateAt: nextAt.Time,
			}) {
				continue
			}
			break
		}
		break
	case 5:
		//每隔N周重复
		for {
			if step > limit {
				break
			}
			step += 1
			if len(timeN) > 0 {
				nextAt = nextAt.SubWeeks(int(timeN[0]))
			} else {
				return
			}
			//如果不是节假日，则跳出
			if skipHoliday && !ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
				DateAt: nextAt.Time,
			}) {
				continue
			}
			break
		}
		break
	case 6:
		//每隔N月重复
		for {
			if step > limit {
				break
			}
			step += 1
			if len(timeN) > 0 {
				nextAt = nextAt.SubMonths(int(timeN[0]))
			} else {
				return
			}
			//如果不是节假日，则跳出
			if skipHoliday && !ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
				DateAt: nextAt.Time,
			}) {
				continue
			}
			break
		}
		break
	case 7:
		//指定星期重复
		if len(timeN) > 0 {
			isOK := false
			for _, vTimeN := range timeN {
				//最多重复7次
				step := 0
				for {
					step += 1
					//推移一天
					nextAt = nextAt.SubDay()
					//如果符合，则考虑安排
					if nextAt.DayOfWeek() == int(vTimeN) {
						isOK = true
						break
					}
					//如果不是节假日，则跳出
					if skipHoliday && !ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
						DateAt: nextAt.Time,
					}) {
						continue
					}
					if step > 7 {
						return
					}
				}
				if isOK {
					break
				}
			}
			if isOK {
				break
			}
		} else {
			return
		}
	case 8:
		//每小时重复
		for {
			if step > limit {
				break
			}
			step += 1
			if len(timeN) > 0 {
				nextAt = nextAt.SubHours(int(timeN[0]))
			} else {
				break
			}
			//如果不是节假日，则跳出
			if skipHoliday && !ToolsHolidaySeason.CheckIsWork(&ToolsHolidaySeason.ArgsCheckIsWork{
				DateAt: nextAt.Time,
			}) {
				continue
			}
			break
		}
	default:
		//无法识别，跳出
		return
	}
	return nextAt, needDeleteConfig, true
}
