package OrgCoreCore

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newSystemData FieldsSystem
)

func TestInitSystem(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	TestCreateOrg(t)
}

func TestSetSystem(t *testing.T) {
	var err error
	newSystemData, err = SetSystem(&ArgsSetSystem{
		OrgID:      orgData.ID,
		SystemMark: "wxx",
		Mark:       "123123",
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newSystemData)
}

func TestGetSystemList(t *testing.T) {
	dataList, dataCount, err := GetSystemList(&ArgsGetSystemList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      -1,
		SystemMark: "",
		Mark:       "",
		IsRemove:   false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestFilterOrgIDsBySystem(t *testing.T) {
	data, err := FilterOrgIDsBySystem(&ArgsFilterOrgIDsBySystem{
		OrgIDs:     []int64{orgData.ID},
		SystemMark: newSystemData.SystemMark,
		Mark:       newSystemData.Mark,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteSystem(t *testing.T) {
	err := DeleteSystem(&ArgsDeleteSystem{
		ID:    newSystemData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearSystem(t *testing.T) {
	TestDeleteWorkTime(t)
	TestDeleteOrg(t)
	TestClear(t)
}
