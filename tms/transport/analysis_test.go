package TMSTransport

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newAnalysisData FieldsAnalysis
)

func TestAnalysisInit(t *testing.T) {
	TestTransportInit(t)
	TestCreateTransport(t)
}

func TestGetAnalysisAvg(t *testing.T) {
	data, err := GetAnalysisAvg(&ArgsGetAnalysisAvg{
		OrgID:       newBindData.OrgID,
		BindID:      newBindData.BindID,
		InfoID:      0,
		UserID:      0,
		TransportID: 0,
		BetweenTime: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHours(2).Time),
			MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().Time),
		},
		TimeType: "hour",
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAnalysisSUM(t *testing.T) {
	data, err := GetAnalysisSUM(&ArgsGetAnalysisSUM{
		OrgID:       newBindData.OrgID,
		BindID:      newBindData.BindID,
		InfoID:      0,
		UserID:      0,
		TransportID: 0,
		BetweenTime: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHours(2).Time),
			MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().Time),
		},
		TimeType: "hour",
	})
	ToolsTest.ReportData(t, err, data)
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
		BindID:      0,
		InfoID:      0,
		UserID:      0,
		TransportID: 0,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newAnalysisData = dataList[0]
	}
}

func TestUpdateAnalysis(t *testing.T) {
	err := UpdateAnalysis(&ArgsUpdateAnalysis{
		TransportID: newAnalysisData.TransportID,
		InfoID:      newAnalysisData.InfoID,
		UserID:      newAnalysisData.UserID,
		Level:       2,
	})
	ToolsTest.ReportError(t, err)
}

func TestAnalysisClear(t *testing.T) {
	TestDeleteTransport(t)
	TestTransportClear(t)
}
