package UserChat

import (
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
	TestOrg.LocalCreateOrg(t)
}

func TestClear(t *testing.T) {
	TestOrg.LocalClear(t)
}
