package IOTMQTT

import (
	"encoding/json"
	"errors"
	"fmt"
	ToolsWeather "gitee.com/weeekj/weeekj_core/v5/tools/weather"
)

// PushWeather 推送给设备天气预报信息
func PushWeather(groupMark, deviceCode string, country int, city int, dataCount int) (err error) {
	//主题地址
	topic := fmt.Sprint("weather/group/", groupMark, "/code/", deviceCode)
	//获取和打包数据
	var data ToolsWeather.DataGetWeather
	data, err = ToolsWeather.GetWeather(&ToolsWeather.ArgsGetWeather{
		Country:  country,
		CityCode: city,
		DayCount: dataCount,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("mqtt push weather data, get weather data, ", err))
		err = MQTTClient.PublishWait(topic, 0, false, nil)
		return
	}
	var dataByte []byte
	dataByte, err = json.Marshal(data)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		err = MQTTClient.PublishWait(topic, 0, false, nil)
		return
	}
	//记录日志
	//CoreLog.Info("mqtt push weather data, device group: ", groupMark, ", code: ", deviceCode, ", country: ", country, ", city: ", city)
	//推送数据
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}
