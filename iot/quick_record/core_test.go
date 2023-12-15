package IOTQuickRecord

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
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
