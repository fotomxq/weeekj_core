package MapRoom

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "gitee.com/weeekj/weeekj_core/v5/core/sql/analysis"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//房间设备数据统计
// 该模块通过绑定关系获取所有绑定设备，然后将设备采集信息进行汇总二次加工处理

// ArgsGetSensorList 获取统计数据列表参数
type ArgsGetSensorList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
}

// GetSensorList 获取统计数据列表参数
func GetSensorList(args *ArgsGetSensorList) (dataList []FieldsSensor, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.RoomID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "room_id = :room_id"
		maps["room_id"] = args.RoomID
	}
	if args.DeviceID > -1 {
		if where != "" {
			where = where + " AND "
		}
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
	tableName := "map_room_iot_sensor"
	if args.IsHistory {
		tableName = "map_room_iot_sensor_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, room_id, device_id, mark, data, data_f, data_s FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetSensorMore 获取一组房间最新的数据参数
type ArgsGetSensorMore struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//房间ID列
	RoomIDs pq.Int64Array `db:"room_ids" json:"roomIDs" check:"ids"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
}

// GetSensorMore 获取一组房间最新的数据
func GetSensorMore(args *ArgsGetSensorMore) (dataList []FieldsSensor, err error) {
	if len(args.RoomIDs) > 30 {
		err = errors.New("too many")
		return
	}
	for _, v := range args.RoomIDs {
		var vData FieldsSensor
		err = Router2SystemConfig.MainDB.Get(&vData, "SELECT id, create_at, device_id, org_id, room_id, mark, data, data_f, data_s FROM map_room_iot_sensor WHERE room_id = $1 AND mark = $2 AND org_id = $3 ORDER BY id DESC LIMIT 1", v, args.Mark, args.OrgID)
		if err != nil || vData.ID < 1 {
			err = nil
			continue
		}
		dataList = append(dataList, vData)
	}
	if len(dataList) < 1 {
		err = errors.New("data is empty")
	}
	return
}

// ArgsGetSensorMoreMarks 获取一组房间最新的数据参数
type ArgsGetSensorMoreMarks struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//房间ID列
	RoomIDs pq.Int64Array `db:"room_ids" json:"roomIDs" check:"ids"`
	//数据标识码
	Marks pq.StringArray `db:"marks" json:"marks" check:"marks"`
}

// GetSensorMoreMarks 获取一组房间最新的数据
// 最大上限反馈600条数据，最多不能超出30个房间
func GetSensorMoreMarks(args *ArgsGetSensorMoreMarks) (dataList []FieldsSensor, err error) {
	if len(args.RoomIDs) > 30 {
		err = errors.New("too many")
		return
	}
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, device_id, org_id, room_id, mark, data, data_f, data_s FROM map_room_iot_sensor WHERE room_id = ANY($1) AND mark = ANY($2) AND org_id = $3 ORDER BY id DESC LIMIT 600", args.RoomIDs, args.Marks, args.OrgID)
	if err != nil {
		return
	}
	if len(dataList) < 1 {
		err = errors.New("data is empty")
	}
	return
}

// ArgsGetSensorAnalysis 获取统计数据SUM参数
type ArgsGetSensorAnalysis struct {
	//查询时间范围
	TimeBetween CoreSQLTime.FieldsCoreTime `json:"timeBetween"`
	//结构方式
	// year / month / day / hour
	TimeType string `json:"timeType"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
}

type DataGetSensorAnalysis struct {
	//时间
	DayTime string `db:"d" json:"dayTime"`
	//数据
	Data  int64   `db:"data" json:"data"`
	DataF float64 `db:"data_f" json:"dataF"`
}

// GetSensorAnalysis 获取统计数据SUM
func GetSensorAnalysis(args *ArgsGetSensorAnalysis) (dataList []DataGetSensorAnalysis, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.RoomID > -1 {
		where = where + " AND room_id = :room_id"
		maps["room_id"] = args.RoomID
	}
	if args.DeviceID > -1 {
		if where != "" {
			where = where + " AND "
		}
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
	tableName := "map_room_iot_sensor"
	if args.IsHistory {
		tableName = "map_room_iot_sensor_history"
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

// GetSensorAnalysisAvg 获取统计数据AVG
func GetSensorAnalysisAvg(args *ArgsGetSensorAnalysis) (dataList []DataGetSensorAnalysis, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.RoomID > -1 {
		where = where + " AND room_id = :room_id"
		maps["room_id"] = args.RoomID
	}
	if args.DeviceID > -1 {
		if where != "" {
			where = where + " AND "
		}
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
	tableName := "map_room_iot_sensor"
	if args.IsHistory {
		tableName = "map_room_iot_sensor_history"
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

// GetSensorAnalysisMax 获取统计数据MAX
func GetSensorAnalysisMax(args *ArgsGetSensorAnalysis) (dataList []DataGetSensorAnalysis, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.RoomID > -1 {
		where = where + " AND room_id = :room_id"
		maps["room_id"] = args.RoomID
	}
	if args.DeviceID > -1 {
		if where != "" {
			where = where + " AND "
		}
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
	tableName := "map_room_iot_sensor"
	if args.IsHistory {
		tableName = "map_room_iot_sensor_history"
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

// GetSensorAnalysisMin 获取统计数据MIN
func GetSensorAnalysisMin(args *ArgsGetSensorAnalysis) (dataList []DataGetSensorAnalysis, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.RoomID > -1 {
		where = where + " AND room_id = :room_id"
		maps["room_id"] = args.RoomID
	}
	if args.DeviceID > -1 {
		if where != "" {
			where = where + " AND "
		}
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
	tableName := "map_room_iot_sensor"
	if args.IsHistory {
		tableName = "map_room_iot_sensor_history"
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

// ArgsDeleteSensorClear 删除统计数据参数
type ArgsDeleteSensorClear struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
}

// DeleteSensorClear 清理指定数据
func DeleteSensorClear(args *ArgsDeleteSensorClear) (err error) {
	if args.OrgID < 1 && args.RoomID < 1 && args.DeviceID < 1 && args.Mark == "" {
		err = errors.New("all args empty")
		return
	}
	if args.IsHistory {
		_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "map_room_iot_sensor_history", "(:org_id < 1 OR org_id = :org_id) AND (:room_id < 1 OR room_id = :room_id) AND (:device_id < 1 OR device_id = :device_id) AND (:mark = '' OR mark = :mark)", args)
	} else {
		_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "map_room_iot_sensor", "(:org_id < 1 OR org_id = :org_id) AND (:room_id < 1 OR room_id = :room_id) AND (:device_id < 1 OR device_id = :device_id) AND (:mark = '' OR mark = :mark)", args)
	}
	return
}
