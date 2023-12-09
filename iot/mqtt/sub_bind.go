package IOTMQTT

import (
	"encoding/json"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// subDeviceBind 请求设备绑定关系
type subDeviceBindData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

func subDeviceBind(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceBindData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device bind, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.MqttError("mqtt sub device bind, key, ", err)
		return
	}
	if err = PushDeviceBind(deviceID); err != nil {
		CoreLog.MqttError("mqtt sub device bind, push bind all data, ", err)
		return
	}
}
