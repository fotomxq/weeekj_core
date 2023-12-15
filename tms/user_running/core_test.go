package TMSUserRunning

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	ToolsTestUserRole "github.com/fotomxq/weeekj_core/v5/tools/test_user_role"
	"testing"
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
	//初始化商户绑定关系
	TestOrg.LocalCreateBind(t)
	//创建用户角色
	ToolsTestUserRole.RoleMark = "tms_running"
	ToolsTestUserRole.CreateUserRole(t)
}

func TestClear(t *testing.T) {
	ToolsTestUserRole.Clear(t)
	TestOrg.LocalClear(t)
}
