package IOTMQTT

import (
	"errors"
)

// 订阅主题
func initSub() (err error) {
	//设备在线情况更正
	if token := MQTTClient.Subscribe("device/online", 0, subDeviceOnline); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//全量更新设备信息
	if token := MQTTClient.Subscribe("device/infos/update", 0, subDeviceInfosUpdate); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//分量更新设备信息
	if token := MQTTClient.Subscribe("device/info/update", 0, subDeviceInfoUpdate); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//查询设备信息
	if token := MQTTClient.Subscribe("device/info/find", 0, subDeviceFind); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//查询设备所在组
	if token := MQTTClient.Subscribe("group/find", 0, subDeviceGroupFind); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//设备掉线遗嘱
	if token := MQTTClient.Subscribe("device/lost", 0, subDeviceLost); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//错误推送
	if token := MQTTClient.Subscribe("device/error", 0, subDeviceError); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//反馈设备任务结果
	if token := MQTTClient.Subscribe("device/mission/result", 0, subDeviceMissionResult); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//设备追踪
	if token := MQTTClient.Subscribe("device/track", 0, subTrack); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//天气预报
	if token := MQTTClient.Subscribe("weather/city", 0, subWeather); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//应用更新
	if token := MQTTClient.Subscribe("app/update", 0, subAppUpdate); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//设备控制关系
	if token := MQTTClient.Subscribe("device/operate", 0, subOperate); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//设备绑定关系
	if token := MQTTClient.Subscribe("device/bind", 0, subDeviceBind); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//v2接口
	//查询设备信息
	if token := MQTTClient.Subscribe("v2/device/info/find", 0, subV2DeviceFind); token.Wait() && token.Error() != nil {
		err = errors.New("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	return
}
