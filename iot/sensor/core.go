package IOTSensor

import (
	CoreHighf "gitee.com/weeekj/weeekj_core/v5/core/highf"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTMQTT "gitee.com/weeekj/weeekj_core/v5/iot/mqtt"
	"github.com/robfig/cron"
)

/**
传感器模块
用于采集设备的传感信息，并提供数据展示接口。
为加快速写效率，数据将自动拆表。
*/

var (
	//定时器
	runTimer       *cron.Cron
	runHistoryLock = false
	//高频拦截器
	blocker CoreHighf.HighFBlocker
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	//初始化mqtt订阅
	if OpenSub {
		IOTMQTT.AppendSubFunc(initSub)
	}
	//初始化拦截器
	blocker.Init(5, false)
}

// 订阅方法集合
func initSub() {
	if token := IOTMQTT.MQTTClient.Subscribe("device/sensor/one", 0, subCreate); token.Wait() && token.Error() != nil {
		CoreLog.Error("iot sensor sub mqtt but token is error, " + token.Error().Error())
		return
	}
}
