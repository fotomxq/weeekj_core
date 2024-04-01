package BaseBPM

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newThemeData FieldsTheme
)

func TestThemeInit(t *testing.T) {
	TestThemeCategoryInit(t)
	TestCreateThemeCategory(t)
	TestGetThemeCategoryByID(t)
}

func TestCreateTheme(t *testing.T) {
	newDataID, err := CreateTheme(&ArgsCreateTheme{
		CategoryID:  newThemeCategory.ID,
		Name:        "测试主题",
		Description: "测试主题描述",
	})
	ToolsTest.ReportData(t, err, newDataID)
	newThemeData.ID = newDataID
}

func TestGetThemeByID(t *testing.T) {
	data, err := GetThemeByID(&ArgsGetThemeByID{
		ID: newThemeData.ID,
	})
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		t.Error("get theme failed, id: ", newThemeData.ID)
	}
	newThemeData = data
}

func TestGetThemeList(t *testing.T) {
	dataList, dataCount, err := GetThemeList(&ArgsGetThemeList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		CategoryID: newThemeCategory.ID,
		IsRemove:   false,
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetThemeCountByCategoryID(t *testing.T) {
	data := GetThemeCountByCategoryID(newThemeCategory.ID)
	ToolsTest.ReportData(t, nil, data)
	if data < 1 {
		t.Error("get theme count failed, category id: ", newThemeCategory.ID, ", count: ", data)
	}
}

func TestUpdateTheme(t *testing.T) {
	err := UpdateTheme(&ArgsUpdateTheme{
		ID:          newThemeData.ID,
		Name:        "测试主题_Update",
		Description: "测试主题描述_Update",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteTheme(t *testing.T) {
	err := DeleteTheme(&ArgsDeleteTheme{
		ID: newThemeData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestThemeClear(t *testing.T) {
	TestDeleteThemeCategory(t)
	TestThemeCategoryClear(t)
}
