package UserCore

import (
	"testing"
)

func TestCreatePermission(t *testing.T) {
	TestInit(t)
	if err := CreatePermission(&ArgsCreatePermission{
		Mark:     "test",
		Name:     "测试权限",
		Des:      "测试权限描述...",
		AllowOrg: false,
	}); err != nil {
		t.Error(err)
	}
}

func TestGetPermissionByMark(t *testing.T) {
	data, err := GetPermissionByMark(&ArgsGetPermissionByMark{
		Mark: "test",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetAllPermission(t *testing.T) {
	data, err := GetAllPermission(&ArgsGetAllPermission{
		AllowOrg: false,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestUpdatePermission(t *testing.T) {
	if err := UpdatePermission(&ArgsUpdatePermission{
		Mark:     "test",
		Name:     "测试权限2",
		Des:      "测试权限描述2。。。",
		AllowOrg: false,
	}); err != nil {
		t.Error(err)
	}
	TestGetPermissionByMark(t)
}

func TestDeletePermission(t *testing.T) {
	if err := DeletePermission(&ArgsDeletePermission{
		Mark: "test",
	}); err != nil {
		t.Error(err)
	}
}

func TestDeleteAllPermission(t *testing.T) {
	if err := DeleteAllPermission(); err != nil {
		t.Error(err)
	}
}
