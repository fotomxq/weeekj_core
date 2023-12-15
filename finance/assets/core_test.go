package FinanceAssets

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit = false
	//虚构来源信息
	userID    int64 = 123
	orgID     int64 = 345
	orgBindID int64 = 234
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestClear(t *testing.T) {
}
