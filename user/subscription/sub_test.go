package UserSubscription

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

func TestInitSub(t *testing.T) {
	TestInitConfig(t)
	TestCreateConfig(t)
}

func TestSetSub(t *testing.T) {
	err := SetSub(&ArgsSetSub{
		OrgID:       TestOrg.OrgData.ID,
		ConfigID:    newConfigData.ID,
		UserID:      TestOrg.UserInfo.ID,
		ExpireAt:    CoreFilter.GetNowTimeCarbon().AddDay().Time,
		HaveExpire:  false,
		UseFrom:     "test",
		UseFromName: "测试",
	})
	ToolsTest.ReportError(t, err)
	err = SetSub(&ArgsSetSub{
		OrgID:       TestOrg.OrgData.ID,
		ConfigID:    newConfigData.ID,
		UserID:      TestOrg.UserInfo.ID,
		ExpireAt:    CoreFilter.GetNowTimeCarbon().AddDays(3).Time,
		HaveExpire:  false,
		UseFrom:     "test",
		UseFromName: "测试",
	})
	ToolsTest.ReportError(t, err)
	err = SetSub(&ArgsSetSub{
		OrgID:       TestOrg.OrgData.ID,
		ConfigID:    newConfigData.ID,
		UserID:      TestOrg.UserInfo.ID - 1,
		ExpireAt:    CoreFilter.GetNowTimeCarbon().AddDay().Time,
		HaveExpire:  false,
		UseFrom:     "test",
		UseFromName: "测试",
	})
	ToolsTest.ReportError(t, err)
}

func TestGetSubList(t *testing.T) {
	dataList, dataCount, err := GetSubList(&ArgsGetSubList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:        -1,
		ConfigID:     -1,
		UserID:       TestOrg.UserInfo.ID,
		NeedIsExpire: false,
		IsExpire:     false,
		IsRemove:     false,
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetSub(t *testing.T) {
	data, err := GetSub(&ArgsGetSub{
		ConfigID: newConfigData.ID,
		UserID:   TestOrg.UserInfo.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestCheckSub(t *testing.T) {
	b := CheckSub(&ArgsCheckSub{
		ConfigID: newConfigData.ID,
		UserID:   TestOrg.UserInfo.ID,
	})
	if !b {
		t.Error("check sub false")
	} else {
		t.Log("check sub true")
	}
}

func TestUseSub(t *testing.T) {
	err := UseSub(&ArgsUseSub{
		ConfigID: newConfigData.ID,
		UserID:   TestOrg.UserInfo.ID,
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Error("config id: ", newConfigData.ID, ", user id: ", TestOrg.UserInfo.ID)
		TestGetSub(t)
	}
}

func TestClearSubByConfig(t *testing.T) {
	err := ClearSubByConfig(&ArgsClearSubByConfig{
		OrgID:    newConfigData.OrgID,
		ConfigID: newConfigData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearSubByUser(t *testing.T) {
	err := ClearSubByUser(&ArgsClearSubByUser{
		OrgID:    newConfigData.OrgID,
		UserID:   TestOrg.UserInfo.ID,
		ConfigID: newConfigData.ID,
	})
	ToolsTest.ReportError(t, err)
}
