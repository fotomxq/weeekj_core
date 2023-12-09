package IOTMQTT

import (
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
	"time"
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
	time.Sleep(time.Second * 5)
}
