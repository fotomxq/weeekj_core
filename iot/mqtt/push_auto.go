package IOTMQTT

import (
	"encoding/json"
	"errors"
	"fmt"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
)

// 推送设备信息变更通知
func pushAuto(deviceID int64, val string, infoData IOTDevice.FieldsAutoInfo) {
	//获取数据
	type dataType struct {
		//数据包
		InfoData IOTDevice.FieldsAutoInfo `json:"infoData"`
		//触发的值
		Val string `json:"val"`
	}
	var data dataType
	var err error
	data.InfoData = infoData
	data.Val = val
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(data)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("device/auto/report/", deviceID)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}
