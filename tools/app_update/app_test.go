package ToolsAppUpdate

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newAppData FieldsApp
)

func TestInitApp(t *testing.T) {
	TestInit(t)
}

func TestCreateApp(t *testing.T) {
	var err error
	newAppData, err = CreateApp(&ArgsCreateApp{
		OrgID:    0,
		Name:     "配送系统PDA测试版",
		Des:      "transport_pda des",
		DesFiles: []int64{},
		AppMark:  "transport_pda",
	})
	ToolsTest.ReportData(t, err, newAppData)
}

func TestGetAppID(t *testing.T) {
	data, err := GetAppID(&ArgsGetAppID{
		ID: newAppData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAppList(t *testing.T) {
	dataList, dataCount, err := GetAppList(&ArgsGetAppList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:  0,
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateApp(t *testing.T) {
	err := UpdateApp(&ArgsUpdateApp{
		ID:       newAppData.ID,
		OrgID:    newAppData.OrgID,
		Name:     newAppData.Name,
		Des:      newAppData.Des,
		DesFiles: newAppData.DesFiles,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteApp(t *testing.T) {
	err := DeleteApp(&ArgsDeleteApp{
		ID:    newAppData.ID,
		OrgID: newAppData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
