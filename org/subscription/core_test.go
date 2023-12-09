package OrgSubscription

import (
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit        = false
	newConfigData FieldsConfig
	newSub        FieldsSub
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
	TestOrg.LocalCreateOrg(t)
}

func TestClear(t *testing.T) {
	TestOrg.LocalClear(t)
}
