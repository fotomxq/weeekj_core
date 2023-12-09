package OrgUser

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

func TestInitUser(t *testing.T) {
	TestInit(t)
}

func TestUpdateUserData(t *testing.T) {
	_, _ = UpdateUserData(&ArgsUpdateUserData{
		OrgID:         TestOrg.OrgData.ID,
		UserID:        TestOrg.UserInfo.ID,
		UserAddressID: -1,
	})
}

func TestGetUserDataList(t *testing.T) {
	dataList, dataCount, err := GetUserDataList(&ArgsGetUserDataList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:       TestOrg.OrgData.ID,
		SearchPhone: "",
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestClearUser(t *testing.T) {
	TestClear(t)
}
