package IOTMQTT

import (
	"encoding/json"
	"fmt"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
)

// PushDeviceRest 请求重启设备
func PushDeviceRest(deviceID int64) {
	pushIOTAction(deviceID, "rest")
}

// 推送通用的action方法
func pushIOTAction(deviceID int64, action string) {
	appendLog := "iot core push mqtt iot action, "
	//打包数据集合
	dataByte, err := json.Marshal(map[string]interface{}{
		"action": action,
	})
	if err != nil {
		CoreLog.Warn(appendLog, "json err: ", err)
		return
	}
	//推送数据
	topic := fmt.Sprint("v2/iot/action/", deviceID)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}
