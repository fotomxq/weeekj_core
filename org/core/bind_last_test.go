package OrgCoreCore

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"testing"
)

func TestInitBindLast(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	TestCreateOrg(t)
	TestCreateGroup(t)
	TestSetBind(t)
}

func TestGetBindLast(t *testing.T) {
	bindData, err := GetBindLast(&ArgsGetBindLast{
		OrgID:   orgData.ID,
		GroupID: 0,
		Mark:    "test_mark",
		Params:  CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(bindData)
	}
}

func TestClearBindLast(t *testing.T) {
	TestDeleteBind(t)
	TestClear(t)
	TestDeleteGroup(t)
	TestDeleteWorkTime(t)
	TestDeleteOrg(t)
}
