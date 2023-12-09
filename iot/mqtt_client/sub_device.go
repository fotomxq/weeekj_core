package IOTMQTTClient

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTMQTT "gitee.com/weeekj/weeekj_core/v5/iot/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SubBaseOnlineHandle 设备在线判定
// device/online/group/[设备组标识码]/code/[设备厂商编码]
type SubBaseOnlineHandle func()

// SubBaseOnline 服务端要求设备更新设备在线状态
func SubBaseOnline(groupMark, deviceCode string, handle SubBaseOnlineHandle) (token mqtt.Token, err error) {
	topic := fmt.Sprint("device/online/group/", groupMark, "/code/", deviceCode)
	token = mqttClient.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		//反馈数据结果
		handle()
	})
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub base online, ", token.Error()))
		return
	}
	return
}

// SubBaseOnlineCancel 取消服务端要求设备更新设备在线状态
func SubBaseOnlineCancel(groupMark, deviceCode string) (token mqtt.Token, err error) {
	topic := fmt.Sprint("device/online/group/", groupMark, "/code/", deviceCode)
	token = mqttClient.SubscribeCancel(topic)
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub base online cancel, ", token.Error()))
		return
	}
	return
}

// SubBaseInfoNeedHandle 下发更新设备信息
// device/info/need/group/[设备组标识码]/code/[设备厂商编码]
type SubBaseInfoNeedHandle func(data IOTMQTT.ArgsPushDeviceNeedInfo)

// SubBaseInfoUpdate 服务端下发最新的设备数据包
func SubBaseInfoUpdate(groupMark, deviceCode string, handle SubBaseInfoNeedHandle) (token mqtt.Token, err error) {
	topic := fmt.Sprint("device/info/need/group/", groupMark, "/code/", deviceCode)
	token = mqttClient.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		var resultData IOTMQTT.ArgsPushDeviceNeedInfo
		resultByte := message.Payload()
		if err := json.Unmarshal(resultByte, &resultData); err != nil {
			CoreLog.MqttError("sub mqtt data json is error, ", err.Error())
			return
		}
		//反馈数据结果
		handle(resultData)
	})
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub base info update, ", token.Error()))
		return
	}
	return
}

// SubBaseInfoUpdateCancel 取消订阅服务端下发最新的设备数据包
func SubBaseInfoUpdateCancel(groupMark, deviceCode string) (token mqtt.Token, err error) {
	topic := fmt.Sprint("device/info/need/group/", groupMark, "/code/", deviceCode)
	token = mqttClient.SubscribeCancel(topic)
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub base info update cancel, ", token.Error()))
		return
	}
	return
}

// SubBaseGroupInfoHandle 下发设备组信息
// group/info[设备组标识码]
type SubBaseGroupInfoHandle func(data IOTMQTT.ArgsPushDeviceNeedGroup)

// SubBaseGroupInfoUpdate 服务端下发设备组信息
func SubBaseGroupInfoUpdate(groupMark string, handle SubBaseGroupInfoHandle) (token mqtt.Token, err error) {
	topic := fmt.Sprint("group/info/", groupMark)
	token = mqttClient.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		var resultData IOTMQTT.ArgsPushDeviceNeedGroup
		resultByte := message.Payload()
		if err := json.Unmarshal(resultByte, &resultData); err != nil {
			CoreLog.MqttError("sub mqtt data json is error, ", err.Error())
			return
		}
		//反馈数据结果
		handle(resultData)
	})
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub base group info update, ", token.Error()))
		return
	}
	return
}

// SubBaseGroupInfoUpdateCancel 取消订阅服务端下发设备组信息
func SubBaseGroupInfoUpdateCancel(groupMark string) (token mqtt.Token, err error) {
	topic := fmt.Sprint("group/info/", groupMark)
	token = mqttClient.SubscribeCancel(topic)
	if token.Wait() && token.Error() != nil {
		err = errors.New(fmt.Sprint("mqtt sub base group info update cancel, ", token.Error()))
		return
	}
	return
}
