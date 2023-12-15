package UserRole

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newRoleData FieldsRole
)

func TestInitRole(t *testing.T) {
	TestInitType(t)
	TestCreateType(t)
}

func TestSetRole(t *testing.T) {
	data, err := SetRole(&ArgsSetRole{
		RoleType:    newTypeData.ID,
		ApplyID:     0,
		UserID:      1,
		Name:        "测试角色",
		Country:     86,
		City:        "10010",
		Gender:      1,
		Phone:       "17777777777",
		CoverFileID: 0,
		CertFiles:   []int64{},
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newRoleData = data
	}
}

func TestGetRoleList(t *testing.T) {
	dataList, dataCount, err := GetRoleList(&ArgsGetRoleList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		RoleType: -1,
		UserID:   -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetRoleID(t *testing.T) {
	data, err := GetRoleID(&ArgsGetRoleID{
		ID: newRoleData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteRole(t *testing.T) {
	err := DeleteRole(&ArgsDeleteRole{
		ID: newRoleData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearRole(t *testing.T) {
	TestDeleteType(t)
}
