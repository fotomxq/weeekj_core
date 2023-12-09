package AnalysisAny

import (
	"fmt"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	"testing"

	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
)

func TestInitAny(t *testing.T) {
	TestInit(t)
	TestInitConfigInit(t)
}

func TestAppendAny(t *testing.T) {
	err := AppendAny(&ArgsAppendAny{
		CreateAt: "",
		OrgID:    1,
		UserID:   2,
		BindID:   3,
		Mark:     "test",
		Data:     1,
		DataVal:  "test data",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetAnyByMark(t *testing.T) {
	data, err := GetAnyByMark(&ArgsGetAnyByMark{
		OrgID: 1,
		Mark:  "test",
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAnyInt64ByMark(t *testing.T) {
	args := ArgsGetAnyByMark{
		OrgID:  1,
		UserID: 2,
		Mark:   "mark",
		BindID: 3,
		Param1: 4,
		Param2: 5,
		BetweenTime: CoreSQLTime.DataCoreTime{
			MinTime: "2022-02",
			MaxTime: "2022-03",
		},
	}
	cacheMark := fmt.Sprint("analysis_any-get-", args)
	t.Log("cacheMark: ", cacheMark)
}

func TestClearAny(t *testing.T) {
	TestClear(t)
}
