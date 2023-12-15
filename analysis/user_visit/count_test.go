package AnalysisUserVisit

import (
	"testing"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

func TestCountInit(t *testing.T) {
	TestInit(t)
}

func TestCreateCount(t *testing.T) {
	err := CreateCount(&ArgsCreateCount{
		OrgID: 1,
		Mark:  0,
		Count: 1,
	})
	ToolsTest.ReportError(t, err)
	err = CreateCount(&ArgsCreateCount{
		OrgID: 1,
		Mark:  0,
		Count: 3,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetCountAnalysis(t *testing.T) {
	data, err := GetCountAnalysis(&ArgsGetCountAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubHour().Time,
			MaxTime: CoreFilter.GetNowTimeCarbon().Time,
		},
		TimeType: "hour",
		OrgID:    1,
		Mark:     0,
	})
	ToolsTest.ReportData(t, err, data)
}
