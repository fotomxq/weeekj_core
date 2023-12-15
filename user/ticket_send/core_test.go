package UserTicketSend

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
		isInit = true
	}
	TestOrg.LocalCreateOrg(t)
}

func TestClear(t *testing.T) {
	TestOrg.LocalClear(t)
}
