package AnalysisIndexEvent

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
	if err := Init(); err != nil {
		t.Error(err)
		return
	}
}

func TestInsertEvent(t *testing.T) {
	err := InsertEvent(&ArgsInsertEvent{
		Code:       "test",
		YearMD:     "2021-01-01",
		Level:      1,
		FromSystem: "test",
		FromID:     1,
		FromType:   "FromType",
		Extend1:    "Extend1",
		Extend2:    "Extend2",
		Extend3:    "Extend3",
		Extend4:    "Extend4",
		Extend5:    "Extend5",
		Threshold:  100,
		IndexVal:   100,
		Remark:     "Remark",
	})
	if err != nil {
		t.Log("TestInsertEvent:", err)
	}
}

func TestClear(t *testing.T) {
}
