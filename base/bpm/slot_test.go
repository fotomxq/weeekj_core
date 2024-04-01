package BaseBPM

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newSlotData FieldsSlot
)

func TestSlotInit(t *testing.T) {
	TestThemeInit(t)
	TestCreateTheme(t)
	TestGetThemeByID(t)
}

func TestCreateSlot(t *testing.T) {
	newDataID, err := CreateSlot(&ArgsCreateSlot{
		Name:            "测试槽位",
		ThemeCategoryID: newThemeCategory.ID,
		ThemeID:         newThemeData.ID,
		ValueType:       "input",
		DefaultValue:    "default_test",
		Params:          "params_test",
	})
	ToolsTest.ReportData(t, err, newDataID)
	newSlotData.ID = newDataID
}

func TestGetSlotByID(t *testing.T) {
	data, err := GetSlotByID(&ArgsGetSlotByID{
		ID: newSlotData.ID,
	})
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		t.Error("get slot failed, id: ", newSlotData.ID, ", err: ", err)
		t.Fail()
		return
	}
	newSlotData = data
}

func TestGetSlotCountByCategoryID(t *testing.T) {
	data := GetSlotCountByCategoryID(newThemeCategory.ID)
	ToolsTest.ReportData(t, nil, data)
	if data < 1 {
		t.Error("get slot count failed, category id: ", newThemeCategory.ID)
	}
}

func TestGetSlotCountByThemeID(t *testing.T) {
	data := GetSlotCountByThemeID(newThemeData.ID)
	ToolsTest.ReportData(t, nil, data)
	if data < 1 {
		t.Error("get slot count failed, theme id: ", newThemeData.ID)
	}
}

func TestGetSlotList(t *testing.T) {
	dataList, dataCount, err := GetSlotList(&ArgsGetSlotList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		ThemeCategoryID: newThemeCategory.ID,
		ThemeID:         newThemeData.ID,
		IsRemove:        false,
		Search:          "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateSlot(t *testing.T) {
	err := UpdateSlot(&ArgsUpdateSlot{
		ID:           newSlotData.ID,
		Name:         "测试槽位_Update",
		ValueType:    "input",
		DefaultValue: "default_test_Update",
		Params:       "params_test_Update",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteSlot(t *testing.T) {
	err := DeleteSlot(&ArgsDeleteSlot{
		ID: newSlotData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestSlotClear(t *testing.T) {
	TestDeleteTheme(t)
	TestThemeClear(t)
}
