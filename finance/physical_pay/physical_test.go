package FinancePhysicalPay

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newPhysicalData FieldsPhysical
)

func TestInitPhysical(t *testing.T) {
	TestInit(t)
}

func TestCreatePhysical(t *testing.T) {
	var err error
	newPhysicalData, err = CreatePhysical(&ArgsCreatePhysical{
		OrgID: TestOrg.OrgData.ID,
		Name:  "测试抵扣物品",
		BindFrom: CoreSQLFrom.FieldsFrom{
			System: "mall",
			ID:     1,
			Mark:   "",
			Name:   "",
		},
		NeedCount:  2,
		LimitCount: 100,
		Params:     nil,
	})
	ToolsTest.ReportData(t, err, newPhysicalData)
}

func TestGetPhysicalByFrom(t *testing.T) {
	data, err := GetPhysicalByFrom(&ArgsGetPhysicalByFrom{
		OrgID: TestOrg.OrgData.ID,
		BindFrom: CoreSQLFrom.FieldsFrom{
			System: "mall",
			ID:     1,
			Mark:   "",
			Name:   "",
		},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetPhysicalID(t *testing.T) {
	data, err := GetPhysicalID(&ArgsGetPhysicalID{
		ID:    newPhysicalData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetPhysicalList(t *testing.T) {
	dataList, dataCount, err := GetPhysicalList(&ArgsGetPhysicalList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdatePhysical(t *testing.T) {
	err := UpdatePhysical(&ArgsUpdatePhysical{
		ID:         newPhysicalData.ID,
		OrgID:      -1,
		Name:       newPhysicalData.Name,
		NeedCount:  newPhysicalData.NeedCount,
		LimitCount: newPhysicalData.LimitCount,
		Params:     nil,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeletePhysical(t *testing.T) {
	err := DeletePhysical(&ArgsDeletePhysical{
		ID:    newPhysicalData.ID,
		OrgID: newPhysicalData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearPhysical(t *testing.T) {
	TestClear(t)
}
