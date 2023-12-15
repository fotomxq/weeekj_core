package ServiceInfoExchange

import (
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	Init()
	if !isInit {
		ToolsTest.Init(t)
		OrgCore.Init()
	}
	isInit = true
	TestOrg.LocalCreateBind(t)
}

func TestClear(t *testing.T) {
	TestOrg.LocalClear(t)
}
