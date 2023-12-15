package IOTMQTTClient

import (
	"encoding/json"
	"errors"
	"fmt"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
)

// ArgsPushDeviceMissionResult 反馈设备任务结果参数
type ArgsPushDeviceMissionResult struct {
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

// PushDeviceMissionResult 反馈设备任务结果
// device/error
func PushDeviceMissionResult(args ArgsPushDeviceMissionResult) (err error) {
	//打包数据集合
	var dataByte []byte
	dataByte, err = json.Marshal(args)
	if err != nil {
		err = errors.New("json error, " + err.Error())
		return
	}
	//推送数据
	topic := "device/mission/result"
	if err = mqttClient.PublishWait(topic, 0, false, dataByte); err != nil {
		err = errors.New(fmt.Sprint("mqtt push base error, ", err))
		return
	}
	return
}
