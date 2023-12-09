package OrgMission

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
		OrgID:       1,
		BindID:      0,
		MissionID:   0,
		ContentMark: "create",
		Content:     "内容测试",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:       -1,
		BindID:      -1,
		MissionID:   -1,
		ContentMark: "",
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}
