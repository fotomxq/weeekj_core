package IOTMQTTClient

import (
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestConnect(t *testing.T) {
	runConnect()
}

func TestPushDeviceFindData(t *testing.T) {
	PushDeviceFindData(ArgsPushDeviceFindData{
		Keys: IOTDevice.ArgsCheckDeviceKey{
			GroupMark: "",
			Code:      "",
			NowTime:   0,
			Rand:      "",
			Key:       "",
		},
	})
}
