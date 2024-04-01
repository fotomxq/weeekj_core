package BaseBPM

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newLogData FieldsLog
)

func TestLogInit(t *testing.T) {
	TestBPMInit(t)
	TestCreateBPM(t)
	TestGetBPMByID(t)
}

func TestCreateLog(t *testing.T) {
	newDataID, err := CreateLog(&ArgsCreateLog{
		OrgID:       1,
		UnitID:      1,
		UserID:      1,
		OrgBindID:   0,
		BPMID:       newBPMData.ID,
		NodeID:      "",
		NodeNumber:  0,
		NodeContent: "",
	})
	ToolsTest.ReportData(t, err, newDataID)
	newLogData.ID = newDataID
}

func TestGetLogByID(t *testing.T) {
	data, err := GetLogByID(&ArgsGetLogByID{
		ID: newLogData.ID,
	})
	ToolsTest.ReportData(t, err, data)
	newLogData = data
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:     -1,
		UnitID:    -1,
		UserID:    -1,
		OrgBindID: -1,
		BPMID:     newBPMData.ID,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestLogClear(t *testing.T) {
	TestDeleteBPM(t)
	TestBPMClear(t)
}
