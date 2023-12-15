package ServiceHousekeeping

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newBindData FieldsBind
)

func TestInitBind(t *testing.T) {
	TestInit(t)
}

func TestSetBind(t *testing.T) {
	var err error
	newBindData, err = SetBind(&ArgsSetBind{
		OrgID:     TestOrg.OrgData.ID,
		BindID:    TestOrg.BindData.ID,
		MapAreaID: 0,
		Params:    nil,
	})
	ToolsTest.ReportData(t, err, newBindData)
}

func TestGetBindList(t *testing.T) {
	data, dataCount, err := GetBindList(&ArgsGetBindList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
	})
	ToolsTest.ReportDataList(t, err, data, dataCount)
}

func TestGetBindByBind(t *testing.T) {
	data, err := GetBindByBind(&ArgsGetBindByBind{
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetBindID(t *testing.T) {
	data, err := GetBindID(&ArgsGetBindID{
		ID:     newBindData.ID,
		OrgID:  TestOrg.OrgData.ID,
		BindID: TestOrg.BindData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteBind(t *testing.T) {
	err := DeleteBind(&ArgsDeleteBind{
		ID:    newBindData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearBind(t *testing.T) {
	TestClear(t)
}
