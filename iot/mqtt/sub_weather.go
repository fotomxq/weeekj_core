package IOTMQTT

import (
	"encoding/json"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type subWeatherData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//城市编码
	CityCode int `db:"city_code" json:"cityCode" check:"intThan0" empty:"true"`
	//查询天数
	// 1 1天 / 3 3天 / 7 7天
	DayCount int `db:"day_count" json:"dayCount" check:"intThan0" empty:"true"`
}

func subWeather(client mqtt.Client, message mqtt.Message) {
	var resultData subWeatherData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub get weather data, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub get weather data, key, ", err)
		return
	}
	if err := PushWeather(resultData.Keys.GroupMark, resultData.Keys.Code, resultData.Country, resultData.CityCode, resultData.DayCount); err != nil {
		CoreLog.MqttError("mqtt sub get weather data, key, ", err)
		return
	} else {
		//CoreLog.Info("mqtt sub get weather data, device group: ", resultData.Keys.GroupMark, ", code: ", resultData.Keys.Code, ", country: ", resultData.Country, ", city: ", resultData.CityCode)
	}
}
