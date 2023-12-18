package UserCore

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	isInit = false

	userData FieldsUserType
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestCreateUser(t *testing.T) {
	TestCreatePermission(t)
	TestCreateGroup(t)
	var err error
	var errCode string
	var parents FieldsUserParents
	userData, errCode, err = CreateUser(&ArgsCreateUser{
		OrgID:                0,
		Name:                 "test_test",
		Password:             "test_testXXX",
		NationCode:           "86",
		Phone:                "17000000001",
		AllowSkipPhoneVerify: true,
		AllowSkipWaitEmail:   true,
		Email:                "xxxxx@qq.com",
		Username:             "username",
		Avatar:               0,
		Status:               2,
		Parents:              parents,
		Groups:               nil,
		Infos:                nil,
		Logins:               nil,
		SortID:               0,
		Tags:                 nil,
	})
	if err != nil {
		//尝试通过手机号获取，可能是违反重复键
		if err.Error() == "get last id pq: 重复键违反唯一约束\"user_core_phone_uindex\"" {
			TestGetUserByPhone(t)
		} else {
			t.Error(errCode, ", err: ", err)
		}
	} else {
		t.Log(userData)
	}
}

func TestGetUserByEmail(t *testing.T) {
	data, err := GetUserByEmail(&ArgsGetUserByEmail{
		OrgID: 0,
		Email: "xxxxx@qq.com",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetUserByID(t *testing.T) {
	var err error
	userData, err = GetUserByID(&ArgsGetUserByID{
		ID: userData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(userData)
	}
}

func TestGetUserByPhone(t *testing.T) {
	var err error
	userData, err = GetUserByPhone(&ArgsGetUserByPhone{
		OrgID:      0,
		NationCode: "86", Phone: "17000000001",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(userData)
	}
}

func TestGetUserByUsername(t *testing.T) {
	data, err := GetUserByUsername(&ArgsGetUserByUsername{
		OrgID:    0,
		Username: "username",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetUserByLogin(t *testing.T) {
	TestInit(t)
	val, err := CoreFilter.GetRandStr3(10)
	if err != nil {
		t.Error(err)
		return
	}
	val = fmt.Sprint("weixin_", val)
	if err := UpdateUserLoginByID(&ArgsUpdateUserLoginByID{
		ID:       userData.ID,
		OrgID:    0,
		Mark:     "weixin",
		Val:      val,
		Config:   "weixin_config_" + CoreFilter.GetRandStr(30),
		IsRemove: false,
	}); err != nil {
		t.Error(err)
		return
	}
	data, err := GetUserByLogin(&ArgsGetUserByLogin{
		OrgID: 0,
		Mark:  "weixin",
		Val:   val,
	})
	if err != nil {
		TestGetUserByID(t)
		t.Log("params: ", &ArgsGetUserByLogin{
			OrgID: 0,
			Mark:  "weixin",
			Val:   val,
		})
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetUserList(t *testing.T) {
	data, dataCount, err := GetUserList(&ArgsGetUserList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1, Max: 10, Sort: "id", Desc: true,
		},
		Status:       -1,
		OrgID:        0,
		ParentSystem: "",
		ParentID:     0,
		SortID:       0,
		Tags:         []int64{},
		IsRemove:     false,
		Search:       "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataCount, data)
	}
}

func TestUpdateUserParent(t *testing.T) {
	var parents FieldsUserParents
	if err := UpdateUserParent(&ArgsUpdateUserParent{
		ID:      userData.ID,
		OrgID:   0,
		Parents: parents,
	}); err != nil {
		t.Error(err)
	}
}

func TestUpdateUserEmailByID(t *testing.T) {
	if err := UpdateUserEmailByID(&ArgsUpdateUserEmailByID{
		ID:        userData.ID,
		OrgID:     0,
		Email:     "xxx2@qq.com",
		AllowSkip: false,
	}); err != nil {
		t.Error(err)
	}
}

func TestUpdateUserGroupByID(t *testing.T) {
	TestCreateGroup(t)
	if err := UpdateUserGroupByID(&ArgsUpdateUserGroupByID{
		ID:       userData.ID,
		OrgID:    0,
		GroupID:  groupData.ID,
		ExpireAt: CoreFilter.GetNowTime().Add(time.Second * 3000),
		IsRemove: false,
	}); err != nil {
		t.Error(err)
	} else {
		TestGetUserByID(t)
		t.Log(userData)
	}
}

func TestUpdateUserInfoByID(t *testing.T) {
	if err := UpdateUserInfoByID(&ArgsUpdateUserInfoByID{
		ID:     userData.ID,
		OrgID:  0,
		Name:   "name_update",
		Avatar: 0,
	}); err != nil {
		t.Error(err)
	}
}

func TestUpdateUserInfosByID(t *testing.T) {
	if err := UpdateUserInfosByID(&ArgsUpdateUserInfosByID{
		ID:       userData.ID,
		OrgID:    0,
		Mark:     "test_info1",
		Val:      "888",
		IsRemove: false,
	}); err == nil {
		TestGetUserByID(t)
		t.Log(userData)
	} else {
		t.Error(err)
	}
	if err := UpdateUserInfosByID(&ArgsUpdateUserInfosByID{
		ID:       userData.ID,
		OrgID:    0,
		Mark:     "test_info1",
		Val:      "999",
		IsRemove: true,
	}); err == nil {
		TestGetUserByID(t)
		t.Log(userData)
	} else {
		t.Error(err)
	}
}

func TestUpdateUserLoginByID(t *testing.T) {
	if err := UpdateUserLoginByID(&ArgsUpdateUserLoginByID{
		ID:       userData.ID,
		OrgID:    0,
		Mark:     "weixin2",
		Val:      "weixin_test2",
		Config:   "weixin_config2",
		IsRemove: false,
	}); err != nil {
		t.Error(err)
	}
	TestGetUserByID(t)
	if err := UpdateUserLoginByID(&ArgsUpdateUserLoginByID{
		ID:       userData.ID,
		OrgID:    0,
		Mark:     "weixin2",
		Val:      "weixin_test2",
		Config:   "weixin_config2",
		IsRemove: true,
	}); err != nil {
		t.Error(err)
	}
	TestGetUserByID(t)
}

func TestUpdateUserPhoneByID(t *testing.T) {
	if err := UpdateUserPhoneByID(&ArgsUpdateUserPhoneByID{
		ID:         userData.ID,
		OrgID:      0,
		NationCode: "86",
		Phone:      "16372635552",
	}); err != nil {
		t.Error(err)
	}
}

func TestUpdateUserStatus(t *testing.T) {
	if err := UpdateUserStatus(&ArgsUpdateUserStatus{
		ID:     userData.ID,
		OrgID:  0,
		Status: 1,
	}); err != nil {
		t.Error(err)
	}
}

func TestUpdateUserUsernameByID(t *testing.T) {
	if err := UpdateUserUsernameByID(&ArgsUpdateUserUsernameByID{
		ID:       userData.ID,
		OrgID:    0,
		Username: "new_username_update_x",
	}); err != nil {
		t.Error(err)
	}
}

func TestUpdateUserPasswordByID(t *testing.T) {
	if err := UpdateUserPasswordByID(&ArgsUpdateUserPasswordByID{
		ID:       userData.ID,
		OrgID:    0,
		Password: "new_password_password",
	}); err != nil {
		t.Error(err)
	}
}

func TestDeleteUserByID(t *testing.T) {
	if err := DeleteUserByID(&ArgsDeleteUserByID{
		ID:    userData.ID,
		OrgID: 0,
	}); err != nil {
		t.Error(err)
	}
	TestDeleteGroup(t)
	TestDeletePermission2(t)
}
