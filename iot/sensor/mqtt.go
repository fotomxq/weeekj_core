package IOTSensor

import (
	"encoding/json"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type subCreateData struct {
	//配对密钥
	Keys IOTDevice.ArgsCheckDeviceKey `json:"keys"`
	//创建时间
	// 如果给空，则以当前时间为主
	// IOS时间
	CreateAt string `db:"create_at" json:"createAt"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//数据
	Data  int64   `db:"data" json:"data"`
	DataF float64 `db:"data_f" json:"dataF"`
	DataS string  `db:"data_s" json:"dataS"`
}

func subCreate(_ mqtt.Client, message mqtt.Message) {
	var resultData subCreateData
	resultByte := message.Payload()
	if err := json.Unmarshal(resultByte, &resultData); err != nil {
		CoreLog.Error("mqtt sub create sensor, json, ", err)
		return
	}
	deviceID, err := IOTDevice.CheckDeviceKeyAndDeviceID(&resultData.Keys)
	if err != nil {
		CoreLog.Error("mqtt sub create sensor, key, ", err)
		return
	}
	//创建数据
	err = Create(&ArgsCreate{
		CreateAt: resultData.CreateAt,
		DeviceID: deviceID,
		Mark:     resultData.Mark,
		Data:     resultData.Data,
		DataF:    resultData.DataF,
		DataS:    resultData.DataS,
	})
	if err != nil {
		CoreLog.Error("mqtt sub create sensor, create, ", err)
		return
	}
}
