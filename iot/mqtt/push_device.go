package IOTMQTT

import (
	"encoding/json"
	"errors"
	"fmt"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
)

// PushDeviceUpdateOnline 要求设备更新设备在线状态
func PushDeviceUpdateOnline(groupMark, deviceCode string) (err error) {
	//推送数据
	topic := fmt.Sprint("device/online/group/", groupMark, "/code/", deviceCode)
	err = MQTTClient.PublishWait(topic, 0, false, nil)
	return
}

// ArgsPushDeviceNeedInfo 主动更新设备信息参数
type ArgsPushDeviceNeedInfo struct {
	DeviceData IOTDevice.FieldsDevice `json:"deviceData"`
}

// PushDeviceNeedInfo 主动更新设备信息
func PushDeviceNeedInfo(groupMark, deviceCode string, args ArgsPushDeviceNeedInfo) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("device/info/need/group/", groupMark, "/code/", deviceCode)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}

// PushDeviceNeedInfoByID 使用设备ID推送设备更新
func PushDeviceNeedInfoByID(deviceID int64) error {
	//获取设备信息
	deviceData, err := IOTDevice.GetDeviceByID(&IOTDevice.ArgsGetDeviceByID{
		ID:    deviceID,
		OrgID: -1,
	})
	if err != nil {
		return err
	}
	//获取设备组
	groupData, err := IOTDevice.GetGroupByID(&IOTDevice.ArgsGetGroupByID{
		ID: deviceData.GroupID,
	})
	if err != nil {
		return err
	}
	//推送请求
	err = PushDeviceNeedInfo(groupData.Mark, deviceData.Code, ArgsPushDeviceNeedInfo{
		DeviceData: deviceData,
	})
	return err
}

// ArgsPushDeviceNeedGroup 主动更新设备组信息参数
type ArgsPushDeviceNeedGroup struct {
	GroupData IOTDevice.FieldsGroup `json:"groupData"`
}

// PushDeviceNeedGroup 主动更新设备组信息
func PushDeviceNeedGroup(groupMark string, args ArgsPushDeviceNeedGroup) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("group/info/", groupMark)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}
