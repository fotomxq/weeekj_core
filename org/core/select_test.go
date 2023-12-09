package OrgCoreCore

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	selectBindData FieldsBind
)

func TestInitSelect(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	TestCreateOrg(t)
	TestCreateGroup(t)
	var err error
	selectBindData, err = SetBind(&ArgsSetBind{
		UserID:   newUserInfo.ID,
		Name:     newUserInfo.Name,
		OrgID:    orgData.ID,
		GroupIDs: []int64{orgGroupData.ID},
		Manager:  []string{"all"},
		Params:   CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(selectBindData)
	}
}

func TestSetSelect(t *testing.T) {
	selectBindData, permissions, err := SetSelect(&ArgsSetSelect{
		UserID: selectBindData.UserID,
		OrgID:  orgData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(selectBindData, permissions)
	}
}

func TestGetSelect(t *testing.T) {
	data, err := GetSelect(&ArgsGetSelect{
		UserID: selectBindData.UserID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestClearSelect(t *testing.T) {
	if err := DeleteBind(&ArgsDeleteBind{
		ID:    selectBindData.ID,
		OrgID: 0,
	}); err != nil {
		t.Error(err, ", bind data id: ", bindData.ID)
	}
	TestClear(t)
	TestDeleteGroup(t)
	TestDeleteWorkTime(t)
	TestDeleteOrg(t)
}
