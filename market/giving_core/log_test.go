package MarketGivingCore

import (
	"testing"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
)

var (
	newLogData FieldsLog
)

func TestInitLog(t *testing.T) {
	TestInitConfig(t)
	TestCreateConfig(t)
}

func TestCreateLog(t *testing.T) {
	var err error
	var errCode string
	newLogData, errCode, err = CreateLog(&ArgsCreateLog{
		OrgID: TestOrg.OrgData.ID,
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "mall",
			ID:     1,
			Mark:   "",
			Name:   "测试商品",
		},
		UserID:         TestOrg.UserInfo.ID,
		ReferrerUserID: 0,
		ReferrerBindID: 0,
		ConfigID:       newConfigData.ID,
		PriceTotal:     1500,
		Des:            "奖励测试",
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
		OrgID:          -1,
		FromInfo:       CoreSQLFrom.FieldsFrom{},
		UserID:         -1,
		ReferrerUserID: -1,
		ReferrerBindID: -1,
		ConfigID:       -1,
		IsRemove:       false,
		Search:         "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
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

func TestDeleteLog(t *testing.T) {
	err := DeleteLog(&ArgsDeleteLog{
		ID:    newLogData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearLog(t *testing.T) {
	TestClearConfig(t)
}
