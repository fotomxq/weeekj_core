package BlogUserRead

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

func TestInitAnalysis(t *testing.T) {
	TestInitLog(t)
	TestCreateLog(t)
}

func TestGetAnalysisList(t *testing.T) {
	dataList, dataCount, err := GetAnalysisList(&ArgsGetAnalysisList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:       -1,
		ChildOrgID:  -1,
		UserID:      -1,
		FromMark:    "",
		FromName:    "",
		IP:          "",
		SortID:      -1,
		ReadTimeMin: -1,
		ReadTimeMax: -1,
		TimeBetween: CoreSQLTime.DataCoreTime{},
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetAnalysisGroupTime(t *testing.T) {
	dataList, err := GetAnalysisGroupTime(&ArgsGetAnalysisGroupTime{
		OrgID:      -1,
		ChildOrgID: -1,
		SortID:     -1,
		TimeType:   "month",
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetAnalysisGroupChildOrgList(t *testing.T) {
	dataList, dataCount, err := GetAnalysisGroupChildOrgList(&ArgsGetAnalysisList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "read_count",
			Desc: false,
		},
		OrgID:       -1,
		ChildOrgID:  -1,
		UserID:      -1,
		FromMark:    "",
		FromName:    "",
		IP:          "",
		SortID:      -1,
		ReadTimeMin: -1,
		ReadTimeMax: -1,
		TimeBetween: CoreSQLTime.DataCoreTime{},
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetAnalysisGroupUserList(t *testing.T) {
	dataList, dataCount, err := GetAnalysisGroupUserList(&ArgsGetAnalysisList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "read_count",
			Desc: false,
		},
		OrgID:       -1,
		ChildOrgID:  -1,
		UserID:      -1,
		FromMark:    "",
		FromName:    "",
		IP:          "",
		SortID:      -1,
		ReadTimeMin: -1,
		ReadTimeMax: -1,
		TimeBetween: CoreSQLTime.DataCoreTime{},
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetAnalysisCount(t *testing.T) {
	data, err := GetAnalysisCount(&ArgsGetAnalysisCount{
		OrgID:       TestOrg.OrgData.ID,
		ChildOrgID:  0,
		UserID:      TestOrg.UserInfo.ID,
		FromMark:    "",
		FromName:    "",
		IP:          "",
		SortID:      -1,
		TimeBetween: CoreSQLTime.DataCoreTime{},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAnalysisTime(t *testing.T) {
	data, err := GetAnalysisTime(&ArgsGetAnalysisCount{
		OrgID:      TestOrg.OrgData.ID,
		ChildOrgID: -1,
		UserID:     -1,
		FromMark:   "web",
		FromName:   "",
		IP:         "",
		SortID:     -1,
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHours(3).Time),
			MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().AddMinutes(1).Time),
		},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAnalysisAvgReadTime(t *testing.T) {
	data, err := GetAnalysisAvgReadTime(&ArgsGetAnalysisAvgReadTime{
		OrgID:       TestOrg.OrgData.ID,
		ChildOrgID:  0,
		UserID:      TestOrg.UserInfo.ID,
		FromMark:    "",
		FromName:    "",
		IP:          "",
		SortID:      -1,
		ContentID:   -1,
		TimeBetween: CoreSQLTime.DataCoreTime{},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestClearAnalysis(t *testing.T) {
	TestClearLog(t)
}
