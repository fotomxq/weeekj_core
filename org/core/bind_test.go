package OrgCoreCore

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	bindData FieldsBind
	userID   int64 = 123
)

func TestInit2(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	TestCreateOrg(t)
	TestCreateGroup(t)
}

func TestSetBind(t *testing.T) {
	var err error
	bindData, err = SetBind(&ArgsSetBind{
		UserID:   userID,
		OrgID:    orgData.ID,
		Name:     "测试绑定姓名",
		GroupIDs: []int64{orgGroupData.ID},
		Manager:  []string{"all"},
		Params:   CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(bindData)
	}
}

// 单独测试超过11个授权，是否支持？
func TestSetBind3(t *testing.T) {
	var err error
	bindData, err = SetBind(&ArgsSetBind{
		UserID:   userID,
		Name:     "测试绑定姓名",
		OrgID:    orgData.ID,
		GroupIDs: []int64{orgGroupData.ID},
		Manager:  []string{"member", "work-organizational-organizer-view", "work-organizational-organizer-edit", "work-organizational-organizer-delete", "work-organizational-organizer-operate", "work-organizational-organizer-config", "work-organizational-organizer-clear", "work-time", "work-time-organizer-view", "work-time-organizer-edit", "work-time-organizer-delete", "work-time-organizer-operate"},
		Params:   nil,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(bindData)
	}
	//返回原先授权
	TestSetBind(t)
}

func TestGetBind(t *testing.T) {
	bindData, err := GetBind(&ArgsGetBind{
		ID:     bindData.ID,
		OrgID:  0,
		UserID: 0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(bindData)
	}
}

func TestGetBinds(t *testing.T) {
	bindList, err := GetBinds(&ArgsGetBinds{
		IDs: []int64{bindData.ID},
	})
	ToolsTest.ReportData(t, err, bindList)
}

func TestGetBindsName(t *testing.T) {
	bindList, err := GetBindsName(&ArgsGetBinds{
		IDs:        []int64{bindData.ID},
		HaveRemove: true,
	})
	ToolsTest.ReportData(t, err, bindList)
}

func TestGetBindList(t *testing.T) {
	dataList, dataCount, err := GetBindList(&ArgsGetBindList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    orgData.ID,
		UserID:   0,
		GroupID:  0,
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

// 带有用户性质的检查
func TestGetBindListByUserID(t *testing.T) {
	dataList, dataCount, err := GetBindList(&ArgsGetBindList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    orgData.ID,
		UserID:   userID,
		GroupID:  0,
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

func TestGetBindByCreateInfoAndOrg(t *testing.T) {
	data, err := GetBindByUser(&ArgsGetBindByUser{
		UserID: bindData.UserID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestSetBindParams(t *testing.T) {
	t.Log("wait bindData id: ", bindData.OrgID, ", user id: ", bindData.UserID, ", bind id: ", bindData.ID)
	err := SetBindParams(&ArgsSetBindParams{
		OrgID:  bindData.OrgID,
		UserID: bindData.UserID,
		Params: CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "a1",
				Val:  "v1",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	TestGetBind(t)
}

func TestSetBind2(t *testing.T) {
	TestDeleteBind(t)
	var err error
	bindData, err = SetBind(&ArgsSetBind{
		UserID:   userID,
		Name:     "测试姓名2",
		OrgID:    orgData.ID,
		GroupIDs: []int64{orgGroupData.ID},
		Manager:  []string{"all"},
		Params:   CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(bindData)
	}
}

func TestCheckBindAndOrg(t *testing.T) {
	if err := CheckBindAndOrg(&ArgsCheckBindAndOrg{
		ID:    bindData.ID,
		OrgID: bindData.OrgID,
	}); err != nil {
		t.Error("not in : ", err)
	}
	if err := CheckBindAndOrg(&ArgsCheckBindAndOrg{
		ID:    bindData.ID,
		OrgID: 12321312312,
	}); err == nil {
		t.Error("in: ", err)
	}
}

func TestDeleteBind(t *testing.T) {
	if err := DeleteBind(&ArgsDeleteBind{
		ID:    bindData.ID,
		OrgID: 0,
	}); err != nil {
		t.Error(err, ", bind data id: ", bindData.ID)
	}
}

func TestClear4(t *testing.T) {
	TestClear(t)
	TestDeleteGroup(t)
	TestDeleteWorkTime(t)
	TestDeleteOrg(t)
}
