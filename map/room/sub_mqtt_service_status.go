package MapRoom

import (
	"encoding/json"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 请求变更呼叫状态
// 配合聊天室完成呼叫处理请求，设备端主动发起该请求
type subRoomServiceStatusData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//房间ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//服务呼叫状态
	// 0 无呼叫; 1 正在呼叫; 2 已经应答并处置
	// 处置完成后将回归0状态
	ServiceStatus int `db:"service_status" json:"serviceStatus" check:"intThan0" empty:"true"`
	//服务工作人员
	// 如果没有指定，将自动分配
	ServiceBindID int64 `db:"service_bind_id" json:"serviceBindID" check:"id" empty:"true"`
	//联动行政任务
	// 任务完成后将自动清除为0，否则将一直挂起
	// 如果没有指定，将自动生成
	ServiceMissionID int64 `db:"service_mission_id" json:"serviceMissionID" check:"id" empty:"true"`
}

func subRoomServiceStatus(_ mqtt.Client, message mqtt.Message) {
	var resultData subRoomServiceStatusData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub update room service status, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.MqttError("mqtt sub update room service status, key, ", err)
		return
	}
	_, err = IOTDevice.CheckOperate(&IOTDevice.ArgsCheckOperate{
		DeviceID: deviceID,
		OrgID:    resultData.OrgID,
	})
	if err != nil {
		CoreLog.MqttError("mqtt sub update room service status, check operate, ", err)
		return
	}
	//查询房间合法性
	roomData, err := GetRoomID(&ArgsGetRoomID{
		ID:    resultData.ID,
		OrgID: resultData.OrgID,
	})
	if err != nil {
		CoreLog.MqttError("mqtt sub update room service status, data room data, ", err)
		return
	}
	//变更呼叫请求
	_, err = UpdateServiceStatus(&ArgsUpdateServiceStatus{
		ID:               roomData.ID,
		OrgID:            resultData.OrgID,
		ServiceStatus:    resultData.ServiceStatus,
		ServiceBindID:    resultData.ServiceBindID,
		ServiceMissionID: resultData.ServiceMissionID,
	})
	if err != nil {
		CoreLog.MqttError("mqtt sub update room service status, update status, ", err)
		return
	}
}
