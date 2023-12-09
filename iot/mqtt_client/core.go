package IOTMQTTClient

import (
	CoreMQTTSimple "gitee.com/weeekj/weeekj_core/v5/core/mqtt"
	"github.com/robfig/cron"
)

var (
	//定时器
	runTimer       *cron.Cron
	runConnectLock = false
	//订阅前缀
	mqttPrefix = "device"
	//mqtt对象
	mqttClient CoreMQTTSimple.MQTTSimple
)
