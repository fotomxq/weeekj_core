package AnalysisAny

import (
	"testing"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if !isInit {
		//ToolsTest.Init(t)
		isInit = true
	}
	//TestOrg.LocalCreateBind(t)
}

func TestClear(t *testing.T) {
	//TestOrg.LocalClear(t)
}
