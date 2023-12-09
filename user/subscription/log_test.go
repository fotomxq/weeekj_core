package UserSubscription

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitLog(t *testing.T) {
	TestInitSub(t)
	TestSetSub(t)
	TestUseSub(t)
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
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetAnalysisUse(t *testing.T) {
	data, err := GetAnalysisUse(&ArgsGetAnalysisUse{
		OrgID: newConfigData.OrgID,
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHour().Time),
			MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().Time),
		},
	})
	ToolsTest.ReportData(t, err, data)
	data, err = GetAnalysisUse(&ArgsGetAnalysisUse{
		OrgID: newConfigData.OrgID,
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHour().Time),
			MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().Time),
		},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAnalysisArea(t *testing.T) {
	data, err := GetAnalysisArea(&ArgsGetAnalysisArea{
		OrgID: newConfigData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestClearLog(t *testing.T) {
	TestClearSubByConfig(t)
	TestDeleteConfig(t)
}
