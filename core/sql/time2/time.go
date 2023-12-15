package CoreSQLTime2

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	"time"
)

// FieldsCoreTime 时间周期选择器和处理器
type FieldsCoreTime struct {
	//最小时间
	MinTime time.Time `json:"minTime"`
	//最大时间
	MaxTime time.Time `json:"maxTime"`
}

// DataCoreTime 外部参数类型
type DataCoreTime struct {
	//最小时间
	// Default时间
	MinTime string `json:"minTime" check:"defaultTime"`
	//最大时间
	// Default时间
	MaxTime string `json:"maxTime" check:"defaultTime"`
}

// GetBetweenByISO 将ISO时间转为标准时间
func (t *DataCoreTime) GetBetweenByISO() (data FieldsCoreTime, err error) {
	data.MinTime, err = CoreFilter.GetTimeByDefault(t.MinTime)
	if err != nil {
		return
	}
	data.MaxTime, err = CoreFilter.GetTimeByDefault(t.MaxTime)
	if err != nil {
		return
	}
	return
}

func (t *DataCoreTime) GetOldTime() (data CoreSQLTime.DataCoreTime, err error) {
	var minTime, maxTime time.Time
	minTime, err = CoreFilter.GetTimeByDefault(t.MinTime)
	if err != nil {
		return
	}
	maxTime, err = CoreFilter.GetTimeByDefault(t.MaxTime)
	if err != nil {
		return
	}
	data = CoreSQLTime.DataCoreTime{
		MinTime: CoreFilter.GetISOByTime(minTime),
		MaxTime: CoreFilter.GetISOByTime(maxTime),
	}
	return
}

// GetBetweenByTime 获取时间范围的mgo
func (t *FieldsCoreTime) GetBetweenByTime(fieldsName string, maps map[string]interface{}) (query string, newMaps map[string]interface{}) {
	if t.MinTime.Unix() < 1 || t.MaxTime.Unix() < 1 {
		return "", maps
	}
	query = fmt.Sprint(fieldsName, " >= :", fieldsName, "_min_time AND ", fieldsName, " <= :", fieldsName, "_max_time")
	maps[fieldsName+"_min_time"] = t.MinTime
	maps[fieldsName+"_max_time"] = t.MaxTime
	return query, maps
}

// GetBetweenByTimeAnd 带有query处理的方法
func (t *FieldsCoreTime) GetBetweenByTimeAnd(fieldsName string, query string, maps map[string]interface{}) (newQuery string, newMaps map[string]interface{}) {
	var newWhere string
	newWhere, newMaps = t.GetBetweenByTime(fieldsName, maps)
	if newWhere != "" {
		if query == "" {
			newQuery = newWhere
		} else {
			newQuery = query + " AND " + newWhere
		}
		return newQuery, maps
	}
	return query, maps
}

// GetOldFields 获取旧的结构方式
func (t *FieldsCoreTime) GetOldFields() CoreSQLTime.FieldsCoreTime {
	return CoreSQLTime.FieldsCoreTime{
		MinTime: t.MinTime,
		MaxTime: t.MaxTime,
	}
}

type DataOldCoreTime struct {
	//最小时间
	// Default时间
	MinTime string `json:"minTime" check:"isoTime"`
	//最大时间
	// Default时间
	MaxTime string `json:"maxTime" check:"isoTime"`
}

func (t *DataOldCoreTime) GetBetweenByISO() (data DataCoreTime) {
	minTime, err := CoreFilter.GetTimeByISO(t.MinTime)
	if err != nil {
		return
	}
	maxTime, err := CoreFilter.GetTimeByISO(t.MaxTime)
	if err != nil {
		return
	}
	data = DataCoreTime{
		MinTime: CoreFilter.GetTimeToDefaultTime(minTime),
		MaxTime: CoreFilter.GetTimeToDefaultTime(maxTime),
	}
	return
}
