package ToolsCommunication

import (
	"encoding/json"
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	BaseQiniu "github.com/fotomxq/weeekj_core/v5/base/qiniu"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	IOTMQTT "github.com/fotomxq/weeekj_core/v5/iot/mqtt"
)

//广播房间的关闭或删除\退出

// 请求发起新的双向聊天
type subNewRoomData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	ConnectType int `db:"connect_type" json:"connectType" check:"intThan0" empty:"true"`
	//通讯类型
	DataType int `db:"data_type" json:"dataType" check:"intThan0"`
	//昵称
	FromName string `db:"from_name" json:"fromName" check:"name"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	FromConnectType int `db:"from_connect_type" json:"fromConnectType" check:"intThan0" empty:"true"`
	//到达系统
	ToSystem int   `db:"to_system" json:"toSystem" check:"intThan0"`
	ToID     int64 `db:"to_id" json:"toID" check:"id"`
	//到达昵称
	ToName string `db:"to_name" json:"toName" check:"name"`
	//链接方式
	// 0 系统自带TCP握手方式; 1 系统自带RTC方式; 2 第三方agora服务
	ToConnectType int `db:"to_connect_type" json:"toConnectType" check:"intThan0" empty:"true"`
}

func subNewRoom(client mqtt.Client, message mqtt.Message) {
	//解析数据
	var resultData subNewRoomData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.Error("mqtt sub append room two, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.Error("mqtt sub append room two, key, ", err)
		return
	}
	//验证设备是否为该组织具备
	_, err = IOTDevice.CheckOperate(&IOTDevice.ArgsCheckOperate{
		DeviceID: deviceID,
		OrgID:    resultData.OrgID,
	})
	if err != nil {
		CoreLog.Error("mqtt sub append room two, check device operate by org, ", err)
		return
	}
	//构建聊天握手
	roomData, fromData, toData, err := AppendRoomTwo(&ArgsAppendRoomTwo{
		OrgID:           resultData.OrgID,
		ConnectType:     resultData.ConnectType,
		DataType:        resultData.DataType,
		SortID:          0,
		Tags:            []int64{},
		Name:            fmt.Sprint("设备[", deviceID, "]聊天室"),
		Des:             "",
		CoverFileID:     0,
		Password:        "",
		FromSystem:      2,
		FromID:          deviceID,
		FromName:        resultData.FromName,
		FromConnectType: resultData.FromConnectType,
		ToSystem:        resultData.ToSystem,
		ToID:            resultData.ToID,
		ToName:          resultData.ToName,
		ToConnectType:   resultData.ToConnectType,
	})
	if err != nil {
		CoreLog.Error("mqtt sub append room two, create new room and two from data, ", err)
		return
	}
	//推送设备1
	if err := pushNew(deviceID, roomData, fromData); err != nil {
		CoreLog.Error("mqtt sub append room two, push new room data to 1, ", err)
		return
	}
	//推送设备2
	if resultData.ToSystem == 2 {
		if err := pushNew(resultData.ToID, roomData, toData); err != nil {
			CoreLog.Error("mqtt sub append room two, push new room data to 2, ", err)
			return
		}
	}
}

// 订阅聊天室信息
type subRoomInfoData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//房屋ID
	RoomID int64 `db:"room_id" json:"roomID"`
}

func subRoomInfo(client mqtt.Client, message mqtt.Message) {
	//解析数据
	var resultData subRoomInfoData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.Error("mqtt sub sub room info, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.Error("mqtt sub sub room info, key, ", err)
		return
	}
	//验证设备是否为该组织具备
	_, err = IOTDevice.CheckOperate(&IOTDevice.ArgsCheckOperate{
		DeviceID: deviceID,
		OrgID:    resultData.OrgID,
	})
	if err != nil {
		CoreLog.Error("mqtt sub sub room info, check operate, ", err)
		return
	}
	//验证此设备和聊天室房间ID的关联性
	fromData, err := CheckFromAndRoom(&ArgsCheckFromAndRoom{
		RoomID:     resultData.RoomID,
		FromSystem: 2,
		FromID:     deviceID,
	})
	if err != nil {
		CoreLog.Error("mqtt sub sub room info, check from and room, ", err)
		return
	}
	//获取房间数据
	roomData, err := GetRoom(&ArgsGetRoom{
		ID:    fromData.RoomID,
		OrgID: -1,
	})
	if err != nil {
		CoreLog.Error("mqtt sub sub room info, get room data, ", err)
		return
	}
	//推送数据
	if err := pushRoomInfo(roomData); err != nil {
		CoreLog.Error("mqtt sub sub room info, push room data, ", err)
		return
	}
}

// 推送聊天室信息
func pushRoomInfo(data FieldsRoom) (err error) {
	type dataType struct {
		Data FieldsRoom `json:"data"`
		//文件集
		// id => url
		FileList map[int64]string `json:"fileList"`
	}
	newData := dataType{
		Data: data,
	}
	var waitFiles []int64
	if data.CoverFileID > 0 {
		waitFiles = append(waitFiles, data.CoverFileID)
	}
	if len(waitFiles) > 0 {
		newData.FileList, _ = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
			ClaimIDList: waitFiles,
			UserID:      0,
			OrgID:       0,
			IsPublic:    true,
		})
	}
	//获取和打包数据
	var dataByte []byte
	dataByte, err = json.Marshal(newData)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("tools/communication/room/info/", data.ID)
	err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}

// 请求设备关联的聊天室
type subRoomByFromData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
}

func subRoomByFrom(client mqtt.Client, message mqtt.Message) {
	//解析数据
	var resultData subRoomByFromData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.Error("mqtt sub sub room info, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.Error("mqtt sub sub room info, key, ", err)
		return
	}
	//获取该设备关联的房间数据包
	fromList, err := GetFrom(&ArgsGetFrom{
		FromSystem: 2,
		FromID:     deviceID,
	})
	if err != nil {
		CoreLog.Error("mqtt sub sub room info, check from and room, ", err)
		return
	}
	//获取房间列
	var waitRooms []int64
	for _, v := range fromList {
		if v.RoomID > 0 {
			waitRooms = append(waitRooms, v.RoomID)
		}
	}
	var roomList []FieldsRoom
	if len(waitRooms) > 0 {
		roomList, _ = GetRoomMore(&ArgsGetRoomMore{
			IDs:        waitRooms,
			HaveRemove: false,
		})
	}
	//推送数据
	if err := pushRoomByFrom(deviceID, roomList, fromList); err != nil {
		CoreLog.Error("mqtt sub sub room info, push room data, ", err)
		return
	}
}

// 请求续约聊天室
type subExpireData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id"`
}

func subExpire(client mqtt.Client, message mqtt.Message) {
	//解析数据
	var resultData subExpireData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.Error("mqtt sub sub room info, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.Error("mqtt sub sub room expire, key, ", err)
		return
	}
	//续租更新
	err = UpdateFromExpire(&ArgsUpdateFromExpire{
		RoomID:     resultData.RoomID,
		FromSystem: 2,
		FromID:     deviceID,
	})
	if err != nil {
		CoreLog.Error("mqtt sub sub room expire, update, ", err)
		return
	} else {
		CoreLog.Info("mqtt sub sub room expire")
	}
}

// 给设备推送新的房间
func pushNew(deviceID int64, roomData FieldsRoom, fromData FieldsFrom) (err error) {
	type dataType struct {
		//房间
		RoomData FieldsRoom `json:"roomData"`
		//成员
		FromData FieldsFrom `json:"fromData"`
	}
	newData := dataType{
		RoomData: roomData,
		FromData: fromData,
	}
	//获取和打包数据
	var dataByte []byte
	dataByte, err = json.Marshal(newData)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("tools/communication/room/new/", deviceID)
	err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}

// 推送设备相关的聊天室
func pushRoomByFrom(deviceID int64, roomList []FieldsRoom, fromList []FieldsFrom) (err error) {
	type dataType struct {
		//房间列
		RoomList []FieldsRoom `json:"roomList"`
		//成员列
		FromList []FieldsFrom `json:"fromList"`
		//文件集
		// id => url
		FileList map[int64]string `json:"fileList"`
	}
	newData := dataType{
		RoomList: roomList,
		FromList: fromList,
	}
	var waitFiles []int64
	for _, v := range roomList {
		if v.CoverFileID > 0 {
			waitFiles = append(waitFiles, v.CoverFileID)
		}
	}
	if len(waitFiles) > 0 {
		newData.FileList, _ = BaseQiniu.GetPublicURLsMap(&BaseQiniu.ArgsGetPublicURLs{
			ClaimIDList: waitFiles,
			UserID:      0,
			OrgID:       0,
			IsPublic:    true,
		})
	}
	//获取和打包数据
	var dataByte []byte
	dataByte, err = json.Marshal(newData)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("tools/communication/room/from/", deviceID)
	err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}

// 推送删除房间或成员退出房间
func pushRoomOrFromDelete(roomID int64, fromID int64) (err error) {
	type dataType struct {
		RoomID int64 `json:"roomID"`
		FromID int64 `json:"fromID"`
	}
	newData := dataType{
		RoomID: roomID,
		FromID: fromID,
	}
	//获取和打包数据
	var dataByte []byte
	dataByte, err = json.Marshal(newData)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("tools/communication/room/delete/", roomID)
	err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte)
	return
}

// 退出聊天室
type subOutRoomData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//来源ID
	ID int64 `db:"id" json:"id" check:"id"`
}

func subOutRoom(client mqtt.Client, message mqtt.Message) {
	//解析数据
	var resultData subOutRoomData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.Error("mqtt sub out room, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.Error("mqtt sub out room, key, ", err)
		return
	}
	//退出房间
	err = OutRoom(&ArgsOutRoom{
		ID:         resultData.ID,
		FromSystem: 2,
		FromID:     deviceID,
	})
	if err != nil {
		CoreLog.Error("mqtt sub out room, out room, ", err)
		return
	}
}
