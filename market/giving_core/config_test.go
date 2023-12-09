package MarketGivingCore

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newConfigData FieldsConfig
)

func TestInitConfig(t *testing.T) {
	TestInit(t)
}

func TestCreateConfig(t *testing.T) {
	var err error
	newConfigData, err = CreateConfig(&ArgsCreateConfig{
		OrgID:             TestOrg.OrgData.ID,
		Name:              "测试推荐处理",
		MarketConfigID:    0,
		LimitTimeType:     0,
		LimitCount:        0,
		UserIntegral:      0,
		UserSubs:          FieldsConfigUserSubs{},
		UserTickets:       FieldsConfigUserTickets{},
		DepositConfigMark: "",
		Price:             0,
		Count:             0,
		Params:            nil,
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestGetConfigByID(t *testing.T) {
	var err error
	newConfigData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    newConfigData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestGetConfigList(t *testing.T) {
	dataList, dataCount, err := GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    TestOrg.OrgData.ID,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetConfigMore(t *testing.T) {
	data, err := GetConfigMore(&ArgsGetConfigMore{
		IDs:        []int64{newConfigData.ID},
		HaveRemove: false,
		OrgID:      TestOrg.OrgData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetConfigMoreMap(t *testing.T) {
	data, err := GetConfigMoreMap(&ArgsGetConfigMore{
		IDs:        []int64{newConfigData.ID},
		HaveRemove: false,
		OrgID:      TestOrg.OrgData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateConfig(t *testing.T) {
	err := UpdateConfig(&ArgsUpdateConfig{
		ID:                newConfigData.ID,
		OrgID:             TestOrg.OrgData.ID,
		Name:              newConfigData.Name,
		MarketConfigID:    newConfigData.MarketConfigID,
		LimitTimeType:     newConfigData.LimitTimeType,
		LimitCount:        newConfigData.LimitCount,
		UserIntegral:      newConfigData.UserIntegral,
		UserSubs:          newConfigData.UserSubs,
		UserTickets:       newConfigData.UserTickets,
		DepositConfigMark: newConfigData.DepositConfigMark,
		Price:             newConfigData.Price,
		Count:             newConfigData.Count,
		Params:            newConfigData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	err := DeleteConfig(&ArgsDeleteConfig{
		ID:    newConfigData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearConfig(t *testing.T) {
	TestClear(t)
}
