package ToolsTestUserRole

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	UserRole "gitee.com/weeekj/weeekj_core/v5/user/role"
	"testing"
)

//用户角色测试工具组合

var (
	//RoleData 用户角色数据包
	RoleData UserRole.FieldsRole
	//RoleType 用户角色类型
	RoleType UserRole.FieldsType
	//RoleMark Mark
	RoleMark = "test"
)

func CreateType(t *testing.T) {
	var err error
	RoleType, err = UserRole.CreateType(&UserRole.ArgsCreateType{
		Mark:     RoleMark,
		Name:     "测试分组",
		GroupIDs: []int64{},
		Params:   CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, RoleType)
}

func CreateUserRole(t *testing.T) {
	if RoleType.ID < 1 {
		CreateType(t)
	}
	var err error
	RoleData, err = UserRole.SetRole(&UserRole.ArgsSetRole{
		RoleType:    RoleType.ID,
		ApplyID:     0,
		UserID:      TestOrg.UserInfo.ID,
		Name:        TestOrg.UserInfo.Name,
		Country:     86,
		City:        "10010",
		Gender:      1,
		Phone:       TestOrg.UserInfo.Phone,
		CoverFileID: TestOrg.UserInfo.Avatar,
		CertFiles:   []int64{},
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, RoleData)
}

func Clear(t *testing.T) {
	if RoleData.ID > 0 {
		err := UserRole.DeleteRole(&UserRole.ArgsDeleteRole{
			ID: RoleData.ID,
		})
		ToolsTest.ReportError(t, err)
	}
	if RoleType.ID > 0 {
		err := UserRole.DeleteType(&UserRole.ArgsDeleteType{
			ID: RoleType.ID,
		})
		ToolsTest.ReportError(t, err)
	}
}
