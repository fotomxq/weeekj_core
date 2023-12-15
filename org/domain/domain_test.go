package OrgDomain

import (
	"testing"

	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
)

var (
	newDomainData FieldsDomain
)

func TestInitDomain(t *testing.T) {
	TestInit(t)
}

func TestCreateDomain(t *testing.T) {
	data, err := CreateDomain(&ArgsCreateDomain{
		OrgID:  TestOrg.OrgData.ID,
		Host:   "localhost:29003",
		Params: nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newDomainData = data
	}
}

func TestGetDomainOrg(t *testing.T) {
	data, _, err := GetDomainOrg(&ArgsGetDomainOrg{
		Host: newDomainData.Host,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetDomainList(t *testing.T) {
	dataList, dataCount, err := GetDomainList(&ArgsGetDomainList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:  newDomainData.OrgID,
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateDomain(t *testing.T) {
	err := UpdateDomain(&ArgsUpdateDomain{
		ID:     newDomainData.ID,
		OrgID:  newDomainData.OrgID,
		Host:   newDomainData.Host,
		Params: newDomainData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteDomain(t *testing.T) {
	err := DeleteDomain(&ArgsDeleteDomain{
		ID:    newDomainData.ID,
		OrgID: newDomainData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearDomain(t *testing.T) {
	TestClear(t)
}
