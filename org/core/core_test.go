package OrgCoreCore

import (
	OrgTime "github.com/fotomxq/weeekj_core/v5/org/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"testing"
)

var (
	isInit       = false
	orgData      FieldsOrg
	workTimeData OrgTime.FieldsWorkTime
	newUserInfo  UserCore.FieldsUserType
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
}

func TestCreateUser(t *testing.T) {
	//创建用户
	var err error
	var errCode string
	newUserInfo, errCode, err = UserCore.CreateUser(&UserCore.ArgsCreateUser{
		Name:       "测试用户",
		Password:   "",
		NationCode: "",
		Phone:      "",
		Email:      "",
		Username:   "",
		Status:     2,
		Parents:    nil,
		Groups:     nil,
		Infos:      nil,
		Logins:     nil,
	})
	if err != nil {
		t.Error(errCode, err)
		return
	} else {
		t.Log("new user data, ", newUserInfo)
		userID = newUserInfo.ID
	}
}

func TestClear(t *testing.T) {
	var err error
	err = UserCore.DeleteUserByID(&UserCore.ArgsDeleteUserByID{
		ID: newUserInfo.ID,
	})
	if err != nil {
		t.Error(err, ", delete user by id: ", newUserInfo.ID)
	}
}
