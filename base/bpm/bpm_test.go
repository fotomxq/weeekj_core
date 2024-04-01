package BaseBPM

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newBPMData FieldsBPM
)

func TestBPMInit(t *testing.T) {
	TestEventInit(t)
	TestCreateEvent(t)
	TestGetEventByID(t)
	TestCreateSlot(t)
	TestGetSlotByID(t)
}

func TestCreateBPM(t *testing.T) {
	newDataID, err := CreateBPM(&ArgsCreateBPM{
		Name:        "测试BPM",
		Description: "测试BPM描述",
		ThemeID:     newThemeData.ID,
		NodeCount:   0,
		JSONNode:    "",
	})
	ToolsTest.ReportData(t, err, newDataID)
	newBPMData.ID = newDataID
}

func TestGetBPMByID(t *testing.T) {
	data, err := GetBPMByID(&ArgsGetBPMByID{
		ID: newBPMData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetBPMCountByThemeID(t *testing.T) {
	data := GetBPMCountByThemeID(newThemeData.ID)
	ToolsTest.ReportData(t, nil, data)
	if data < 1 {
		t.Fatal("GetBPMCountByThemeID error.")
	}
}

func TestGetBPMList(t *testing.T) {
	dataList, dataCount, err := GetBPMList(&ArgsGetBPMList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		ThemeID:  newThemeData.ID,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateBPM(t *testing.T) {
	err := UpdateBPM(&ArgsUpdateBPM{
		ID:          newBPMData.ID,
		Name:        "测试BPM-修改",
		Description: "测试BPM描述-修改",
		ThemeID:     newThemeData.ID,
		NodeCount:   0,
		JSONNode:    "",
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteBPM(t *testing.T) {
	err := DeleteBPM(&ArgsDeleteBPM{
		ID: newBPMData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestBPMClear(t *testing.T) {
	TestDeleteEvent(t)
	TestDeleteSlot(t)
	TestEventClear(t)
}
