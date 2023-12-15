package OrgCoreCore

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	"testing"
)

var (
	operateData FieldsOperate
)

func TestInit9(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	TestCreateOrg(t)
	TestCreateGroup(t)
	TestSetBind(t)
}

func TestSetOperate(t *testing.T) {
	err := SetOperate(&ArgsSetOperate{
		OrgID: orgData.ID,
		BindInfo: CoreSQLFrom.FieldsFrom{
			System: "bind",
			ID:     bindData.ID,
			Mark:   "",
			Name:   bindData.Name,
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "device",
			ID:     bindData.ID,
			Mark:   "",
			Name:   "测试设备控制",
		},
		Manager: []string{"all"},
		Params: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "test_1",
				Val:  "value_1",
			},
		},
	})
	if err != nil {
		t.Error(err, ", org id: ", orgData.ID)
	}
}

func TestGetOperate(t *testing.T) {
	var err error
	operateData, err = GetOperate(&ArgsGetOperate{
		OrgID: orgData.ID,
		BindInfo: CoreSQLFrom.FieldsFrom{
			System: "bind",
			ID:     bindData.ID,
			Mark:   "",
			Name:   bindData.Name,
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "device",
			ID:     bindData.ID,
			Mark:   "",
			Name:   "测试设备控制",
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(operateData)
	}
}

func TestGetOperateList(t *testing.T) {
	dataList, dataCount, err := GetOperateList(&ArgsGetOperateList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    0,
		BindInfo: CoreSQLFrom.FieldsFrom{},
		From:     CoreSQLFrom.FieldsFrom{},
		Manager:  "",
		IsRemove: false,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList, dataCount)
	}
}

func TestCheckOperate(t *testing.T) {
	data, b, err := CheckOperate(&ArgsCheckOperate{
		OrgID: orgData.ID,
		BindInfo: CoreSQLFrom.FieldsFrom{
			System: "bind",
			ID:     bindData.ID,
			Mark:   "",
			Name:   bindData.Name,
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "device",
			ID:     bindData.ID,
			Mark:   "",
			Name:   "测试设备控制",
		},
		Managers: []string{"all"},
		Filter:   "or",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(b, data)
	}
}

func TestCheckOperateOnlyBool(t *testing.T) {
	b := CheckOperateOnlyBool(&ArgsCheckOperate{
		OrgID: orgData.ID,
		BindInfo: CoreSQLFrom.FieldsFrom{
			System: "bind",
			ID:     bindData.ID,
			Mark:   "",
			Name:   bindData.Name,
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "device",
			ID:     bindData.ID,
			Mark:   "",
			Name:   "测试设备控制",
		},
		Managers: []string{"all"},
		Filter:   "or",
	})
	if !b {
		t.Error(b)
	}
}

func TestDeleteOperate(t *testing.T) {
	err := DeleteOperate(&ArgsDeleteOperate{
		ID:    operateData.ID,
		OrgID: orgData.ID,
	})
	if err != nil {
		t.Error(err, ", id: ", operateData.ID)
	}
}

func TestClear3(t *testing.T) {
	TestDeleteBind(t)
	TestDeleteGroup(t)
	TestDeleteWorkTime(t)
	TestDeleteOrg(t)
	TestClear(t)
}
