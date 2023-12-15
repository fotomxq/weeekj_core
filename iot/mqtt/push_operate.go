package IOTMQTT

import (
	"encoding/json"
	"errors"
	"fmt"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
)

func PushOperate(groupMark, deviceCode string, deviceID int64) (err error) {
	var operateList []IOTDevice.FieldsOperate
	operateList, err = IOTDevice.GetOperateByDeviceID(&IOTDevice.ArgsGetOperateByDeviceID{
		DeviceID: deviceID,
	})
	if err != nil {
		return
	}
	var dataByte []byte
	dataByte, err = json.Marshal(operateList)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("device/operate/group/", groupMark, "/code/", deviceCode)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}
