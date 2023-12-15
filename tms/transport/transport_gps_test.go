package TMSTransport

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestTransportGPSInit(t *testing.T) {
	TestBindInit(t)
	TestSetBind(t)
	TestCreateTransport(t)
	TestUpdateTransportGPS(t)
}

func TestGetTransportGPSGroup(t *testing.T) {
	data, err := GetTransportGPSGroup(&ArgsGetTransportGPSGroup{
		OrgID:       newTransportData.OrgID,
		TransportID: newTransportData.ID,
		MinTime:     CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHours(2).Time),
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetTransportGPSLast(t *testing.T) {
	data, err := GetTransportGPSLast(&ArgsGetTransportGPSLast{
		OrgID:       newTransportData.OrgID,
		TransportID: newTransportData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetTransportGPSList(t *testing.T) {
	dataList, dataCount, err := GetTransportGPSList(&ArgsGetTransportGPSList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:       -1,
		TransportID: -1,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestTransportGPSClear(t *testing.T) {
	TestDeleteTransport(t)
	TestDeleteBind(t)
	TestBindClear(t)
}
