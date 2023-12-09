package IOTSensor

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsSensor, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.DeviceID > -1 {
		where = where + "device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_sensor"
	if args.IsHistory {
		tableName = "iot_sensor_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, device_id, mark, data, data_f, data_s FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetListTime 获取时间之间的数据集合参数
type ArgsGetListTime struct {
	//时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
}

// GetListTime 获取时间之间的数据集合
func GetListTime(args *ArgsGetListTime) (dataList []FieldsSensor, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, device_id, mark, data, data_f, data_s FROM iot_sensor WHERE create_at >= $1 AND create_at <= $2 AND ($3 < 1 OR device_id = $3) AND ($4 = '' OR mark = $4)", args.TimeBetween.MinTime, args.TimeBetween.MaxTime, args.DeviceID, args.Mark)
	return
}

// ArgsGetAnalysis 获取统计数据结构参数
type ArgsGetAnalysis struct {
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType" check:"mark"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
}

type DataGetAnalysis struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//数据
	Data  int64   `db:"data" json:"data"`
	DataF float64 `db:"data_f" json:"dataF"`
}

// GetAnalysis 获取统计数据结构
func GetAnalysis(args *ArgsGetAnalysis) (dataList []DataGetAnalysis, err error) {
	where := "mark = :mark AND device_id = :device_id"
	maps := map[string]interface{}{
		"mark":      args.Mark,
		"device_id": args.DeviceID,
	}
	tableName := "iot_sensor"
	if args.IsHistory {
		tableName = "iot_sensor_history"
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", SUM(data) as data, SUM(data_f) as data_f FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// GetAnalysisAvg 获取统计数据结构平均值
// 只支持float计算
func GetAnalysisAvg(args *ArgsGetAnalysis) (dataList []DataGetAnalysis, err error) {
	where := "mark = :mark AND device_id = :device_id"
	maps := map[string]interface{}{
		"mark":      args.Mark,
		"device_id": args.DeviceID,
	}
	tableName := "iot_sensor"
	if args.IsHistory {
		tableName = "iot_sensor_history"
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", AVG(data_f) as data_f FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

func GetAnalysisMax(args *ArgsGetAnalysis) (dataList []DataGetAnalysis, err error) {
	where := "mark = :mark AND device_id = :device_id"
	maps := map[string]interface{}{
		"mark":      args.Mark,
		"device_id": args.DeviceID,
	}
	tableName := "iot_sensor"
	if args.IsHistory {
		tableName = "iot_sensor_history"
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", MAX(data) as data, MAX(data_f) as data_f FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

func GetAnalysisMin(args *ArgsGetAnalysis) (dataList []DataGetAnalysis, err error) {
	where := "mark = :mark AND device_id = :device_id"
	maps := map[string]interface{}{
		"mark":      args.Mark,
		"device_id": args.DeviceID,
	}
	tableName := "iot_sensor"
	if args.IsHistory {
		tableName = "iot_sensor_history"
	}
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", args.TimeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", MIN(data) as data, MIN(data_f) as data_f FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	return
}

// ArgsCreate 添加数据参数
type ArgsCreate struct {
	//创建时间
	// 如果给空，则以当前时间为主
	// IOS时间
	CreateAt string `db:"create_at" json:"createAt"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//数据
	Data  int64   `db:"data" json:"data"`
	DataF float64 `db:"data_f" json:"dataF"`
	DataS string  `db:"data_s" json:"dataS"`
}

// Create 添加数据
func Create(args *ArgsCreate) (err error) {
	if b := blocker.Check(args.DeviceID, args.Mark, fmt.Sprint(args.Data, args.DataF, args.DataS)); !b {
		//err = errors.New("")
		//不反馈错误信息，外部可认为保存成功
		return
	}
	appendAt := CoreFilter.GetNowTime()
	if args.CreateAt != "" {
		appendAt, err = CoreFilter.GetTimeByISO(args.CreateAt)
		if err != nil {
			return
		}
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO iot_sensor (create_at, device_id, mark, data, data_f, data_s) VALUES (:create_at,:device_id,:mark,:data,:data_f,:data_s)", map[string]interface{}{
		"create_at": appendAt,
		"device_id": args.DeviceID,
		"mark":      args.Mark,
		"data":      args.Data,
		"data_f":    args.DataF,
		"data_s":    args.DataS,
	})
	return
}

// ArgsDeleteClear 清理指定数据参数
// 注意，如果全部留空将失败
type ArgsDeleteClear struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
}

// DeleteClear 清理指定数据
func DeleteClear(args *ArgsDeleteClear) (err error) {
	if args.DeviceID < 1 && args.Mark == "" {
		err = errors.New("device id or mark empty")
		return
	}
	if args.IsHistory {
		_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_sensor_history", "(:device_id < 1 OR device_id = :device_id) AND (:mark = '' OR mark = :mark)", args)
	} else {
		_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_sensor", "(:device_id < 1 OR device_id = :device_id) AND (:mark = '' OR mark = :mark)", args)
	}
	return
}
