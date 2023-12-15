package IOTMQTT

import (
	"errors"
	"fmt"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	IOTMission "github.com/fotomxq/weeekj_core/v5/iot/mission"
)

// 推送目标消息
// 内部函数
func sendAutoDevice(val string, infoData IOTDevice.FieldsAutoInfo) (err error) {
	//如果存在任务action，则走任务模式，否则直接将数据包推送给目标设备
	if infoData.SendAction != "" {
		_, _, err = IOTMission.CreateMission(&IOTMission.ArgsCreateMission{
			OrgID:      infoData.OrgID,
			DeviceID:   infoData.ReportDeviceID,
			ParamsData: infoData.ParamsData,
			Action:     infoData.SendAction,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("create device mission, ", err))
			return
		}
		return
	} else {
		pushAuto(infoData.ReportDeviceID, val, infoData)
		return
	}
}
