package UserTicket

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitLog(t *testing.T) {
	TestInitTicket(t)
	TestAddTicket(t)
	TestUseTicket(t)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		ConfigID: -1,
		UserID:   -1,
		Mode:     -1,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetAnalysisUse(t *testing.T) {
	data, err := GetAnalysisUse(&ArgsGetAnalysisUse{
		OrgID: newConfigData.OrgID,
		Mode:  1,
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHour().Time),
			MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().Time),
		},
	})
	ToolsTest.ReportData(t, err, data)
	data, err = GetAnalysisUse(&ArgsGetAnalysisUse{
		OrgID: newConfigData.OrgID,
		Mode:  2,
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHour().Time),
			MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().Time),
		},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestClearLog(t *testing.T) {
	TestClearTicket(t)
}
