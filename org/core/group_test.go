package OrgCoreCore

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"testing"
)

var (
	orgGroupData FieldsGroup
)

func TestInit4(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	TestCreateOrg(t)
}

func TestCreateGroup(t *testing.T) {
	var err error
	orgGroupData, err = CreateGroup(&ArgsCreateGroup{
		OrgID:   orgData.ID,
		Name:    "测试分组",
		Manager: []string{"all"},
		Params:  []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(orgGroupData)
		TestCreateWorkTime(t)
	}
}

func TestGetGroup(t *testing.T) {
	data, err := GetGroup(&ArgsGetGroup{
		ID:    orgGroupData.ID,
		OrgID: orgData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestUpdateGroup(t *testing.T) {
	err := UpdateGroup(&ArgsUpdateGroup{
		ID:      orgGroupData.ID,
		OrgID:   orgGroupData.OrgID,
		Name:    orgGroupData.Name + "edit",
		Manager: []string{"all", "member"},
		Params:  nil,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteGroup(t *testing.T) {
	err := DeleteGroup(&ArgsDeleteGroup{
		ID:    orgGroupData.ID,
		OrgID: orgGroupData.OrgID,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClear2(t *testing.T) {
	TestDeleteWorkTime(t)
	TestDeleteOrg(t)
	TestClear(t)
}
