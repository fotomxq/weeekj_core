package BaseBPM

import (
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newThemeCategory FieldsThemeCategory
)

func TestThemeCategoryInit(t *testing.T) {
	TestInit(t)
}

func TestCreateThemeCategory(t *testing.T) {
	newDataID, err := CreateThemeCategory(&ArgsCreateThemeCategory{
		Name:        "测试主题分类",
		Description: "测试主题分类描述",
	})
	ToolsTest.ReportData(t, err, newDataID)
	newThemeCategory.ID = newDataID
}

func TestGetThemeCategoryByID(t *testing.T) {
	data, err := GetThemeCategoryByID(&ArgsGetThemeCategoryByID{
		ID: newThemeCategory.ID,
	})
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		t.Error("get theme category failed, id: ", newThemeCategory.ID)
	}
	newThemeCategory = data
}

func TestGetThemeCategoryList(t *testing.T) {
	dataList, dataCount, err := GetThemeCategoryList(&ArgsGetThemeCategoryList{
		Pages: CoreSQL2.ArgsPages{
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

func TestUpdateThemeCategory(t *testing.T) {
	err := UpdateThemeCategory(&ArgsUpdateThemeCategory{
		ID:          newThemeCategory.ID,
		Name:        fmt.Sprint(newThemeCategory.Name, "_Update"),
		Description: fmt.Sprint(newThemeCategory.Description, "_Update"),
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteThemeCategory(t *testing.T) {
	err := DeleteThemeCategory(&ArgsDeleteThemeCategory{
		ID: newThemeCategory.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestThemeCategoryClear(t *testing.T) {
	TestClear(t)
}
