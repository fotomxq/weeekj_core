package OrgCoreCore

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitPermissionFunc(t *testing.T) {
	TestInit(t)
}

func TestSetPermissionFunc(t *testing.T) {
	err := SetPermissionFunc(&ArgsSetPermissionFunc{
		Mark:        "test",
		Name:        "测试",
		Des:         "测试描述",
		ParentMarks: []string{"test1"},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetAllPermissionFunc(t *testing.T) {
	data, err := GetAllPermissionFunc()
	ToolsTest.ReportData(t, err, data)
}

func TestDeletePermissionFunc(t *testing.T) {
	err := DeletePermissionFunc(&ArgsDeletePermissionFunc{
		Mark: "test",
	})
	ToolsTest.ReportError(t, err)
}

func TestClearPermissionFunc(t *testing.T) {

}
