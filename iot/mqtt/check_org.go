package IOTMQTT

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	"reflect"
)

//CheckDeviceAndOrg 组织专用检查程序
/**
1. 解析数据
2. 检查设备key是否正确
3. 检查设备授权给组织情况
*/
func CheckDeviceAndOrg(message mqtt.Message, resultData interface{}) (b bool) {
	//如果全局启动了mqtt_debug模式，则忽略权限检查
	mqttDebug, err := BaseConfig.GetDataBool("MQTTDebug")
	if err != nil {
		mqttDebug = false
	}
	if mqttDebug {
		return true
	}
	//解析参数数据
	if err := json.Unmarshal(message.Payload(), &resultData); err != nil {
		CoreLog.MqttError("mqtt get params, json, ", err)
		return
	}
	var deviceKeys IOTDevice.ArgsCheckDeviceKey
	var orgID int64
	//核对key
	for key := 0; key < reflect.TypeOf(resultData).Elem().NumField(); key += 1 {
		field := reflect.TypeOf(resultData).Elem().Field(key)
		val := reflect.ValueOf(resultData).Elem().Field(key)
		fieldMark := field.Tag.Get("json")
		switch fieldMark {
		case "keys":
			deviceKeys = val.Interface().(IOTDevice.ArgsCheckDeviceKey)
			break
		case "orgID":
			orgID = val.Interface().(int64)
			break
		}
	}
	//通过数据集，获取设备信息
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&deviceKeys)
	if err != nil {
		CoreLog.MqttError("mqtt get params, key, ", err)
		return
	}
	//检查设备授权
	_, err = IOTDevice.CheckOperate(&IOTDevice.ArgsCheckOperate{
		DeviceID: deviceID,
		OrgID:    orgID,
	})
	if err != nil {
		CoreLog.MqttError("mqtt check operate, topic: ", message.Topic(), ", check operate, deviceID: ", deviceID, ", orgID: ", orgID, ", err: ", err)
		return
	}
	//反馈成功
	return true
}
