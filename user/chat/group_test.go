package UserChat

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newGroupData FieldsGroup
)

func TestInitGroup(t *testing.T) {
	TestInit(t)
}

func TestCreateGroup(t *testing.T) {
	var err error
	newGroupData, err = CreateGroup(&ArgsCreateGroup{
		OrgID:            TestOrg.OrgData.ID,
		Name:             "测试房间",
		UserID:           TestOrg.UserInfo.ID,
		OnlyCreateInvite: false,
		Params:           nil,
	})
	ToolsTest.ReportData(t, err, newGroupData)
}

func TestGetGroupList(t *testing.T) {
	dataList, dataCount, err := GetGroupList(&ArgsGetGroupList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:  -1,
		UserID: -1,
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDeleteGroup(t *testing.T) {
	err := DeleteGroup(&ArgsDeleteGroup{
		ID:     newGroupData.ID,
		OrgID:  newGroupData.OrgID,
		UserID: newGroupData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearGroup(t *testing.T) {
	TestClear(t)
}
