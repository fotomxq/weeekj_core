package MapRoom

import (
	"encoding/json"
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTMQTT "gitee.com/weeekj/weeekj_core/v5/iot/mqtt"
	OrgRecord "gitee.com/weeekj/weeekj_core/v5/org/record"
)

// 发布紧急呼叫
func pushRoomEmergencyCall(orgID, deviceID, id int64) {
	type newDataType struct {
		ID int64 `json:"id"`
	}
	newData := newDataType{
		ID: id,
	}
	//获取和打包数据
	dataByte, err := json.Marshal(newData)
	if err != nil {
		CoreLog.MqttError("mqtt push data, json error, " + err.Error())
		return
	}
	//记录日志
	haveData, _ := AppendWarning(&ArgsAppendWarning{
		OrgID:    orgID,
		RoomID:   id,
		DeviceID: deviceID,
		NeedMQTT: false,
		CallType: 0,
	})
	if haveData {
		return
	}
	_ = OrgRecord.Create(&OrgRecord.ArgsCreate{
		OrgID:       orgID,
		BindID:      0,
		ContentMark: "map_room_emergency_call",
		Content:     fmt.Sprint("房间[", id, "]发布新的紧急呼叫请求"),
	})
	//推送数据
	topic := fmt.Sprint("map/room/emergency_call_new/org/", orgID)
	err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		CoreLog.MqttError(fmt.Sprint("mqtt push data, ", err))
		return
	}
	return
}

// 发布解除应急呼叫
func pushRoomUnEmergencyCall(orgID, deviceID, id int64) {
	type newDataType struct {
		ID int64 `json:"id"`
	}
	newData := newDataType{
		ID: id,
	}
	//获取和打包数据
	dataByte, err := json.Marshal(newData)
	if err != nil {
		CoreLog.MqttError("mqtt push data, json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("map/room/emergency_call_un/org/", orgID)
	err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte)
	if err != nil {
		CoreLog.MqttError(fmt.Sprint("mqtt push data, ", err))
		return
	}
	return
}
