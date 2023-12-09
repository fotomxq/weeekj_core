package IOTMQTT

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	"strings"
	"sync"
	"time"
)

// 推送MQTT列队处理机制
// 该模块会持续重复的发送MQTT数据，确保指定设备收到数据包
// 数据不会持久化，重启服务后将重置
// 每隔10秒、30秒、1分钟、5分钟、10分钟自动重发，最后每次10分钟重发，重试400次结束
var (
	//数据列队
	waitSendList []waitSendType
	waitSendLock sync.Mutex
)

// 数据列队结构
type waitSendType struct {
	//数据标识码
	// 重复数据会被替换为最新的数据包
	Mark string
	//等待发送时间
	// 避免高峰推送数据，造成数据堆叠异常
	WaitAt time.Time
	//重试次数
	TryCount int
	//topic
	Topic string
	//qos
	Qos byte
	//数据内容
	Data []byte
}

// PushWait 推送新的MQTT到列队
func PushWait(mark string, topic string, qos byte, data []byte) {
	waitData := waitSendType{
		Mark:     mark,
		WaitAt:   time.Time{},
		TryCount: 0,
		Topic:    topic,
		Qos:      qos,
		Data:     data,
	}
	pushMqttWait(&waitData)
	waitSendLock.Lock()
	defer waitSendLock.Unlock()
	for k, v := range waitSendList {
		if v.Mark == mark {
			waitSendList[k].Data = data
			waitSendList[k].TryCount = 0
			return
		}
	}
	waitSendList = append(waitSendList, waitData)
}

// DeletePushWait 删除等待数据
func DeletePushWait(mark string) {
	waitSendLock.Lock()
	defer waitSendLock.Unlock()
	var newList []waitSendType
	for _, v := range waitSendList {
		if v.Mark == mark {
			continue
		}
		newList = append(newList, v)
	}
	waitSendList = newList
}

func DeleteSearchPushWait(mark string) {
	waitSendLock.Lock()
	defer waitSendLock.Unlock()
	var newList []waitSendType
	for _, v := range waitSendList {
		if strings.Contains(v.Mark, mark) {
			continue
		}
		newList = append(newList, v)
	}
	waitSendList = newList
}

// 定时推送模块
func runWaitSend() {
	waitSendLock.Lock()
	defer waitSendLock.Unlock()
	nowAt := CoreFilter.GetNowTimeCarbon()
	var waitDeleteKey []int
	for k, v := range waitSendList {
		//如果没到数据，则跳过
		if v.WaitAt.Unix() > nowAt.Time.Unix() {
			continue
		}
		var nextAt time.Time
		//计算下一次执行时间
		switch v.TryCount {
		case 0:
			//延迟推送
			nextAt = nowAt.AddSeconds(10).Time
		case 1:
			//延迟推送
			nextAt = nowAt.AddSeconds(30).Time
		case 2:
			//延迟推送
			nextAt = nowAt.AddMinutes(5).Time
		case 9:
			//延迟推送
			nextAt = nowAt.AddMinutes(10).Time
		default:
			//重试30次，延迟推送
			nextAt = nowAt.AddMinutes(30).Time
			if v.TryCount >= 30 {
				waitDeleteKey = append(waitDeleteKey, k)
				continue
			}
		}
		//推送数据
		pushMqttWait(&v)
		//延迟等待
		waitSendList[k].WaitAt = nextAt
		waitSendList[k].TryCount += 1
	}
	var newList []waitSendType
	for _, v := range waitDeleteKey {
		for k2, v2 := range waitSendList {
			if v == k2 {
				continue
			}
			newList = append(newList, v2)
		}
	}
	waitSendList = newList
}

// pushMqttWait 推送wait MQTT
func pushMqttWait(data *waitSendType) {
	_ = MQTTClient.Publish(data.Topic, data.Qos, false, data.Data)
}
