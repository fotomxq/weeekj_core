package UserCore

import (
	"testing"
)

var(
	groupData FieldsGroupType
)

func TestInit2(t *testing.T) {
	TestInit(t)
	TestCreatePermission(t)
}

func TestCreateGroup(t *testing.T) {
	var err error
	groupData, err = CreateGroup(&ArgsCreateGroup{
		OrgID:       0,
		Name:        "测试用户组",
		Des:         "测试用户组描述...",
		Permissions: []string{"test"},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetAllGroup(t *testing.T) {
	data, err := GetAllGroup(&ArgsGetAllGroup{
		OrgID: 0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetGroupByMark(t *testing.T) {
	data, err := GetGroup(&ArgsGetGroup{
		ID: groupData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestUpdateGroup(t *testing.T) {
	if err := UpdateGroup(&ArgsUpdateGroup{
		ID:          groupData.ID,
		OrgID:       0,
		Name:        "测试用户组2",
		Des:         "test_2",
		Permissions: []string{"test"},
	}); err != nil {
		t.Error(err)
	}
}

func TestDeleteGroup(t *testing.T) {
	if err := DeleteGroup(&ArgsDeleteGroup{
		ID:    groupData.ID,
		OrgID: 0,
	}); err != nil {
		t.Error(err)
	}
}

func TestDeletePermission2(t *testing.T) {
	TestDeletePermission(t)
}
