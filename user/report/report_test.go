package UserReport

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newReportData FieldsReport
)

func TestInitReport(t *testing.T) {
	TestInit(t)
}

func TestCreateReport(t *testing.T) {
	err := CreateReport(&ArgsCreateReport{
		OrgID:  TestOrg.OrgData.ID,
		UserID: TestOrg.UserInfo.ID,
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "mall",
			ID:     1,
			Mark:   "",
			Name:   "",
		},
		IP:         "0.0.0.0",
		UserName:   TestOrg.UserInfo.Name,
		NationCode: TestOrg.UserInfo.NationCode,
		Phone:      TestOrg.UserInfo.Phone,
		Email:      TestOrg.UserInfo.Email,
		Files:      []int64{},
		Content:    "测试举报",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetReportList(t *testing.T) {
	dataList, dataCount, err := GetReportList(&ArgsGetReportList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newReportData = dataList[0]
	}
}

func TestGetReport(t *testing.T) {
	data, err := GetReport(&ArgsGetReport{
		ID:    newReportData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestReportData(t *testing.T) {
	err := ReportData(&ArgsReportData{
		ID:            newReportData.ID,
		OrgID:         -1,
		ReportUserID:  newReportData.UserID,
		ReportContent: "反馈内容",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteReport(t *testing.T) {
	err := DeleteReport(&ArgsDeleteReport{
		ID:    newReportData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearReport(t *testing.T) {
	TestClear(t)
}
