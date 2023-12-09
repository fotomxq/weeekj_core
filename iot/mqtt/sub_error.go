package IOTMQTT

import (
	"encoding/json"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	IOTError "gitee.com/weeekj/weeekj_core/v5/iot/error"
	IOTLog "gitee.com/weeekj/weeekj_core/v5/iot/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// subDeviceError 推送错误信息
type subDeviceErrorData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//是否推送了预警信息
	SendEW bool `db:"send_ew" json:"sendEW"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//错误标识码
	Code string `db:"code" json:"code"`
	//日志内容
	Content string `db:"content" json:"content"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

func subDeviceError(client mqtt.Client, message mqtt.Message) {
	var resultData subDeviceErrorData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub device create error, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub device create error, key, ", err)
		return
	}
	deviceData, err := IOTDevice.GetOperateAndDevice(&IOTDevice.ArgsGetOperateAndDevice{
		DeviceID: resultData.DeviceID,
	})
	if err != nil {
		CoreLog.MqttError("mqtt sub device create error, device or operate not exist, ", err)
		return
	}
	for _, v := range deviceData {
		if err := IOTError.Create(&IOTError.ArgsCreate{
			SendEW:   resultData.SendEW,
			OrgID:    v.OrgID,
			GroupID:  v.GroupID,
			DeviceID: v.ID,
			Code:     resultData.Code,
			Content:  resultData.Content,
			Params:   resultData.Params,
		}); err != nil {
			CoreLog.MqttError("mqtt sub device create error, ", err)
			return
		}
		IOTLog.Append(&IOTLog.ArgsAppend{
			OrgID:    v.OrgID,
			GroupID:  v.GroupID,
			DeviceID: v.ID,
			Mark:     "error",
			Content:  "设备发生异常",
		})
	}
}
