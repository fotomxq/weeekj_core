package TMSTransport

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestBindGPSInit(t *testing.T) {
	TestBindInit(t)
	TestSetBind(t)
	TestCreateTransport(t)
	TestUpdateTransportGPS(t)
}

func TestGetBindGPSGroup(t *testing.T) {
	data, err := GetBindGPSGroup(&ArgsGetBindGPSGroup{
		OrgID:   newBindData.OrgID,
		BindID:  newBindData.BindID,
		MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubHours(2).Time),
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetBindGPSLast(t *testing.T) {
	data, err := GetBindGPSLast(&ArgsGetBindGPSLast{
		OrgID:  newBindData.OrgID,
		BindID: newBindData.BindID,
	})
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		t.Log("org id: ", newBindData.OrgID, ", bind id: ", newBindData.BindID)
	}
}

func TestGetBindGPSList(t *testing.T) {
	dataList, dataCount, err := GetBindGPSList(&ArgsGetBindGPSList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:  -1,
		BindID: -1,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestBindGPSClear(t *testing.T) {
	TestDeleteTransport(t)
	TestDeleteBind(t)
	TestBindClear(t)
}
