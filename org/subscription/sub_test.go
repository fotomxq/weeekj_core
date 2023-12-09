package OrgSubscription

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

func TestInitSub(t *testing.T) {
	TestInit(t)
}

func TestSetSub(t *testing.T) {
	TestInit(t)
	TestSetConfig(t)
	errCode, err := SetSub(&ArgsSetSub{
		ConfigUnit: 2,
		OrgID:      TestOrg.OrgData.ID,
		ConfigID:   newConfigData.ID,
		Params:     []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error("code: ", errCode, ", err: ", err)
	} else {
		TestGetSubList(t)
		t.Log("设置成功")
	}
	TestDeleteSub(t)
}

func TestGetSubList(t *testing.T) {
	dataList, dataCount, err := GetSubList(&ArgsGetSubList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:        0,
		ConfigID:     0,
		NeedIsExpire: false,
		IsExpire:     false,
		IsRemove:     false,
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newSub = dataList[0]
		t.Log("expire: ", newSub.ExpireAt)
	}
}

func TestCheckSub(t *testing.T) {
	expireAt, b := CheckSub(&ArgsCheckSub{
		OrgID:    123,
		ConfigID: newConfigData.ID,
	})
	if !b {
		t.Error(expireAt)
	} else {
		t.Log("b: ", b, ", expire at: ", expireAt)
	}
}

func TestDeleteSub(t *testing.T) {
	err := DeleteSub(&ArgsDeleteSub{
		ID: newSub.ID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteConfig(t)
}
