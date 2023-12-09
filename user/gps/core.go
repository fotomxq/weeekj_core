package UserGPS

import (
	"errors"
	"github.com/robfig/cron"
	"sync"
	"time"
)

//用户定位服务
var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
	//每隔N时间更新配置
	minSaveConfigTime = int64(time.Minute * 1)
	//最短记录距离
	userGPSRecordLimitDistance         float64
	userGPSRecordLimitDistanceLastTime time.Time
	userGPSRecordLimitDistanceLock     sync.Mutex
	//最短记录时间
	userGPSRecordLimitTime         int64
	userGPSRecordLimitTimeLastTime time.Time
	userGPSRecordLimitTimeLock     sync.Mutex
	//最短记录极限时间
	userGPSRecordLimitMaxTime         int64
	userGPSRecordLimitMaxTimeLastTime time.Time
	userGPSRecordLimitMaxTimeLock     sync.Mutex
)

//checkMapType 检查地图类型
func checkMapType(mapType int) (err error) {
	switch mapType {
	case 0:
	case 1:
	case 2:
	default:
		err = errors.New("map type is error")
		return
	}
	return
}
