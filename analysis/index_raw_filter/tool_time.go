package AnalysisIndexRawFilter

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"time"
)

// timeDifference 计算两个时间点之间的差距，根据提供的单位返回差距值
func timeDifference(t1, t2 time.Time, unit string) (float64, error) {
	duration := t2.Sub(t1)
	switch unit {
	case "sec":
		return duration.Seconds(), nil
	case "min":
		return duration.Minutes(), nil
	case "hour":
		return duration.Hours(), nil
	case "day":
		return CoreFilter.GetRound(duration.Hours() / 24), nil
	case "year":
		year1, _, _ := t1.Date()
		year2, _, _ := t2.Date()
		return float64(year2 - year1), nil
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}
}
