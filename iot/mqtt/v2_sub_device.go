package IOTMQTT

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	IOTBind "github.com/fotomxq/weeekj_core/v5/iot/bind"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
)

// subDeviceFind 查询设备
type subV2DeviceFindData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

type subV2DeviceFindReport struct {
	//设备ID
	DeviceID int64 `json:"deviceID"`
	//状态
	// 0 public 公共可用 / 1 private 私有 / 2 ban 停用
	Status int `db:"status" json:"status"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
	//设备组过期时间
	GroupExpireTime int64 `json:"groupExpireTime"`
	//扩展参数
	GroupParams CoreSQLConfig.FieldsConfigsType `json:"groupParams"`
	//授权商户ID
	BindOrgID int64 `json:"bindOrgID"`
	//绑定内容，只反馈第一条
	//附加模块
	BindFromInfo CoreSQLFrom.FieldsFrom `json:"bindFromInfo"`
	//扩展参数
	BindParams CoreSQLConfig.FieldsConfigsType `json:"bindParams"`
}

func subV2DeviceFind(client mqtt.Client, message mqtt.Message) {
	var resultData subV2DeviceFindData
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
	if err != nil {
		CoreLog.MqttError("mqtt sub device find, device not exist, ", err)
		return
	}
	//获取分组数据
	groupData, err := IOTDevice.GetGroupByID(&IOTDevice.ArgsGetGroupByID{
		ID: deviceData.GroupID,
	})
	if err != nil {
		CoreLog.MqttError("mqtt sub device find, group not exist, ", err)
		return
	}
	//获取授权
	operateList, err := IOTDevice.GetOperateByDeviceID(&IOTDevice.ArgsGetOperateByDeviceID{
		DeviceID: deviceData.ID,
	})
	var orgID int64
	if err == nil && len(operateList) > 0 {
		orgID = operateList[0].OrgID
	}
	//绑定关系
	bindData, _ := IOTBind.GetBindDevice(orgID, deviceData.ID)
	//推送消息
	postData := subV2DeviceFindReport{
		DeviceID:        deviceData.ID,
		Status:          deviceData.Status,
		GroupID:         deviceData.GroupID,
		Params:          deviceData.Params,
		GroupExpireTime: groupData.ExpireTime,
		GroupParams:     groupData.Params,
		BindOrgID:       orgID,
		BindFromInfo:    bindData.FromInfo,
		BindParams:      bindData.Params,
	}
	var dataByte []byte
	dataByte, err = json.Marshal(postData)
	if err != nil {
		CoreLog.MqttError("mqtt sub device find, json, ", err)
		return
	}
	topic := fmt.Sprint("v2/device/info/need/group/", resultData.Keys.GroupMark, "/code/", deviceData.Code)
	err = MQTTClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		CoreLog.MqttError("mqtt sub device find, mqtt push device info, ", err)
		return
	}
}
