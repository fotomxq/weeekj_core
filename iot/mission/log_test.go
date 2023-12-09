package IOTMission

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitLog(t *testing.T) {
	TestInit(t)
}

func TestCreateLog(t *testing.T) {
	err := CreateLog(&ArgsCreateLog{
		MissionID: 1,
		Status:    1,
		Mark:      "test",
		Content:   "test des",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		MissionID: -1,
		Status:    -1,
		Mark:      "",
		Search:    "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}
