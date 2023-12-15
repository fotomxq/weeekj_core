package BaseSafe

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitLog(t *testing.T) {
	TestInit(t)
}

func TestCreateLog(t *testing.T) {
	CreateLog(&ArgsCreateLog{
		System: "test",
		Level:  1,
		IP:     "0.0.0.1",
		UserID: 123,
		OrgID:  23,
		Des:    "des",
	})
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:  -1,
		UserID: -1,
		System: "",
		Level:  -1,
		IP:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}
