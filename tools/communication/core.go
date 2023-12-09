package ToolsCommunication

import (
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	ClassTag "gitee.com/weeekj/weeekj_core/v5/class/tag"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTMQTT "gitee.com/weeekj/weeekj_core/v5/iot/mqtt"
	"github.com/robfig/cron"
)

//通讯处理模块
// 该模块主要提供通讯衔接和处理优化功能

var (
	//Sort 分类
	Sort = ClassSort.Sort{
		SortTableName: "tools_communication_sort",
	}
	//Tag 标签
	Tag = ClassTag.Tag{
		TagTableName: "tools_communication_tag",
	}
	//定时器
	runTimer *cron.Cron
	runLock  = false
	//OpenSub 是否启动订阅
	OpenSub = false
)

// Init 初始化
func Init() {
	if OpenSub {
		//初始化mqtt订阅
		IOTMQTT.AppendSubFunc(initMQTT)
	}
}

// 订阅方法集合
func initMQTT() {
	//请求发起新的双向聊天
	if token := IOTMQTT.MQTTClient.Subscribe("tools/communication/room/new", 0, subNewRoom); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//请求聊天室信息
	if token := IOTMQTT.MQTTClient.Subscribe("tools/communication/room/info", 0, subRoomInfo); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//请求设备绑定的所有房间
	if token := IOTMQTT.MQTTClient.Subscribe("tools/communication/room/from", 0, subRoomInfo); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//续租聊天
	if token := IOTMQTT.MQTTClient.Subscribe("tools/communication/expire", 0, subExpire); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
	//退出聊天室
	if token := IOTMQTT.MQTTClient.Subscribe("tools/communication/out", 0, subOutRoom); token.Wait() && token.Error() != nil {
		CoreLog.Error("sub mqtt but token is error, " + token.Error().Error())
		return
	}
}
