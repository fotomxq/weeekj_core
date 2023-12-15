package IOTMQTTClient

import (
	CoreMQTTSimple "github.com/fotomxq/weeekj_core/v5/core/mqtt"
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
