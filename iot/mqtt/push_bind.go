package IOTMQTT

import (
	"encoding/json"
	"errors"
	"fmt"
	IOTBind "github.com/fotomxq/weeekj_core/v5/iot/bind"
)

// PushDeviceBind 推送设备所有绑定关系
func PushDeviceBind(deviceID int64) (err error) {
	//获取数据
	type dataType struct {
		//数据结构
		DataList []IOTBind.FieldsBind `json:"dataList"`
	}
	var data dataType
	var rawData IOTBind.FieldsBind
	rawData, err = IOTBind.GetBindByDeviceID(deviceID)
	if err != nil {
		err = errors.New("no data, " + err.Error())
		return
	}
	data.DataList = append(data.DataList, rawData)
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(data)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("device/bind/all/", deviceID)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}
