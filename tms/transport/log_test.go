package TMSTransport

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

func TestLogInit(t *testing.T) {
	TestTransportInit(t)
	TestCreateTransport(t)
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:           -1,
		BindID:          0,
		TransportID:     0,
		TransportBindID: 0,
		Mark:            "",
		Search:          "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestLogClear(t *testing.T) {
	TestDeleteTransport(t)
	TestTransportClear(t)
}
