package IOTTrack

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

// getIOTGPSRecordLimitDistance 最短记录距离
func getIOTGPSRecordLimitDistance() (data float64, err error) {
	if CoreFilter.GetNowTime().Unix()-iotGPSRecordLimitDistanceLastTime.Unix() < minSaveConfigTime {
		return iotGPSRecordLimitDistance, nil
	}
	data, err = BaseConfig.GetDataFloat64("IOTTrackRecordLimitDistance")
	if err == nil {
		iotGPSRecordLimitDistanceLock.Lock()
		defer iotGPSRecordLimitDistanceLock.Unlock()
		iotGPSRecordLimitDistance = data
		iotGPSRecordLimitDistanceLastTime = CoreFilter.GetNowTime()
	}
	return
}

// getIOTGPSRecordLimitTime 最短记录时间
func getIOTGPSRecordLimitTime() (data int64, err error) {
	if CoreFilter.GetNowTime().Unix()-iotGPSRecordLimitTimeLastTime.Unix() < minSaveConfigTime {
		return iotGPSRecordLimitTime, nil
	}
	data, err = BaseConfig.GetDataInt64("IOTTrackRecordLimitTime")
	if err == nil {
		iotGPSRecordLimitTimeLock.Lock()
		defer iotGPSRecordLimitTimeLock.Unlock()
		iotGPSRecordLimitTime = data
		iotGPSRecordLimitTimeLastTime = CoreFilter.GetNowTime()
	}
	return
}

// getIOTGPSRecordLimitMaxTime 最短记录极限时间
func getIOTGPSRecordLimitMaxTime() (data int64, err error) {
	if CoreFilter.GetNowTime().Unix()-iotGPSRecordLimitMaxTimeLastTime.Unix() < minSaveConfigTime {
		return iotGPSRecordLimitMaxTime, nil
	}
	data, err = BaseConfig.GetDataInt64("IOTTrackRecordLimitMaxTime")
	if err == nil {
		iotGPSRecordLimitMaxTimeLock.Lock()
		defer iotGPSRecordLimitMaxTimeLock.Unlock()
		iotGPSRecordLimitMaxTime = data
		iotGPSRecordLimitMaxTimeLastTime = CoreFilter.GetNowTime()
	}
	return
}
