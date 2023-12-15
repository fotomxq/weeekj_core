package IOTTrack

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	deviceID int64 = 234
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestCreate(t *testing.T) {
	err := Create(&ArgsCreate{
		DeviceID:    deviceID,
		MapType:     0,
		Longitude:   1.2,
		Latitude:    1.3,
		StationInfo: "test station info",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		DeviceID: 0,
		MapType:  0,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetLast(t *testing.T) {
	data, err := GetLast(&ArgsGetLast{
		DeviceID: deviceID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteByDevice(t *testing.T) {
	err := DeleteByDevice(&ArgsDeleteByDevice{
		DeviceID: deviceID,
	})
	ToolsTest.ReportError(t, err)
}
