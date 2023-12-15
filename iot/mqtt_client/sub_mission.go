package IOTMQTTClient

import (
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTMQTT "github.com/fotomxq/weeekj_core/v5/iot/mqtt"
)

// SubMissionSendGroupHandle 设备组消息
// device/mission/send/group/[设备组标识码]
type SubMissionSendGroupHandle func(data IOTMQTT.ArgsPushMissionToDevice)

// SubMissionSendGroup 设备组消息
func SubMissionSendGroup(groupMark string, handle SubMissionSendGroupHandle) (token mqtt.Token, err error) {
	topic := fmt.Sprint("device/mission/send/group/", groupMark)
	token = mqttClient.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		var resultData IOTMQTT.ArgsPushMissionToDevice
		resultByte := message.Payload()
		if err := json.Unmarshal(resultByte, &resultData); err != nil {
			CoreLog.MqttError("sub mqtt data json is error, ", err.Error())
			return
		}
		//反馈数据结果
		handle(resultData)
	})
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub mission send group, ", token.Error()))
		return
	}
	return
}

// SubMissionSendGroupCancel 取消订阅设备组消息
func SubMissionSendGroupCancel(groupMark string) (token mqtt.Token, err error) {
	topic := fmt.Sprint("device/mission/send/group/", groupMark)
	token = mqttClient.SubscribeCancel(topic)
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub mission send group cancel, ", token.Error()))
		return
	}
	return
}

// SubMissionSendDeviceHandle 设备消息
// device/mission/send/group/[设备组标识码]/code/[设备厂商编码]
type SubMissionSendDeviceHandle func(data IOTMQTT.ArgsPushMissionToGroup)

// SubMissionSendDevice 设备任务消息
func SubMissionSendDevice(groupMark, deviceCode string, handle SubMissionSendDeviceHandle) (token mqtt.Token, err error) {
	topic := fmt.Sprint("device/mission/send/group/", groupMark, "/code/", deviceCode)
	token = mqttClient.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		var resultData IOTMQTT.ArgsPushMissionToGroup
		resultByte := message.Payload()
		if err := json.Unmarshal(resultByte, &resultData); err != nil {
			CoreLog.MqttError("sub mqtt data json is error, ", err.Error())
			return
		}
		//反馈数据结果
		handle(resultData)
	})
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("sub mqtt but token is error, ", token.Error()))
		return
	}
	return
}

// SubMissionSendDeviceCancel 取消设备任务消息
func SubMissionSendDeviceCancel(groupMark, deviceCode string) (token mqtt.Token, err error) {
	topic := fmt.Sprint("device/mission/send/group/", groupMark, "/code/", deviceCode)
	token = mqttClient.SubscribeCancel(topic)
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub mission send device, ", token.Error()))
		return
	}
	return
}
