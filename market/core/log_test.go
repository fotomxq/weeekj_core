package MarketCore

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newLogData FieldsLog
)

func TestInitLog(t *testing.T) {
	TestInitBind(t)
	TestCreateBind(t)
}

func TestCreateLog(t *testing.T) {
	var err error
	var errCode string
	newLogData, errCode, err = CreateLog(&ArgsCreateLog{
		OrgID:      TestOrg.OrgData.ID,
		UserID:     TestOrg.UserInfo.ID,
		BindID:     TestOrg.BindData.ID,
		BindUserID: TestOrg.UserInfo.ID,
		ConfigID:   newConfigData.ID,
		PriceTotal: 1500,
		Des:        "测试推荐",
	})
	ToolsTest.ReportData(t, err, newLogData)
	if err != nil {
		t.Error("errCode: ", errCode)
	}
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      -1,
		UserID:     -1,
		BindID:     -1,
		BindUserID: -1,
		ConfigID:   -1,
		IsRemove:   false,
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetAnalysisBind(t *testing.T) {
	data, err := GetAnalysisBind(&ArgsGetAnalysisBind{
		OrgID:  newConfigData.OrgID,
		UserID: -1,
		BindID: -1,
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubHour().Time,
			MaxTime: CoreFilter.GetNowTimeCarbon().Time,
		},
		ConfigID: -1,
		TimeType: "day",
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAnalysisPriceBind(t *testing.T) {
	data, err := GetAnalysisPriceBind(&ArgsGetAnalysisPriceBind{
		OrgID: newConfigData.OrgID,
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubHour().Time,
			MaxTime: CoreFilter.GetNowTimeCarbon().Time,
		},
		ConfigID: -1,
		TimeType: "day",
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAnalysisPriceTotal(t *testing.T) {
	data, err := GetAnalysisPriceTotal(&ArgsGetAnalysisPriceTotal{
		OrgID: newConfigData.OrgID,
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubHour().Time,
			MaxTime: CoreFilter.GetNowTimeCarbon().Time,
		},
		ConfigID: -1,
		TimeType: "day",
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAnalysisCountBind(t *testing.T) {
	dataList, dataCount, err := GetAnalysisCountBind(&ArgsGetAnalysisCountBind{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "count_count",
			Desc: true,
		},
		OrgID:    newConfigData.OrgID,
		ConfigID: -1,
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubHour().Time,
			MaxTime: CoreFilter.GetNowTimeCarbon().Time,
		},
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDeleteLog(t *testing.T) {
	err := DeleteLog(&ArgsDeleteLog{
		ID:    newLogData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearLog(t *testing.T) {
	TestDeleteBind(t)
	TestClearBind(t)
}
