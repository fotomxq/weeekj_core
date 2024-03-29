package BaseBPM

import (
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

func TestGetThemeCategory(t *testing.T) {
	data, err := GetThemeCategoryByID(&ArgsGetThemeCategoryByID{
		ID: newThemeCategory.ID,
	})
	ToolsTest.ReportData(t, err, data)
	newThemeCategory = data
}

func TestUpdateThemeCategory(t *testing.T) {
	err := UpdateThemeCategory(&ArgsUpdateThemeCategory{
		ID:          newThemeCategory.ID,
		Name:        newThemeCategory.Name,
		Description: newThemeCategory.Description,
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
