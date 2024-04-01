package BaseBPM

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newEventData FieldsEvent
)

func TestEventInit(t *testing.T) {
	TestSlotInit(t)
}

func TestCreateEvent(t *testing.T) {
	newDataID, err := CreateEvent(&ArgsCreateEvent{
		Name:            "测试事件",
		Description:     "测试事件描述",
		ThemeCategoryID: newThemeCategory.ID,
		ThemeID:         newThemeData.ID,
		Code:            "A001",
		EventType:       "nats",
		EventURL:        "/test/01",
		EventParams:     "test_params",
	})
	ToolsTest.ReportData(t, err, newDataID)
	newEventData.ID = newDataID
}

func TestGetEventByID(t *testing.T) {
	data, err := GetEventByID(&ArgsGetEventByID{
		ID: newEventData.ID,
	})
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		t.Error("get event failed, id: ", newEventData.ID)
	}
	newEventData = data
}

func TestGetEventCountByCategoryID(t *testing.T) {
	data := GetEventCountByCategoryID(newThemeCategory.ID)
	ToolsTest.ReportData(t, nil, data)
	if data < 1 {
		t.Error("get event count failed, category id: ", newThemeCategory.ID)
	}
}

func TestGetEventCountByThemeID(t *testing.T) {
	data := GetEventCountByThemeID(newThemeData.ID)
	ToolsTest.ReportData(t, nil, data)
	if data < 1 {
		t.Error("get event count failed, theme id: ", newThemeData.ID)
	}
}

func TestGetEventList(t *testing.T) {
	dataList, dataCount, err := GetEventList(&ArgsGetEventList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		ThemeCategoryID: newThemeCategory.ID,
		ThemeID:         newThemeData.ID,
		Code:            "",
		IsRemove:        false,
		Search:          "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err != nil {
		t.Error("get event list failed, theme id: ", newThemeData.ID)
	}
}

func TestUpdateEvent(t *testing.T) {
	err := UpdateEvent(&ArgsUpdateEvent{
		ID:              newEventData.ID,
		Name:            "测试事件_Update",
		Description:     "测试事件描述_Update",
		ThemeCategoryID: newThemeCategory.ID,
		ThemeID:         newThemeData.ID,
		Code:            "A002",
		EventType:       "nats",
		EventURL:        "/test/02",
		EventParams:     "test_params_update",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteEvent(t *testing.T) {
	err := DeleteEvent(&ArgsDeleteEvent{
		ID: newEventData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestEventClear(t *testing.T) {
	TestSlotClear(t)
}
