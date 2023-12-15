package IOTMQTT

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	ToolsAppUpdate "github.com/fotomxq/weeekj_core/v5/tools/app_update"
)

// subDeviceOnline 设备在线情况更正
type subAppUpdateData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//运行环境
	// android_phone / android_pad / ios_phone / ios_ipad / windows / osx / linux
	// 或者特定品牌的定制
	System string `db:"system" json:"system"`
	//环境的最低版本
	// 如果给与指定专供版本，则该设定无效
	// [7, 1, 4] => version 7.1.4
	SystemVersion string `db:"system_version" json:"systemVersion"`
	//APP标识码
	AppMark string `db:"app_mark" json:"appMark"`
	//版本号
	// [7, 1, 4] => version 7.1.4
	Version string `db:"version" json:"version"`
}

func subAppUpdate(_ mqtt.Client, message mqtt.Message) {
	var resultData subAppUpdateData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.MqttError("mqtt sub app update, json, ", err)
		return
	}
	if err := IOTDevice.CheckDeviceKey(&resultData.Keys); err != nil {
		CoreLog.MqttError("mqtt sub app update, key, ", err)
		return
	}
	data, b := ToolsAppUpdate.CheckUpdate(&ToolsAppUpdate.ArgsCheckUpdate{
		System:        resultData.System,
		SystemVersion: resultData.SystemVersion,
		AppMark:       resultData.AppMark,
		Version:       resultData.Version,
	})
	if b {
		if err := PushAppUpdate(resultData.Keys.GroupMark, resultData.Keys.Code, "", data); err != nil {
			CoreLog.MqttError("mqtt sub app update, push update data, ", err)
			return
		}
	}
}
