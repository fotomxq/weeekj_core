package MallCore

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newTransportData FieldsTransport
)

func TestInitTransport(t *testing.T) {
	TestInit(t)
}

func TestCreateTransport(t *testing.T) {
	var err error
	newTransportData, err = CreateTransport(&ArgsCreateTransport{
		OrgID:      TestOrg.OrgData.ID,
		Name:       "测试配送模版",
		Rules:      0,
		RulesUnit:  0,
		RulesPrice: 0,
		AddType:    0,
		AddUnit:    0,
		AddPrice:   0,
		FreeType:   0,
		FreeUnit:   0,
	})
	ToolsTest.ReportData(t, err, newTransportData)
}

func TestGetTransportID(t *testing.T) {
	var err error
	newTransportData, err = GetTransportID(&ArgsGetTransportID{
		ID:    newTransportData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportData(t, err, newTransportData)
}

func TestGetTransportList(t *testing.T) {
	dataList, dataCount, err := GetTransportList(&ArgsGetTransportList{
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

func TestGetTransports(t *testing.T) {
	data, err := GetTransports(&ArgsGetTransports{
		IDs:        []int64{newTransportData.ID},
		HaveRemove: false,
		OrgID:      TestOrg.OrgData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateTransport(t *testing.T) {
	err := UpdateTransport(&ArgsUpdateTransport{
		ID:         newTransportData.ID,
		OrgID:      newTransportData.OrgID,
		Name:       newTransportData.Name,
		Rules:      newTransportData.Rules,
		RulesUnit:  newTransportData.RulesUnit,
		RulesPrice: newTransportData.RulesPrice,
		AddType:    newTransportData.AddType,
		AddUnit:    newTransportData.AddUnit,
		AddPrice:   newTransportData.AddPrice,
		FreeType:   newTransportData.FreeType,
		FreeUnit:   newTransportData.FreeUnit,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportFree(t *testing.T) {
	err := UpdateTransport(&ArgsUpdateTransport{
		ID:         newTransportData.ID,
		OrgID:      newTransportData.OrgID,
		Name:       newTransportData.Name,
		Rules:      0,
		RulesUnit:  0,
		RulesPrice: 0,
		AddType:    0,
		AddUnit:    0,
		AddPrice:   0,
		FreeType:   0,
		FreeUnit:   0,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportUnit1(t *testing.T) {
	err := UpdateTransport(&ArgsUpdateTransport{
		ID:         newTransportData.ID,
		OrgID:      newTransportData.OrgID,
		Name:       newTransportData.Name,
		Rules:      1,
		RulesUnit:  1,
		RulesPrice: 10,
		AddType:    1,
		AddUnit:    1,
		AddPrice:   5,
		FreeType:   1,
		FreeUnit:   30,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportUnit2(t *testing.T) {
	err := UpdateTransport(&ArgsUpdateTransport{
		ID:         newTransportData.ID,
		OrgID:      newTransportData.OrgID,
		Name:       newTransportData.Name,
		Rules:      2,
		RulesUnit:  1,
		RulesPrice: 10,
		AddType:    2,
		AddUnit:    1,
		AddPrice:   5,
		FreeType:   2,
		FreeUnit:   30,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportUnit3(t *testing.T) {
	err := UpdateTransport(&ArgsUpdateTransport{
		ID:         newTransportData.ID,
		OrgID:      newTransportData.OrgID,
		Name:       newTransportData.Name,
		Rules:      3,
		RulesUnit:  1,
		RulesPrice: 10,
		AddType:    3,
		AddUnit:    1,
		AddPrice:   5,
		FreeType:   3,
		FreeUnit:   30,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteTransport(t *testing.T) {
	err := DeleteTransport(&ArgsDeleteTransport{
		ID:    newTransportData.ID,
		OrgID: newTransportData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearTransport(t *testing.T) {
	TestClear(t)
}
