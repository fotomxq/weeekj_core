package IOTSensor

import (
	"testing"
	"time"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

func TestInitSensor(t *testing.T) {
	TestInit(t)
}

func TestCreate(t *testing.T) {
	err := Create(&ArgsCreate{
		CreateAt: "",
		DeviceID: 1,
		Mark:     "test1",
		Data:     1,
		DataF:    1.1,
		DataS:    "ttt",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		DeviceID:  -1,
		Mark:      "",
		IsHistory: false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetAnalysis(t *testing.T) {
	dataList, err := GetAnalysis(&ArgsGetAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubDays(1).Time,
			MaxTime: CoreFilter.GetNowTime().Add(time.Second * 1),
		},
		TimeType:  "hour",
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetAnalysisAvg(t *testing.T) {
	dataList, err := GetAnalysisAvg(&ArgsGetAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubDays(1).Time,
			MaxTime: CoreFilter.GetNowTime().Add(time.Second * 1),
		},
		TimeType:  "hour",
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetAnalysisMax(t *testing.T) {
	dataList, err := GetAnalysisMax(&ArgsGetAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubDays(1).Time,
			MaxTime: CoreFilter.GetNowTime().Add(time.Second * 1),
		},
		TimeType:  "hour",
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetAnalysisMin(t *testing.T) {
	dataList, err := GetAnalysisMin(&ArgsGetAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubDays(1).Time,
			MaxTime: CoreFilter.GetNowTime().Add(time.Second * 1),
		},
		TimeType:  "hour",
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestDeleteClear(t *testing.T) {
	err := DeleteClear(&ArgsDeleteClear{
		DeviceID:  1,
		Mark:      "test1",
		IsHistory: false,
	})
	ToolsTest.ReportError(t, err)
}
