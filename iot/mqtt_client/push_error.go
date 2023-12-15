package IOTMQTTClient

import (
	"encoding/json"
	"errors"
	"fmt"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
)

// ArgsPushError 广播错误采集器参数
type ArgsPushError struct {
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

// PushError 广播错误采集器参数
// device/error
func PushError(args ArgsPushError) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := "device/error"
	if err = mqttClient.PublishWait(topic, 0, false, dataByte); err != nil {
		err = errors.New(fmt.Sprint("mqtt push base error, ", err))
		return
	}
	return
}
