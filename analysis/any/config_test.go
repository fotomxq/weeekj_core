package AnalysisAny

import (
	"testing"

	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

func TestInitConfig(t *testing.T) {
	TestInit(t)
}

func TestInitConfigInit(t *testing.T) {
	err := InitConfig(&ArgsInitConfig{
		Mark:     "test",
		FileDay:  3,
		MqttOrg:  true,
		MqttUser: false,
		MqttBind: false,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearConfig(t *testing.T) {
	TestClear(t)
}
