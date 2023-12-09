package ServiceUserInfoCost

import (
	"testing"

	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
)

func TestInitCost(t *testing.T) {
	TestInitConfig(t)
}

func TestSetCost(t *testing.T) {
	err := SetCost(&ArgsSetCost{
		CreateAt:     CoreFilter.GetISOByTime(CoreFilter.GetNowTime()),
		OrgID:        TestOrg.OrgData.ID,
		RoomID:       1,
		InfoID:       1,
		RoomBindMark: "room_ele",
		SensorMark:   "ele",
		Unit:         1,
		Currency:     86,
		Price:        310,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetCostList(t *testing.T) {
	dataList, dataCount, err := GetCostList(&ArgsGetCostList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:        0,
		RoomBindMark: "",
		SensorMark:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetCostLast(t *testing.T) {
	data, err := GetCostLast(&ArgsGetCostLast{
		OrgID:        TestOrg.OrgData.ID,
		RoomID:       1,
		InfoID:       1,
		RoomBindMark: "room_ele",
		SensorMark:   "ele",
	})
	ToolsTest.ReportData(t, err, data)
}
