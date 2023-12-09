package IOTLog

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit   = false
	newLogID int64
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestAppend(t *testing.T) {
	Append(&ArgsAppend{
		OrgID:    1,
		GroupID:  1,
		DeviceID: 1,
		Mark:     "test",
		Content:  "test_content",
	})
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:       0,
		GroupID:     0,
		DeviceID:    0,
		TimeBetween: CoreSQLTime.FieldsCoreTime{},
		IsHistory:   false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if len(dataList) > 0 {
		newLogID = dataList[0].ID
	}
}

func TestGetByID(t *testing.T) {
	data, err := GetByID(&ArgsGetByID{
		ID:    newLogID,
		OrgID: -1,
	})
	ToolsTest.ReportData(t, err, data)
}
