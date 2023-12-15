package IOTMQTT

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	IOTLog "github.com/fotomxq/weeekj_core/v5/iot/log"
	IOTMission "github.com/fotomxq/weeekj_core/v5/iot/mission"
)

type subDeviceMissionResultData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//设备ID
	ID int64 `db:"id" json:"id" check:"id"`
	//任务状态
	// 0 wait 等待发起 / 1 send 已经发送 / 3 failed 已经失败 / 4 cancel 取消
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//回收数据
	// 回收数据如果过大，将不会被存储到本地
	ReportData []byte `db:"report_data" json:"reportData"`
	//行为标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//日志内容
	Content string `db:"content" json:"content" check:"des" min:"1" max:"1000"`
}

func subDeviceMissionResult(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceMissionResultData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device update mission, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub device update mission, key, ", err)
		return
	}
	if resultData.Status == 2 {
		if err := IOTMission.UpdateMissionFinish(&IOTMission.ArgsUpdateMissionFinish{
			ID:         resultData.ID,
			ReportData: resultData.ReportData,
		}); err != nil {
			CoreLog.MqttError("mqtt sub device update mission, json, ", err)
			return
		} else {
			deviceData, err := IOTDevice.GetOperateAndDevice(&IOTDevice.ArgsGetOperateAndDevice{
				DeviceID: resultData.ID,
			})
			if err != nil {
				CoreLog.MqttError("mqtt sub device update mission, device or operate not exist, ", err)
				return
			}
			for _, v := range deviceData {
				IOTLog.Append(&IOTLog.ArgsAppend{
					OrgID:    v.OrgID,
					GroupID:  v.GroupID,
					DeviceID: v.ID,
					Mark:     "finish",
					Content:  "完成目标任务",
				})
			}
		}
	} else {
		if err := IOTMission.UpdateMissionStatus(&IOTMission.ArgsUpdateMissionStatus{
			ID:     resultData.ID,
			Status: resultData.Status,
		}); err != nil {
			CoreLog.MqttError("mqtt sub device update mission, json, ", err)
			return
		} else {
			deviceData, err := IOTDevice.GetOperateAndDevice(&IOTDevice.ArgsGetOperateAndDevice{
				DeviceID: resultData.ID,
			})
			if err != nil {
				CoreLog.MqttError("mqtt sub device update mission, device or operate not exist, ", err)
				return
			}
			for _, v := range deviceData {
				IOTLog.Append(&IOTLog.ArgsAppend{
					OrgID:    v.OrgID,
					GroupID:  v.GroupID,
					DeviceID: v.ID,
					Mark:     resultData.Mark,
					Content:  resultData.Content,
				})
			}
		}
	}
}
