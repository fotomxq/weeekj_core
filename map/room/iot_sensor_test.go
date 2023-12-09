package MapRoom

import (
	"testing"
	"time"

	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
)

func TestInitIOTSensor(t *testing.T) {
	TestInit(t)
}

func TestGetSensorList(t *testing.T) {
	dataList, dataCount, err := GetSensorList(&ArgsGetSensorList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:     -1,
		RoomID:    -1,
		DeviceID:  -1,
		Mark:      "",
		IsHistory: false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetSensorAnalysis(t *testing.T) {
	dataList, err := GetSensorAnalysis(&ArgsGetSensorAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubDays(1).Time,
			MaxTime: CoreFilter.GetNowTime().Add(time.Second * 1),
		},
		TimeType:  "hour",
		OrgID:     -1,
		RoomID:    -1,
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetSensorAnalysisAvg(t *testing.T) {
	dataList, err := GetSensorAnalysisAvg(&ArgsGetSensorAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubDays(1).Time,
			MaxTime: CoreFilter.GetNowTime().Add(time.Second * 1),
		},
		TimeType:  "hour",
		OrgID:     -1,
		RoomID:    -1,
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetSensorAnalysisMax(t *testing.T) {
	dataList, err := GetSensorAnalysisMax(&ArgsGetSensorAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubDays(1).Time,
			MaxTime: CoreFilter.GetNowTime().Add(time.Second * 1),
		},
		TimeType:  "hour",
		OrgID:     -1,
		RoomID:    -1,
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetSensorAnalysisMin(t *testing.T) {
	dataList, err := GetSensorAnalysisMin(&ArgsGetSensorAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubDays(1).Time,
			MaxTime: CoreFilter.GetNowTime().Add(time.Second * 1),
		},
		TimeType:  "hour",
		OrgID:     -1,
		RoomID:    -1,
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestDeleteSensorClear(t *testing.T) {
	err := DeleteSensorClear(&ArgsDeleteSensorClear{
		OrgID:     -1,
		RoomID:    -1,
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportError(t, err)
}
