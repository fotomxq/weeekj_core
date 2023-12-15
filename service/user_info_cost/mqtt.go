package ServiceUserInfoCost

import (
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	IOTMQTT "github.com/fotomxq/weeekj_core/v5/iot/mqtt"
	MapRoom "github.com/fotomxq/weeekj_core/v5/map/room"
	ServiceUserInfo "github.com/fotomxq/weeekj_core/v5/service/user_info"
)

// 请求信息档案列表
type subUserInfoCostListData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark" check:"mark" empty:"true"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark" check:"mark" empty:"true"`
}

func subUserInfoCostList(client mqtt.Client, message mqtt.Message) {
	logAppend := "mqtt sub get server user info cost list, "
	var resultData subUserInfoCostListData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError(logAppend, "json, ", err)
		return
	}
	if _, _, b := IOTMQTT.SubBefore(resultData.Keys, resultData.OrgID, logAppend); !b {
		return
	}
	dataList, dataCount, err := GetCostList(&ArgsGetCostList{
		Pages:        resultData.Pages,
		OrgID:        resultData.OrgID,
		RoomID:       resultData.RoomID,
		InfoID:       resultData.InfoID,
		RoomBindMark: resultData.RoomBindMark,
		SensorMark:   resultData.SensorMark,
	})
	if err != nil {
		CoreLog.MqttError(logAppend, "get list data, ", err)
		return
	}
	if err := pushUserInfoCostList(resultData.Keys.GroupMark, resultData.Keys.Code, dataList, dataCount); err != nil {
		CoreLog.MqttError(logAppend, "push data, ", err)
		return
	}
}

func pushUserInfoCostList(groupMark, deviceCode string, dataList []FieldsCost, dataCount int64) (err error) {
	//重组数据
	type newDataType struct {
		//数据集合
		DataList []FieldsCost `json:"dataList"`
		//房间信息
		Rooms map[int64]string `json:"rooms"`
		//信息人员列
		InfoNames map[int64]string `json:"infoNames"`
		//配置列
		ConfigNames map[int64]string `json:"configNames"`
	}
	var newDataList newDataType
	if err == nil {
		newDataList.DataList = dataList
		var waitRooms, waitInfo, waitConfigs []int64
		for _, v := range dataList {
			if v.RoomID > 0 {
				waitRooms = append(waitRooms, v.RoomID)
			}
			if v.InfoID > 0 {
				waitInfo = append(waitInfo, v.InfoID)
			}
			if v.ConfigID > 0 {
				waitConfigs = append(waitConfigs, v.ConfigID)
			}
		}
		if len(waitRooms) > 0 {
			newDataList.Rooms, _ = MapRoom.GetRoomsName(&MapRoom.ArgsGetRooms{
				IDs:        waitInfo,
				HaveRemove: false,
			})
		}
		if len(waitInfo) > 0 {
			newDataList.InfoNames, _ = ServiceUserInfo.GetInfoMoreNames(&ServiceUserInfo.ArgsGetInfoMore{
				IDs:        waitInfo,
				HaveRemove: false,
			})
		}
		if len(waitConfigs) > 0 {
			newDataList.ConfigNames, _ = GetConfigsName(&ArgsGetConfigs{
				IDs:        waitConfigs,
				HaveRemove: false,
				OrgID:      -1,
			})
		}
	}
	//获取和打包数据
	var dataByte []byte
	dataByte, err = json.Marshal(newDataList)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := fmt.Sprint("service/user/info/cost/list/group/", groupMark, "/code/", deviceCode)
	if err = IOTMQTT.MQTTClient.PublishWait(topic, 0, false, dataByte); err != nil {
		err = errors.New(fmt.Sprint("mqtt push data, ", err))
		return
	}
	return
}

// 请求最新的一条数据
type subUserInfoCostLastData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomIDs []int64 `db:"room_ids" json:"roomIDs" check:"ids" empty:"true"`
	//信息ID
	InfoIDs []int64 `db:"info_ids" json:"infoIDs" check:"ids" empty:"true"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark" check:"mark" empty:"true"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark" check:"mark" empty:"true"`
}

func subUserInfoCostLast(client mqtt.Client, message mqtt.Message) {
	//前置处理
	logAppend := "mqtt sub get server user info cost last, "
	var resultData subUserInfoCostLastData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError(logAppend, "json, ", err)
		return
	}
	if _, _, b := IOTMQTT.SubBefore(resultData.Keys, resultData.OrgID, logAppend); !b {
		return
	}
	//构建数据
	var dataList []FieldsCost
	for _, v := range resultData.RoomIDs {
		data, err := GetCostLast(&ArgsGetCostLast{
			OrgID:        resultData.OrgID,
			RoomID:       v,
			InfoID:       0,
			RoomBindMark: resultData.RoomBindMark,
			SensorMark:   resultData.SensorMark,
		})
		if err != nil {
			continue
		}
		dataList = append(dataList, data)
	}
	for _, v := range resultData.InfoIDs {
		data, err := GetCostLast(&ArgsGetCostLast{
			OrgID:        resultData.OrgID,
			RoomID:       0,
			InfoID:       v,
			RoomBindMark: resultData.RoomBindMark,
			SensorMark:   resultData.SensorMark,
		})
		if err != nil {
			continue
		}
		dataList = append(dataList, data)
	}
	if err := pushUserInfoCostList(resultData.Keys.GroupMark, resultData.Keys.Code, dataList, int64(len(dataList))); err != nil {
		CoreLog.MqttError(logAppend, "push data, ", err)
		return
	}
}
