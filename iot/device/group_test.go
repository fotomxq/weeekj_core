package IOTDevice

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newGroupData FieldsGroup
)

func TestInitGroup(t *testing.T) {
	TestInit(t)
}

func TestCreateGroup(t *testing.T) {
	TestCreateAction(t)
	var err error
	newGroupData, err = CreateGroup(&ArgsCreateGroup{
		Mark:       "test_mark",
		Name:       "test",
		Des:        "test des",
		CoverFiles: []int64{},
		Action:     []int64{newActionData.ID},
		ExpireTime: 60,
		UseType:    0,
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newGroupData)
}

func TestGetGroupList(t *testing.T) {
	dataList, dataCount, err := GetGroupList(&ArgsGetGroupList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetGroupByID(t *testing.T) {
	data, err := GetGroupByID(&ArgsGetGroupByID{
		ID: newGroupData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetGroupMore(t *testing.T) {
	dataList, err := GetGroupMore(&ArgsGetGroupMore{
		IDs:        []int64{newGroupData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetGroupMoreMap(t *testing.T) {
	data, err := GetGroupMoreMap(&ArgsGetGroupMore{
		IDs:        []int64{newGroupData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateGroup(t *testing.T) {
	err := UpdateGroup(&ArgsUpdateGroup{
		ID:         newGroupData.ID,
		Mark:       "test_mark",
		Name:       "test",
		Des:        "test des",
		CoverFiles: []int64{},
		Action:     []int64{},
		ExpireTime: 60,
		UseType:    0,
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteGroup(t *testing.T) {
	err := DeleteGroup(&ArgsDeleteGroup{
		ID: newGroupData.ID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteAction(t)
}
