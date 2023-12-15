package CoreSQL2

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"time"
)

type ArgsTimeBetween struct {
	//最小时间
	// Default时间
	MinTime string `json:"minTime" check:"defaultTime"`
	//最大时间
	// Default时间
	MaxTime string `json:"maxTime" check:"defaultTime"`
}

func (t *ArgsTimeBetween) GetFields() (FieldsTimeBetween, error) {
	minAt, err := CoreFilter.GetTimeByDefault(t.MinTime)
	if err != nil {
		return FieldsTimeBetween{}, err
	}
	maxAt, err := CoreFilter.GetTimeByDefault(t.MaxTime)
	if err != nil {
		return FieldsTimeBetween{}, err
	}
	result := FieldsTimeBetween{
		MinTime: minAt,
		MaxTime: maxAt,
	}
	return result, nil
}

type FieldsTimeBetween struct {
	//最小时间
	MinTime time.Time `json:"minTime"`
	//最大时间
	MaxTime time.Time `json:"maxTime"`
}

func (t *FieldsTimeBetween) GetData() ArgsTimeBetween {
	return ArgsTimeBetween{
		MinTime: CoreFilter.GetTimeToDefaultTime(t.MinTime),
		MaxTime: CoreFilter.GetTimeToDefaultTime(t.MaxTime),
	}
}
