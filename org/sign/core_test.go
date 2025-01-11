package OrgSign

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit bool
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
	TestOrg.LocalInit()
	TestOrg.LocalCreateBind(t)
	if err := Init(); err != nil {
		return
	}
}

func TestClear(t *testing.T) {
	TestOrg.LocalClear(t)
}
