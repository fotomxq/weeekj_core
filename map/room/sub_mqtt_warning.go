package MapRoom

import (
	"encoding/json"
	"fmt"

	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTBind "gitee.com/weeekj/weeekj_core/v5/iot/bind"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	OrgRecord "gitee.com/weeekj/weeekj_core/v5/org/record"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 请求紧急呼叫
type subRoomEmergencyCallNeedData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

func subRoomEmergencyCallNeed(_ mqtt.Client, message mqtt.Message) {
	var resultData subRoomEmergencyCallNeedData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub need emergency call, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.MqttError("mqtt sub need emergency call, key, ", err)
		return
	}
	//获取授权
	operateList, err := IOTDevice.GetOperateByDeviceID(&IOTDevice.ArgsGetOperateByDeviceID{
		DeviceID: deviceID,
	})
	var orgID int64
	if err == nil && len(operateList) > 0 {
		orgID = operateList[0].OrgID
	} else {
		CoreLog.MqttError("mqtt sub need emergency call, device no operate, ", err)
		return
	}
	//绑定关系
	bindData, err := IOTBind.GetBindDevice(orgID, deviceID)
	if err != nil {
		CoreLog.MqttError("mqtt sub need emergency call, device no bind, ", err)
		return
	}
	if bindData.FromInfo.System != "room" {
		CoreLog.MqttError("mqtt sub need emergency call, bind from not room, ", err)
		return
	}
	//查询房间合法性
	roomData, err := GetRoomID(&ArgsGetRoomID{
		ID:    bindData.FromInfo.ID,
		OrgID: orgID,
	})
	if err != nil {
		CoreLog.MqttError("mqtt sub need emergency call, data room data, ", err)
		return
	}
	//发布紧急呼叫
	pushRoomEmergencyCall(orgID, deviceID, roomData.ID)
}

// 订阅解除应急呼叫
func subRoomUnEmergencyCall(_ mqtt.Client, message mqtt.Message) {
	type newDataType struct {
		ID int64 `json:"id"`
	}
	var resultData newDataType
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub un emergency call, json, ", err)
		return
	}
	//查询房间合法性
	roomData, err := GetRoomID(&ArgsGetRoomID{
		ID:    resultData.ID,
		OrgID: -1,
	})
	if err != nil {
		CoreLog.MqttError("mqtt sub un emergency call, data room data, ", err)
		return
	}
	//记录日志
	_ = OrgRecord.Create(&OrgRecord.ArgsCreate{
		OrgID:       roomData.OrgID,
		BindID:      0,
		ContentMark: "map_room_un_emergency_call",
		Content:     fmt.Sprint("房间[", roomData.ID, "]发布解除呼叫请求"),
	})
	_ = UnWarning(&ArgsAppendWarning{
		OrgID:    roomData.OrgID,
		RoomID:   roomData.ID,
		DeviceID: 0,
		NeedMQTT: false,
		CallType: 0,
	})
}
