package IOTTrack

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"math"
	"sync"
	"time"
)

var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
	//每隔N时间更新配置
	minSaveConfigTime = int64(time.Minute * 1)
	//最短记录距离
	iotGPSRecordLimitDistance         float64
	iotGPSRecordLimitDistanceLastTime time.Time
	iotGPSRecordLimitDistanceLock     sync.Mutex
	//最短记录时间
	iotGPSRecordLimitTime         int64
	iotGPSRecordLimitTimeLastTime time.Time
	iotGPSRecordLimitTimeLock     sync.Mutex
	//最短记录极限时间
	iotGPSRecordLimitMaxTime         int64
	iotGPSRecordLimitMaxTimeLastTime time.Time
	iotGPSRecordLimitMaxTimeLock     sync.Mutex
)

// ArgsGetList 获取定位追踪列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//所属设备
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
}

// GetList 获取定位追踪列表
func GetList(args *ArgsGetList) (dataList []FieldsTrack, dataCount int64, err error) {
	if err = checkMapType(args.MapType); err != nil {
		return
	}
	var where string
	maps := map[string]interface{}{}
	if args.DeviceID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	if args.MapType > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "map_type = :map_type"
		maps["map_type"] = args.MapType
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"iot_track",
		"id",
		"SELECT id, create_at, device_id, map_type, longitude, latitude, station_info FROM iot_track WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetLast 获取最新的定位参数
type ArgsGetLast struct {
	//所属设备
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
}

// GetLast 获取最新的定位
func GetLast(args *ArgsGetLast) (data FieldsTrack, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, device_id, map_type, longitude, latitude FROM iot_track WHERE device_id = $1 ORDER BY id DESC LIMIT 1", args.DeviceID)
	return
}

// ArgsCreate 添加新的定位参数
type ArgsCreate struct {
	//所属设备
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//地图制式
	// 0 / 1 / 2 / 3
	// WGS-84 / GCJ-02 / BD-09 / 2000中国大地坐标系
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//基站信息
	StationInfo string `db:"station_info" json:"stationInfo"`
}

// Create 添加新的定位
func Create(args *ArgsCreate) (err error) {
	//检查地图类型
	if err = checkMapType(args.MapType); err != nil {
		return
	}
	//检查该记录是否可以写入？
	var data FieldsTrack
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT create_at, map_type, longitude, latitude FROM iot_track WHERE device_id = $1 ORDER BY id DESC LIMIT 1", args.DeviceID)
	//如果存在数据
	if err == nil {
		//当前时间
		nowTime := CoreFilter.GetNowTime()
		//极限记录时间，如果符合则跳过位移判断
		var recordLimitMaxTime int64
		recordLimitMaxTime, err = getIOTGPSRecordLimitMaxTime()
		if nowTime.Unix()-data.CreateAt.Unix() < recordLimitMaxTime {
			//检查位移是否符合记录条件
			var recordLimitDistance float64
			recordLimitDistance, err = getIOTGPSRecordLimitDistance()
			if err != nil {
				err = errors.New("get config recordLimitDistance, " + err.Error())
				return
			}
			if math.Abs(args.Longitude-data.Longitude) < recordLimitDistance || math.Abs(args.Latitude-data.Latitude) < recordLimitDistance {
				return
			}
		}
		//检查时间
		var recordLimitTime int64
		recordLimitTime, err = getIOTGPSRecordLimitTime()
		if err != nil {
			err = errors.New("get config recordLimitTime, " + err.Error())
			return
		}
		if nowTime.Unix()-data.CreateAt.Unix() < recordLimitTime {
			return
		}
	}
	//创建新的记录
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO iot_track (device_id, map_type, longitude, latitude, station_info) VALUES (:device_id, :map_type,:longitude,:latitude,:station_info)", args)
	return
}

// ArgsDeleteByDevice 删除指定用户的所有定位参数
type ArgsDeleteByDevice struct {
	//所属设备
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
}

// DeleteByDevice 删除指定用户的所有定位
func DeleteByDevice(args *ArgsDeleteByDevice) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_track", "device_id = :device_id", args)
	return
}

// checkMapType 检查地图类型
func checkMapType(mapType int) (err error) {
	switch mapType {
	case 0:
	case 1:
	case 2:
	case 3:
	default:
		err = errors.New("map type is error")
		return
	}
	return
}
