package BasePageStyle

import (
	"testing"

	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
		OrgCore.Init(true, true)
	}
	isInit = true
	TestOrg.LocalCreateBind(t)
}

func TestClear(t *testing.T) {
	TestOrg.LocalClear(t)
}
