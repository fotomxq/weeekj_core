package OrgCoreCore

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

// 注意，本测试必定经过了组织创建
func TestInit5(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	TestCreateOrg(t)
	TestCreateGroup(t)
	TestSetBind(t)
}

func TestGetPermissionsByOrg(t *testing.T) {
	data := GetPermissionsByOrg(orgData.ID)
	ToolsTest.ReportData(t, nil, data)
}

func TestCheckPermissionsByBindOrGroup(t *testing.T) {
	reply, err := CheckPermissionsByBind(&ArgsCheckPermissionsByBind{
		BindID:      bindData.ID,
		Permissions: []string{"org_view"},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(reply)
	}
}

func TestCheckPermissionsByBindOrGroupOnlyBool(t *testing.T) {
	b := CheckPermissionsByBindOrGroupOnlyBool(&ArgsCheckPermissionsByBind{
		BindID:      bindData.ID,
		Permissions: []string{"org_view"},
	})
	if !b {
		t.Error(b)
	}
}

func TestClear6(t *testing.T) {
	TestDeleteBind(t)
	TestDeleteGroup(t)
	TestDeleteWorkTime(t)
	TestDeleteOrg(t)
	TestClear(t)
}
