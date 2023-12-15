package MarketGivingQrcode

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	MarketGivingCore "github.com/fotomxq/weeekj_core/v5/market/giving_core"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newConditionsData FieldsConditions
	newConfigData     MarketGivingCore.FieldsConfig
)

func TestInitConditions(t *testing.T) {
	TestInit(t)
}

func TestCreateConfig(t *testing.T) {
	var err error
	newConfigData, err = MarketGivingCore.CreateConfig(&MarketGivingCore.ArgsCreateConfig{
		OrgID:             TestOrg.OrgData.ID,
		Name:              "测试推荐处理",
		MarketConfigID:    0,
		LimitTimeType:     0,
		LimitCount:        0,
		UserIntegral:      0,
		UserSubs:          MarketGivingCore.FieldsConfigUserSubs{},
		UserTickets:       MarketGivingCore.FieldsConfigUserTickets{},
		DepositConfigMark: "",
		Price:             0,
		Count:             0,
		Params:            nil,
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestCreateConditions(t *testing.T) {
	var err error
	newConditionsData, err = CreateConditions(&ArgsCreateConditions{
		OrgID:     TestOrg.OrgData.ID,
		Name:      "测试条件",
		ConfigID:  newConfigData.ID,
		HavePhone: false,
		HaveOrder: false,
		Params:    nil,
	})
	ToolsTest.ReportData(t, err, newConditionsData)
}

func TestGetConditionsList(t *testing.T) {
	dataList, dataCount, err := GetConditionsList(&ArgsGetConditionsList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		ConfigID: -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateConditions(t *testing.T) {
	err := UpdateConditions(&ArgsUpdateConditions{
		ID:        newConditionsData.ID,
		OrgID:     newConditionsData.OrgID,
		Name:      newConditionsData.Name,
		ConfigID:  newConditionsData.ConfigID,
		HavePhone: newConditionsData.HavePhone,
		HaveOrder: newConditionsData.HaveOrder,
		Params:    newConditionsData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteConditions(t *testing.T) {
	err := DeleteConditions(&ArgsDeleteConditions{
		ID:    newConditionsData.ID,
		OrgID: newConditionsData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteConfig(t *testing.T) {
	err := MarketGivingCore.DeleteConfig(&MarketGivingCore.ArgsDeleteConfig{
		ID:    newConfigData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearConditions(t *testing.T) {
	TestClear(t)
}
