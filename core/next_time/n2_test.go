package CoreNextTime

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"testing"
)

func TestGetTimeByTimeN2Now(t *testing.T) {
	timeAt, b := GetTimeByTimeN2Now(CoreFilter.GetNowTimeCarbon().SubWeek(), 0, []int64{})
	t.Log(timeAt, ", ", b)
	timeAt, b = GetTimeByTimeN2Now(CoreFilter.GetNowTimeCarbon().SubWeek().SubHour(), 0, []int64{})
	t.Log(timeAt, ", ", b)
	timeAt, b = GetTimeByTimeN2Now(timeAt, 0, []int64{})
	t.Log(timeAt, ", ", b)
}
