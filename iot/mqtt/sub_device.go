package IOTMQTT

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	IOTLog "gitee.com/weeekj/weeekj_core/v5/iot/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// subDeviceOnline 设备在线情况更正
type subDeviceOnlineData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//目标设备ID
	DeviceID int64 `json:"deviceID"`
	//是否在线
	IsOnline bool `json:"isOnline"`
}

func subDeviceOnline(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceOnlineData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device online, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub device online, key, ", err)
		return
	}
	if err := IOTDevice.UpdateDeviceOnline(&IOTDevice.ArgsUpdateDeviceOnline{
		ID:       resultData.DeviceID,
		IsOnline: resultData.IsOnline,
	}); err != nil {
		CoreLog.MqttError("mqtt sub device online, update device online, ", err)
		return
	}
	//必须推送请求，告知设备在线成功，设备收到后可更新时间等信息，作为心跳包处理机制之一
	if resultData.IsOnline {
		if err := PushDeviceUpdateOnline(resultData.Keys.GroupMark, resultData.Keys.Code); err != nil {
			CoreLog.MqttError("mqtt sub device online, send device online, ", err)
			return
		}
	}
}

// subDeviceInfosUpdate 全量更新设备数据
type subDeviceInfosUpdateData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//数据集合
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

func subDeviceInfosUpdate(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceInfosUpdateData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device infos update, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.MqttError("mqtt sub device infos update, key, ", err)
		return
	}
	if err := IOTDevice.UpdateDeviceInfos(&IOTDevice.ArgsUpdateDeviceInfos{
		ID:     deviceID,
		Params: resultData.Params,
	}); err != nil {
		CoreLog.MqttError("mqtt sub device infos update, update device infos, ", err, ", device id: ", deviceID)
		return
	} else {
		deviceData, err := IOTDevice.GetOperateAndDevice(&IOTDevice.ArgsGetOperateAndDevice{
			DeviceID: deviceID,
		})
		if err != nil {
			CoreLog.MqttError("mqtt sub device infos update, device or operate not exist, ", err)
			return
		}
		for _, v := range deviceData {
			IOTLog.Append(&IOTLog.ArgsAppend{
				OrgID:    v.OrgID,
				GroupID:  v.GroupID,
				DeviceID: v.ID,
				Mark:     "update_infos",
				Content:  "全量更新设备信息",
			})
		}
		//推送auto设计
		for _, v := range resultData.Params {
			reportDataList, b, err := IOTDevice.OpenAutoInfo(&IOTDevice.ArgsOpenAutoInfo{
				OrgID:    -1,
				DeviceID: deviceID,
				Mark:     v.Mark,
				Val:      v.Val,
			})
			if err != nil {
				CoreLog.MqttError("mqtt sub device infos update, open auto info, ", err)
				continue
			}
			if b {
				for _, v2 := range reportDataList {
					err = sendAutoDevice(v.Val, v2)
					if err != nil {
						CoreLog.MqttError("mqtt sub device infos update, open auto info, send auto device, ", err)
					}
				}
			}
		}
	}
}

// subDeviceInfoUpdate 分量更新设备数据
type subDeviceInfoUpdateData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//数据集合
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

func subDeviceInfoUpdate(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceInfoUpdateData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device info update, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.MqttError("mqtt sub device info update, key, ", err)
		return
	}
	if err := IOTDevice.UpdateDeviceInfo(&IOTDevice.ArgsUpdateDeviceInfo{
		ID:     deviceID,
		Params: resultData.Params,
	}); err != nil {
		CoreLog.MqttError("mqtt sub device info update, update device info, ", err, ", device id: ", deviceID)
		return
	} else {
		deviceData, err := IOTDevice.GetOperateAndDevice(&IOTDevice.ArgsGetOperateAndDevice{
			DeviceID: deviceID,
		})
		if err != nil {
			CoreLog.MqttError("mqtt sub device info update, device or operate not exist, ", err)
			return
		}
		for _, v := range deviceData {
			IOTLog.Append(&IOTLog.ArgsAppend{
				OrgID:    v.OrgID,
				GroupID:  v.GroupID,
				DeviceID: v.ID,
				Mark:     "update_info",
				Content:  "分量更新设备信息",
			})
		}
		//推送auto设计
		for _, v := range resultData.Params {
			reportDataList, b, err := IOTDevice.OpenAutoInfo(&IOTDevice.ArgsOpenAutoInfo{
				OrgID:    -1,
				DeviceID: deviceID,
				Mark:     v.Mark,
				Val:      v.Val,
			})
			if err != nil {
				CoreLog.MqttError("mqtt sub device info update, open auto info, ", err)
				continue
			}
			if b {
				for _, v2 := range reportDataList {
					err = sendAutoDevice(v.Val, v2)
					if err != nil {
						CoreLog.MqttError("mqtt sub device info update, open auto info, send auto device, ", err)
					}
				}
			}
		}
	}
}

// subDeviceFind 查询设备
type subDeviceFindData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

type subDeviceFindReport struct {
	DeviceData IOTDevice.FieldsDevice `json:"deviceData"`
}

func subDeviceFind(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceFindData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device find, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub device find, key, ", err)
		return
	}
	//获取设备信息
	deviceData, err := IOTDevice.GetDeviceByCode(&IOTDevice.ArgsGetDeviceByCode{
		GroupMark: resultData.Keys.GroupMark,
		Code:      resultData.Keys.Code,
	})
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device find, not exist, ", err)
		return
	}
	//推送消息
	postData := subDeviceFindReport{
		DeviceData: deviceData,
	}
	var dataByte []byte
	dataByte, err = json.Marshal(postData)
	if err != nil {
		CoreLog.MqttError("mqtt sub device find, json, ", err)
		return
	}
	topic := fmt.Sprint("device/info/need/group/", resultData.Keys.GroupMark, "/code/", deviceData.Code)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		err = errors.New(fmt.Sprint("mqtt push device info, ", err))
		return
	}
}

// subDeviceGroupFind 寻找设备的分组数据
type subDeviceGroupFindData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//设备组标识码
	GroupMark string `json:"groupMark" check:"mark"`
}

type subDeviceGroupFindReport struct {
	GroupData IOTDevice.FieldsGroup `json:"groupData"`
}

func subDeviceGroupFind(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceGroupFindData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device group find, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub device group find, key, ", err)
		return
	}
	//获取设备信息
	groupData, err := IOTDevice.GetGroupByMark(&IOTDevice.ArgsGetGroupByMark{
		Mark: resultData.GroupMark,
	})
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device group find, not exist, ", err)
		return
	}
	//推送消息
	postData := subDeviceGroupFindReport{
		GroupData: groupData,
	}
	var dataByte []byte
	dataByte, err = json.Marshal(postData)
	if err != nil {
		CoreLog.MqttError("mqtt sub device group find, json, ", err)
		return
	}
	topic := fmt.Sprint("group/info/", groupData.Mark)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		err = errors.New(fmt.Sprint("mqtt push device group, ", err))
		return
	}
}

// subDeviceLostData 寻找设备的分组数据
type subDeviceLostData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

// subDeviceLost 掉线遗嘱
func subDeviceLost(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceGroupFindData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device lost, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.MqttError("mqtt sub device lost, key, ", err)
		return
	}
	//标记为掉线
	if err := IOTDevice.UpdateDeviceOnline(&IOTDevice.ArgsUpdateDeviceOnline{
		ID:       deviceID,
		IsOnline: false,
	}); err != nil {
		CoreLog.MqttError("mqtt sub device lost, update online, ", err)
		return
	}
}
