package IOTMQTT

import (
	"encoding/json"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 请求获取当前设备控制列表
type subOperateData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

func subOperate(client mqtt.Client, message mqtt.Message) {
	var resultData subOperateData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub get operate list, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub get operate list, key, ", err)
		return
	}
	//获取设备的信息
	deviceData, err := IOTDevice.GetDeviceByCode(&IOTDevice.ArgsGetDeviceByCode{
		GroupMark: resultData.Keys.GroupMark,
		Code:      resultData.Keys.Code,
	})
	if err != nil {
		CoreLog.MqttError("mqtt sub get operate list, device not exist, ", err)
		return
	}
	//推送信息
	if err := PushOperate(resultData.Keys.GroupMark, resultData.Keys.Code, deviceData.ID); err != nil {
		CoreLog.MqttError("mqtt sub get operate list, push operate list, ", err)
		return
	}
}
