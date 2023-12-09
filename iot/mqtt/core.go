package IOTMQTT

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreMQTTSimple "gitee.com/weeekj/weeekj_core/v5/core/mqtt"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	"github.com/robfig/cron"
)

//MQTT服务端程序
// 该模块发送和接收设备信息，并完成对应的任务处理核心

var (
	//定时器
	runTimer          *cron.Cron
	runConnectLock    = false
	runMissionLock    = false
	runUpdateDataLock = false
	runWaitSendLock   = false
	//订阅前缀
	mqttPrefix = "device"
	//MQTTClient mqtt对象
	MQTTClient CoreMQTTSimple.MQTTSimple
	//MQTTIsConnect mqtt初始化和连接标记
	MQTTIsConnect = false
	//等待激活的订阅函数组
	subFunc []func()
	//重试次数
	connectTryCount = 0
	//OpenBaseMission 是否启动任务广播环节
	OpenBaseMission = true
)

// AppendSubFunc 写入等待订阅函数包
func AppendSubFunc(sF func()) {
	subFunc = append(subFunc, sF)
}

// SubBefore 检查商户的所有权等一揽子前缀处理
func SubBefore(keys IOTDevice.ArgsCheckDeviceKey, orgID int64, logAppend interface{}) (deviceID int64, operateData IOTDevice.FieldsOperate, b bool) {
	var err error
	if deviceID, err = IOTDevice.CheckDeviceKeyAndDeviceID(&keys); err != nil {
		CoreLog.Error(logAppend, "key, ", err)
		return
	}
	if operateData, err = IOTDevice.CheckOperate(&IOTDevice.ArgsCheckOperate{
		DeviceID: deviceID,
		OrgID:    orgID,
	}); err != nil {
		CoreLog.Error(logAppend, "check operate, ", err)
		return
	}
	b = true
	return
}
