package IOTMQTT

import (
	"encoding/json"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	IOTTrack "gitee.com/weeekj/weeekj_core/v5/iot/track"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type subTrackData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//设备ID
	ID int64 `json:"id" check:"id"`
	//地图制式
	// 0 / 1 / 2 / 3
	// WGS-84 / GCJ-02 / BD-09 / 2000中国大地坐标系
	MapType int `db:"map_type" json:"mapType"`
	//定位信息
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	//基站信息
	StationInfo string `json:"stationInfo"`
}

func subTrack(client mqtt.Client, message mqtt.Message) {
	var resultData subTrackData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub create device track, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub create device track, key, ", err)
		return
	}
	//检查设备是否存在
	if err := IOTDevice.CheckDeviceCode(&IOTDevice.ArgsCheckDeviceCode{
		DeviceID:  resultData.ID,
		GroupMark: resultData.Keys.GroupMark,
		Code:      resultData.Keys.Code,
	}); err != nil {
		CoreLog.MqttError("mqtt sub create device track, device not exist, ", err)
		return
	}
	if err := IOTTrack.Create(&IOTTrack.ArgsCreate{
		DeviceID:    resultData.ID,
		MapType:     resultData.MapType,
		Longitude:   resultData.Longitude,
		Latitude:    resultData.Latitude,
		StationInfo: resultData.StationInfo,
	}); err != nil {
		CoreLog.MqttError("mqtt sub create device track, create, ", err)
	}
}
