package IOTMQTT

import (
	"encoding/json"
	"errors"
	"fmt"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	"time"
)

// ArgsPushMissionToDevice 推送新的任务参数
type ArgsPushMissionToDevice struct {
	//设备ID
	ID int64 `json:"id"`
	//任务ID
	MissionID int64 `json:"missionID"`
	//过期时间
	ExpireAt time.Time `json:"expireAt"`
	//任务动作
	Action string `json:"action"`
	//发送请求数据集合
	ParamsData []byte `json:"paramsData"`
}

// PushMissionToDevice 推送新的任务
// device/mission
func PushMissionToDevice(args *ArgsPushMissionToDevice) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	var deviceData IOTDevice.GetDeviceGroupData
	deviceData, err = IOTDevice.GetDeviceGroup(&IOTDevice.ArgsGetDeviceGroup{
		DeviceID: args.ID,
	})
	if err != nil {
		err = errors.New("get device data, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("device/mission/send/group/", deviceData.GroupMark, "/code/", deviceData.Code)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}

// ArgsPushMissionToGroup 推送新的任务到分组参数
type ArgsPushMissionToGroup struct {
	//设备ID
	ID int64 `json:"id"`
	//任务ID
	MissionID int64 `json:"missionID"`
	//分组ID
	GroupID int64 `json:"groupID"`
	//过期时间
	ExpireAt time.Time `json:"expireAt"`
	//任务动作
	Action string `json:"action"`
	//发送请求数据集合
	ParamsData []byte `json:"paramsData"`
}

// PushMissionToGroup 推送新的任务到分组
func PushMissionToGroup(args ArgsPushMissionToGroup) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	var deviceData IOTDevice.GetDeviceGroupData
	deviceData, err = IOTDevice.GetDeviceGroup(&IOTDevice.ArgsGetDeviceGroup{
		DeviceID: args.ID,
	})
	if err != nil {
		err = errors.New("get device data, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("device/mission/send/group/", deviceData.GroupMark)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}
