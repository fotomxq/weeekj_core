package IOTMQTT

import (
	"encoding/json"
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
)

func runUpdateData() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device update data run error, ", r)
		}
	}()
	waitPushUpdateDateListLock.Lock()
	defer waitPushUpdateDateListLock.Unlock()
	type dataType struct {
		//动作类型
		// 具体更新的数据来源
		Action string `json:"action"`
		//更新ID
		ID int64 `json:"id"`
	}
	for _, v := range waitPushUpdateDateList {
		//获取数据
		var data dataType
		data.Action = v.Action
		data.ID = v.UpdateID
		//打包数据集合
		var dataByte []byte
		dataByte, err := json.Marshal(data)
		if err != nil {
			CoreLog.MqttError("iot device update data run error, json error, ", err)
			continue
		}
		//推送数据
		topic := fmt.Sprint("v2/device/update/data/", v.OrgID)
		err = MQTTClient.PublishWait(topic, 0, false, dataByte)
		if err != nil {
			CoreLog.MqttError("iot device update data run error, push data, ", err)
			continue
		}
	}
	waitPushUpdateDateList = []waitPushUpdateDateListType{}
}
