package UserGPS

import (
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

// getUserGPSRecordLimitDistance 最短记录距离
func getUserGPSRecordLimitDistance() (data float64, err error) {
	if CoreFilter.GetNowTime().Unix()-userGPSRecordLimitDistanceLastTime.Unix() < minSaveConfigTime {
		return userGPSRecordLimitDistance, nil
	}
	data, err = BaseConfig.GetDataFloat64("UserGPSRecordLimitDistance")
	if err == nil {
		userGPSRecordLimitDistanceLock.Lock()
		defer userGPSRecordLimitDistanceLock.Unlock()
		userGPSRecordLimitDistance = data
		userGPSRecordLimitDistanceLastTime = CoreFilter.GetNowTime()
	}
	return
}

// getUserGPSRecordLimitTime 最短记录时间
func getUserGPSRecordLimitTime() (data int64, err error) {
	if CoreFilter.GetNowTime().Unix()-userGPSRecordLimitTimeLastTime.Unix() < minSaveConfigTime {
		return userGPSRecordLimitTime, nil
	}
	data, err = BaseConfig.GetDataInt64("UserGPSRecordLimitTime")
	if err == nil {
		userGPSRecordLimitTimeLock.Lock()
		defer userGPSRecordLimitTimeLock.Unlock()
		userGPSRecordLimitTime = data
		userGPSRecordLimitTimeLastTime = CoreFilter.GetNowTime()
	}
	return
}

// getUserGPSRecordLimitMaxTime 最短记录极限时间
func getUserGPSRecordLimitMaxTime() (data int64, err error) {
	if CoreFilter.GetNowTime().Unix()-userGPSRecordLimitMaxTimeLastTime.Unix() < minSaveConfigTime {
		return userGPSRecordLimitMaxTime, nil
	}
	data, err = BaseConfig.GetDataInt64("UserGPSRecordLimitMaxTime")
	if err == nil {
		userGPSRecordLimitMaxTimeLock.Lock()
		defer userGPSRecordLimitMaxTimeLock.Unlock()
		userGPSRecordLimitMaxTime = data
		userGPSRecordLimitMaxTimeLastTime = CoreFilter.GetNowTime()
	}
	return
}
