package CoreSQL

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"time"
)

// CheckTimeHaveData 检查sql语句的delete类型时间是否存在数据
func CheckTimeHaveData(d time.Time) bool {
	return d.Unix() > 1000000
}

// CheckTimeThanNow 检查特定时间是否存在或是否满足当前时间？
func CheckTimeThanNow(d time.Time) bool {
	if d.Unix() < 1000000 || d.Unix() >= CoreFilter.GetNowTime().Unix() {
		return true
	}
	return false
}
