package IOTMQTTClient

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
)

// ArgsPushDeviceOnline 设备在线情况更正参数
type ArgsPushDeviceOnline struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//目标设备ID
	DeviceID int64 `json:"deviceID"`
	//是否在线
	IsOnline bool `json:"isOnline"`
}

// PushDeviceOnline 设备在线情况更正
// device/online
func PushDeviceOnline(args ArgsPushDeviceOnline) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := "device/online"
	err = mqttClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		err = errors.New(fmt.Sprint("mqtt push base online, ", err))
		return
	}
	return
}

// ArgsPushDeviceInfosUpdate 全量更新设备数据
type ArgsPushDeviceInfosUpdate struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//目标设备ID
	DeviceID int64 `json:"deviceID"`
	//数据集合
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

// PushDeviceInfosUpdate 全量更新设备数据
// device/infos/update
func PushDeviceInfosUpdate(args ArgsPushDeviceInfosUpdate) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := "device/infos/update"
	err = mqttClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		err = errors.New(fmt.Sprint("mqtt push base online, ", err))
		return
	}
	return
}

// ArgsPushDeviceInfoUpdate 分量更新设备数据参数
type ArgsPushDeviceInfoUpdate struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//目标设备ID
	DeviceID int64 `json:"deviceID"`
	//数据集合
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

// PushDeviceInfoUpdate 分量更新设备数据
// device/info/update
func PushDeviceInfoUpdate(args ArgsPushDeviceInfoUpdate) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := "device/info/update"
	err = mqttClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		err = errors.New(fmt.Sprint("mqtt push base online, ", err))
		return
	}
	return
}

// ArgsPushDeviceFindData 分量更新设备数据参数
type ArgsPushDeviceFindData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

// PushDeviceFindData 分量更新设备数据
// device/find
func PushDeviceFindData(args ArgsPushDeviceFindData) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := "device/find"
	err = mqttClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		err = errors.New(fmt.Sprint("mqtt push base online, ", err))
		return
	}
	return
}

// ArgsPushDeviceGroupFindData 分量更新设备数据参数
type ArgsPushDeviceGroupFindData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

// PushDeviceGroupFindData 分量更新设备数据
// group/find
func PushDeviceGroupFindData(args ArgsPushDeviceGroupFindData) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := "group/find"
	err = mqttClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		err = errors.New(fmt.Sprint("mqtt push base online, ", err))
		return
	}
	return
}
