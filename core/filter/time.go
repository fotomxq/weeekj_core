package CoreFilter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon"
)

// LoadTimeLocation 加载时区文件
func LoadTimeLocation(location string) *time.Location {
	l := LoadTimeLocationChild(location)
	if l == nil {
		l2 := time.FixedZone("CST", 8*60*60)
		return l2
	}
	return l
}

func LoadTimeLocationChild(location string) *time.Location {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	switch location {
	case "Asia/Shanghai":
		l, err := time.LoadLocation("Asia/Shanghai")
		if err == nil {
			return l
		}
		l2 := time.FixedZone("CST", 8*60*60)
		return l2
	}
	return time.FixedZone("CST", 0)
}

// GetNowTime 获取当前时间
func GetNowTime() time.Time {
	return time.Now().In(LoadTimeLocation("Asia/Shanghai"))
}

func GetNowTimeCarbon() carbon.Carbon {
	return GetCarbonByTime(time.Now())
}

func GetCarbonByTime(t time.Time) carbon.Carbon {
	return carbon.SetTimezone("Asia/Shanghai").CreateFromGoTime(t)
}

// GetTimeByAdd 根据推移变量，获取实际的时间结构
// 直接放对应的时间+单位即可实现
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func GetTimeByAdd(addStr string) (time.Time, error) {
	addTime, err := time.ParseDuration(addStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(LoadTimeLocation("Asia/Shanghai")).Add(addTime), nil
}

// GetTimeByTimeN 根据时间类型和时间长度，计算指定未来的时间
// timeType: 0 小时 1 天 2 周 3 月 4 年
func GetTimeByTimeN(addTime carbon.Carbon, timeType int, timeN int) time.Time {
	switch timeType {
	case 0:
		return addTime.AddHours(timeN).Time
	case 1:
		return addTime.AddDays(timeN).Time
	case 2:
		return addTime.AddWeeks(timeN).Time
	case 3:
		return addTime.AddMonths(timeN).Time
	case 4:
		return addTime.AddYears(timeN).Time
	default:
		return addTime.Time
	}
}

// GetUnixStartTime 获取unix时间戳起点
func GetUnixStartTime() time.Time {
	return time.Unix(0, 0)
}

// GetTimeBetweenAdd 获取某个ADD时间点距离现在的秒数
// 注意，如果向前偏移，则可能出现负数！
func GetTimeBetweenAdd(addStr string) (int64, error) {
	addTime, err := GetTimeByAdd(addStr)
	if err != nil {
		return 0, errors.New("cannot get add time, " + err.Error())
	}
	return addTime.Unix() - GetNowTime().Unix(), nil
}

// GetTimeBy30DayList 获取最近30天的列表
func GetTimeBy30DayList() []string {
	var dataList []string
	step := 29
	for {
		if step < 0 {
			break
		}
		dataList = append(dataList, GetNowTime().AddDate(0, 0, 0-step).Format("20060102"))
		step -= 1
	}
	return dataList
}

// GetTimeByISO 将ISO时间转为time
func GetTimeByISO(newTime string) (timeAt time.Time, err error) {
	timeAt, err = time.Parse(time.RFC3339Nano, newTime)
	if err == nil {
		timeAt = timeAt.In(LoadTimeLocation("Asia/Shanghai"))
	}
	return
}

// GetISOByTime 将时间转为ISO格式输出
func GetISOByTime(timeAt time.Time) (newTime string) {
	newTime = timeAt.In(LoadTimeLocation("Asia/Shanghai")).Format(time.RFC3339Nano)
	return
}

// GetTimeByDefault 将普通时间格式转为时间结构
func GetTimeByDefault(str string) (timeAt time.Time, err error) {
	if !strings.Contains(str, ":") {
		str = fmt.Sprint(str, " 00:00:00")
	}
	if strings.Contains(str, "-") {
		timeAt, err = time.ParseInLocation("2006-01-02 15:04:05", str, LoadTimeLocation("Asia/Shanghai"))
		if err != nil {
			timeAt, err = time.ParseInLocation("2006-1-2 15:04:05", str, LoadTimeLocation("Asia/Shanghai"))
		}
	} else {
		timeAt, err = time.ParseInLocation("2006/01/02 15:04:05", str, LoadTimeLocation("Asia/Shanghai"))
		if err != nil {
			timeAt, err = time.ParseInLocation("2006/1/2 15:04:05", str, LoadTimeLocation("Asia/Shanghai"))
		}
	}
	return
}

// GetTimeByDefaultNoErr 将普通时间格式转为时间结构(不反馈错误)
func GetTimeByDefaultNoErr(str string) (timeAt time.Time) {
	var err error
	if !strings.Contains(str, ":") {
		str = fmt.Sprint(str, " 00:00:00")
	}
	if strings.Contains(str, "-") {
		timeAt, err = time.ParseInLocation("2006-01-02 15:04:05", str, LoadTimeLocation("Asia/Shanghai"))
		if err != nil {
			timeAt, err = time.ParseInLocation("2006-1-2 15:04:05", str, LoadTimeLocation("Asia/Shanghai"))
		}
	} else {
		timeAt, err = time.ParseInLocation("2006/01/02 15:04:05", str, LoadTimeLocation("Asia/Shanghai"))
		if err != nil {
			timeAt, err = time.ParseInLocation("2006/1/2 15:04:05", str, LoadTimeLocation("Asia/Shanghai"))
		}
	}
	return
}

// GetTimeCarbonByDefault 获取carbon时间
func GetTimeCarbonByDefault(str string) (timeAt carbon.Carbon, err error) {
	var t time.Time
	t, err = GetTimeByDefault(str)
	if err != nil {
		return
	}
	timeAt = GetCarbonByTime(t)
	return
}

// GetTimeToDefaultTime 将时间格式输出为标准时间结构
func GetTimeToDefaultTime(timeAt time.Time) (str string) {
	return timeAt.Format("2006-01-02 15:04:05")
}

// GetTimeToDefaultDate 将时间格式输出为日期结构
func GetTimeToDefaultDate(timeAt time.Time) (str string) {
	return timeAt.Format("2006-01-02")
}

// GetTimeByDefaultTime 获取两个标准时间的时间
func GetTimeByDefaultTime(a, b string) (time.Time, time.Time, error) {
	aT, err := GetTimeByDefault(a)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	bT, err := GetTimeByDefault(b)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return aT, bT, nil
}

// GetTimeByUnix 获取unix时间戳的时间结构体
func GetTimeByUnix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec).In(LoadTimeLocation("Asia/Shanghai"))
}

// CheckHaveTime 判断是否具备时间
// 和sql模块内部旧的处理判断方法一致，用于检查是否具备时间，例如deleteAt是否已经删除等
func CheckHaveTime(d time.Time) bool {
	return d.Unix() > 1000000
}

// CheckTimeAfterNow 判断时间是否大于当前时间
func CheckTimeAfterNow(d time.Time) bool {
	return d.After(GetNowTime())
}

// GetWeekOfMonthByTimeCarbon 获取当前时间是本月第几周
func GetWeekOfMonthByTimeCarbon() string {
	//当前时间
	date := time.Now()
	// 获取该日期所在月份的第一天
	firstOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	// 计算该日期是本月的第几周
	days := date.Sub(firstOfMonth).Hours() / 24
	weekNumber := int((days + float64(firstOfMonth.Weekday()) + 5) / 7)
	nowAtCarbon := GetNowTimeCarbon().Format("2006-01")
	return fmt.Sprint(nowAtCarbon, "-W", weekNumber)
}

// GetQuarterByTimeCarbon 获取当前时间的季度日期
func GetQuarterByTimeCarbon() string {
	d := GetNowTimeCarbon()
	return fmt.Sprint(d.Year(), "-Q", d.Month()/4+1)
}

// GetHalfYearByTimeCarbon 获取当前时间是上半年还是下半年
func GetHalfYearByTimeCarbon() string {
	d := GetNowTimeCarbon()
	if d.Quarter() <= 2 {
		return fmt.Sprint(d.Year(), "-H1")
	}
	return fmt.Sprint(d.Year(), "-H2")
}

// 根据年月获取当月的第一天时间
// 例如：2016-07
func GetFirstDayByYearMonth(yearMonth string) (firstDay string) {
	// 解析年月字符串
	year, month, err := parseYearMonth(yearMonth)
	if err != nil {
		return ""
	}
	// 获取该月份的第一天
	firstDayTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	firstDay = firstDayTime.Format("2006-01-02")
	return firstDay
}

// 根据年月获取当月的最后一天时间
func GetLastDayByYearMonth(yearMonth string) (lastDay string) {
	// 解析年月字符串
	year, month, err := parseYearMonth(yearMonth)
	if err != nil {
		return ""
	}
	// 获取该月份的第一天
	firstDayTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// 获取下一个月的第一天
	nextMonthTime := firstDayTime.AddDate(0, 1, 0)
	// 获取这个月的最后一天
	lastDayTime := nextMonthTime.Add(-24 * time.Hour)
	lastDay = lastDayTime.Format("2006-01-02")

	return lastDay
}

// 解析年月字符串为年份和月份
func parseYearMonth(yearMonth string) (year int, month int, err error) {
	if len(yearMonth) != 7 || yearMonth[4] != '-' {
		return 0, 0, fmt.Errorf("invalid yearMonth format, should be YYYY-MM")
	}
	// 解析年份
	year, err = strconv.Atoi(yearMonth[:4])
	if err != nil {
		return 0, 0, err
	}
	// 解析月份
	month, err = strconv.Atoi(yearMonth[5:])
	if err != nil {
		return 0, 0, err
	}
	// 检查月份范围
	if month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("month should be between 1 and 12")
	}
	return year, month, nil
}

// GetTimeByParseExcelDate 解析 Excel 格式的日期字符串，并返回 time.Time 类型
func GetTimeByParseExcelDate(dateStr string) time.Time {
	// 定义可能的时间格式
	formats := []string{
		"2006/01/02 15:04:05", // 格式：2020/08/27 14:30:00
		"1/2/06 3:04:05 PM",   // 格式：8/27/20 2:30:00 PM
		"2006/01/02",          // 格式：2020/08/27
		"1/2/06",              // 格式：8/27/20
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, dateStr)
		if err == nil {
			return parsedTime // 成功解析，返回时间和成功标志
		}
	}

	return time.Time{} // 解析失败，返回默认时间和失败标志
}
