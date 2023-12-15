package IOTDevice

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newActionData FieldsAction
)

func TestInitAction(t *testing.T) {
	TestInit(t)
}

func TestCreateAction(t *testing.T) {
	var err error
	newActionData, err = CreateAction(&ArgsCreateAction{
		Mark:        "test_mark",
		Name:        "测试",
		Des:         "test mark",
		ExpireTime:  60,
		ConnectType: "none",
		Configs:     []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newActionData)
}

func TestGetActionList(t *testing.T) {
	dataList, dataCount, err := GetActionList(&ArgsGetActionList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetActionMore(t *testing.T) {
	data, err := GetActionMore(&ArgsGetActionMore{
		IDs:        []int64{newActionData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetActionMoreMap(t *testing.T) {
	data, err := GetActionMoreMap(&ArgsGetActionMore{
		IDs:        []int64{newActionData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateAction(t *testing.T) {
	err := UpdateAction(&ArgsUpdateAction{
		ID:          newActionData.ID,
		Mark:        "test_mark",
		Name:        "测试",
		Des:         "test mark",
		ExpireTime:  60,
		ConnectType: "none",
		Configs:     []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteAction(t *testing.T) {
	err := DeleteAction(&ArgsDeleteAction{
		ID: newActionData.ID,
	})
	ToolsTest.ReportError(t, err)
}
