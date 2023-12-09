package CoreFilter

import (
	"github.com/golang-module/carbon"
	"strings"
)

// GetAnalysisSQLTimeRange 获取统计的SQL时间范围
// analysisType 统计周期类型
// day 日统计；week 周统计；month 月统计; quarter 季度统计；year_h 半年统计；year 年统计；
// analysisAt 统计日期
// eg:
// day: 2023-06-01
// week: 2023-06-W1 / 2023-06-W2 / 2023-06-W3 / 2023-06-W4
// month: 2023-06
// quarter: 2023-Q1 / 2023-Q2 / 2023-Q3 / 2023-Q4
// year_h: 2023-H1 / 2023-H2
// year: 2023
func GetAnalysisSQLTimeRange(analysisType string, analysisAt string) (startAt carbon.Carbon, endAt carbon.Carbon) {
	defer func() {
		if err := recover(); err != nil {
			startAt = carbon.Now()
			endAt = carbon.Now()
		}
	}()
	switch analysisType {
	case "day":
		analysisAtSub := strings.Split(analysisAt, "-")
		startAt = GetNowTimeCarbon().StartOfDay().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(GetIntByStringNoErr(analysisAtSub[1])).SetDay(GetIntByStringNoErr(analysisAtSub[2]))
		endAt = startAt.EndOfDay()
	case "week":
		analysisAtSub := strings.Split(analysisAt, "-")
		switch analysisAtSub[2] {
		case "W1":
			startAt = GetNowTimeCarbon().StartOfMonth().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(GetIntByStringNoErr(analysisAtSub[1])).SetDay(1)
			endAt = startAt.AddDays(6).EndOfDay()
		case "W2":
			startAt = GetNowTimeCarbon().StartOfMonth().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(GetIntByStringNoErr(analysisAtSub[1])).SetDay(8)
			endAt = startAt.AddDays(6).EndOfDay()
		case "W3":
			startAt = GetNowTimeCarbon().StartOfMonth().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(GetIntByStringNoErr(analysisAtSub[1])).SetDay(15)
			endAt = startAt.AddDays(6).EndOfDay()
		case "W4":
			startAt = GetNowTimeCarbon().StartOfMonth().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(GetIntByStringNoErr(analysisAtSub[1])).SetDay(22)
			endAt = startAt.AddDays(6).EndOfDay()
		}
	case "month":
		analysisAtSub := strings.Split(analysisAt, "-")
		startAt = GetNowTimeCarbon().StartOfMonth().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(GetIntByStringNoErr(analysisAtSub[1]))
		endAt = startAt.EndOfMonth()
	case "quarter":
		analysisAtSub := strings.Split(analysisAt, "-")
		switch analysisAtSub[1] {
		case "Q1":
			startAt = GetNowTimeCarbon().StartOfYear().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(1).SetDay(1)
			endAt = startAt.AddMonths(3).SubDay().EndOfMonth()
		case "Q2":
			startAt = GetNowTimeCarbon().StartOfYear().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(4).SetDay(1)
			endAt = startAt.AddMonths(3).SubDay().EndOfMonth()
		case "Q3":
			startAt = GetNowTimeCarbon().StartOfYear().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(7).SetDay(1)
			endAt = startAt.AddMonths(3).SubDay().EndOfMonth()
		case "Q4":
			startAt = GetNowTimeCarbon().StartOfYear().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(10).SetDay(1)
			endAt = startAt.AddMonths(3).SubDay().EndOfMonth()
		}
	case "year_h":
		analysisAtSub := strings.Split(analysisAt, "-")
		switch analysisAtSub[1] {
		case "H1":
			startAt = GetNowTimeCarbon().StartOfYear().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(1).SetDay(1)
			endAt = startAt.AddMonths(6).SubDay()
		case "H2":
			startAt = GetNowTimeCarbon().StartOfYear().SetYear(GetIntByStringNoErr(analysisAtSub[0])).SetMonth(7).SetDay(1)
			endAt = startAt.AddMonths(6).SubDay().EndOfMonth()
		}
	case "year":
		startAt = GetNowTimeCarbon().StartOfYear().SetYear(GetIntByStringNoErr(analysisAt)).SetMonth(1).SetDay(1)
		endAt = startAt.EndOfYear()
	}
	return
}
