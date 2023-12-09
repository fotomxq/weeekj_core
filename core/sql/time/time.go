package CoreSQLTime

import (
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"time"
)

// FieldsCoreTime 时间周期选择器和处理器
// Deprecated
type FieldsCoreTime struct {
	//最小时间
	MinTime time.Time `json:"minTime"`
	//最大时间
	MaxTime time.Time `json:"maxTime"`
}

// DataCoreTime 外部参数类型
// Deprecated
type DataCoreTime struct {
	//最小时间
	// ISO时间
	MinTime string `json:"minTime" check:"isoTime"`
	//最大时间
	// ISO时间
	MaxTime string `json:"maxTime" check:"isoTime"`
}

// GetBetweenByISO 将ISO时间转为标准时间
// Deprecated
func GetBetweenByISO(args DataCoreTime) (data FieldsCoreTime, err error) {
	data.MinTime, err = CoreFilter.GetTimeByISO(args.MinTime)
	if err != nil {
		return
	}
	data.MaxTime, err = CoreFilter.GetTimeByISO(args.MaxTime)
	return
}

// GetBetweenByTime 获取时间范围的mgo
// Deprecated
func GetBetweenByTime(fieldsName string, args FieldsCoreTime, maps map[string]interface{}) (query string, newMaps map[string]interface{}) {
	if args.MinTime.Unix() < 1 || args.MaxTime.Unix() < 1 {
		return "", maps
	}
	query = fmt.Sprint(fieldsName, " >= :", fieldsName, "_min_time AND ", fieldsName, " <= :", fieldsName, "_max_time")
	maps[fieldsName+"_min_time"] = args.MinTime
	maps[fieldsName+"_max_time"] = args.MaxTime
	return query, maps
}

// GetBetweenByTimeAnd 带有query处理的方法
// Deprecated
func GetBetweenByTimeAnd(fieldsName string, args FieldsCoreTime, query string, maps map[string]interface{}) (newQuery string, newMaps map[string]interface{}) {
	var newWhere string
	newWhere, newMaps = GetBetweenByTime(fieldsName, args, maps)
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
