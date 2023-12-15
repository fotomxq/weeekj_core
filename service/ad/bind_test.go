package ServiceAD

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	TestOrgArea "github.com/fotomxq/weeekj_core/v5/tools/test_org_area"
	"testing"
)

var (
	bindData FieldsBind
)

func TestInitBind(t *testing.T) {
	TestInit(t)
}

func TestSetBind(t *testing.T) {
	var err error
	bindData, err = SetBind(&ArgsSetBind{
		StartAt: CoreFilter.GetNowTimeCarbon().Time,
		EndAt:   CoreFilter.GetNowTimeCarbon().AddMonth().Time,
		OrgID:   TestOrg.OrgData.ID,
		AreaID:  TestOrgArea.AreaData.ID,
		AdID:    adData.ID,
		Factor:  1,
		Params:  []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, bindData)
}

func TestGetBindList(t *testing.T) {
	dataList, dataCount, err := GetBindList(&ArgsGetBindList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		AdID:     -1,
		AreaID:   -1,
		IsRemove: false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDeleteBind(t *testing.T) {
	err := DeleteBind(&ArgsDeleteBind{
		ID:    bindData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearBind(t *testing.T) {
	TestClear(t)
}
