package UserChat

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"testing"
)

var (
	newUserData UserCore.FieldsUserType
)

func TestInitChat(t *testing.T) {
	TestInitGroup(t)
	TestCreateGroup(t)
	//新的用户
	var err error
	newUserData, _, err = UserCore.CreateUser(&UserCore.ArgsCreateUser{
		OrgID:              TestOrg.OrgData.ID,
		Name:               "测试用户",
		Password:           "",
		NationCode:         "",
		Phone:              "",
		AllowSkipWaitEmail: false,
		Email:              "",
		Username:           fmt.Sprint("test_", CoreFilter.GetRandNumber(1, 1000)),
		Avatar:             0,
		Status:             2,
		Parents:            []UserCore.FieldsUserParent{},
		Groups:             []UserCore.FieldsUserGroupType{},
		Infos:              []CoreSQLConfig.FieldsConfigType{},
		Logins:             []UserCore.FieldsUserLoginType{},
		SortID:             0,
		Tags:               []int64{},
	})
	ToolsTest.ReportData(t, err, newUserData)
}

func TestInviteUser(t *testing.T) {
	err := InviteUser(&ArgsInviteUser{
		GroupID:      newGroupData.ID,
		UserID:       TestOrg.UserInfo.ID,
		InviteUserID: newUserData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetChatList(t *testing.T) {
	dataList, dataCount, err := GetChatList(&ArgsGetChatList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		GroupID:   newGroupData.ID,
		UserID:    -1,
		HaveLeave: false,
		Search:    "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateChatName(t *testing.T) {
	err := UpdateChatName(&ArgsUpdateChatName{
		GroupID: newGroupData.ID,
		UserID:  newUserData.ID,
		Name:    "测试姓名2",
	})
	ToolsTest.ReportError(t, err)
}

func TestOutChat(t *testing.T) {
	err := OutChat(&ArgsOutChat{
		GroupID:      newGroupData.ID,
		UserID:       newGroupData.UserID,
		InviteUserID: newUserData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearChat(t *testing.T) {
	TestDeleteGroup(t)
	TestClearGroup(t)
}
