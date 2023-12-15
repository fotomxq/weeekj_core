package ServiceUserInfoCost

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	testConfigData FieldsConfig
)

func TestInitConfig(t *testing.T) {
	TestInit(t)
	TestOrg.LocalCreateBind(t)
}

func TestCreateConfig(t *testing.T) {
	data, err := CreateConfig(&ArgsCreateConfig{
		OrgID:        TestOrg.OrgData.ID,
		Name:         "测试标题",
		RoomBindMark: "room_ele",
		SensorMark:   "ele",
		CountType:    0,
		EachUnit:     1,
		Currency:     86,
		EachPrice:    100,
		Params:       nil,
	})
	if err == nil {
		testConfigData = data
	}
	ToolsTest.ReportData(t, err, data)
}

func TestGetConfigList(t *testing.T) {
	dataList, dataCount, err := GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:        0,
		RoomBindMark: "",
		SensorMark:   "",
		IsRemove:     false,
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetConfig(t *testing.T) {
	data, err := GetConfig(&ArgsGetConfig{
		ID:    testConfigData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	if err == nil {
		testConfigData = data
		t.Log("testConfigData id: ", testConfigData.ID, ", org id: ", testConfigData.OrgID)
	}
	ToolsTest.ReportData(t, err, data)
}

func TestGetConfigs(t *testing.T) {
	data, err := GetConfigs(&ArgsGetConfigs{
		IDs:        []int64{testConfigData.ID},
		HaveRemove: false,
		OrgID:      testConfigData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetConfigsName(t *testing.T) {
	data, err := GetConfigsName(&ArgsGetConfigs{
		IDs:        []int64{testConfigData.ID},
		HaveRemove: false,
		OrgID:      testConfigData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateConfig(t *testing.T) {
	err := UpdateConfig(&ArgsUpdateConfig{
		ID:           testConfigData.ID,
		OrgID:        testConfigData.OrgID,
		Name:         testConfigData.Name,
		RoomBindMark: testConfigData.RoomBindMark,
		SensorMark:   testConfigData.SensorMark,
		CountType:    testConfigData.CountType,
		EachUnit:     testConfigData.EachUnit,
		Currency:     testConfigData.Currency,
		EachPrice:    testConfigData.EachPrice,
		Params:       testConfigData.Params,
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Log("testConfigData id: ", testConfigData.ID, ", org id: ", testConfigData.OrgID, ", delete at: ", testConfigData.DeleteAt)
	}
}

func TestDeleteConfig(t *testing.T) {
	err := DeleteConfig(&ArgsDeleteConfig{
		ID:    testConfigData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportError(t, err)
}
