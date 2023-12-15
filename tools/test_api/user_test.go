package TestAPI

import (
	"testing"

	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

var (
	newUserInfo UserCore.FieldsUserType
)

func TestInit4(t *testing.T) {
	TestInit(t)
}

func TestLoginUserByPassword(t *testing.T) {
	newToken, newKey, _, err := LoginUserByPassword()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new token: ", newToken, " | ", newKey)
	}
}

func TestCreateUser(t *testing.T) {
	data, err := CreateUser()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
		newUserInfo = data
	}
}

func TestUpdateUserGroups(t *testing.T) {
	if err := UpdateUserGroups(newUserInfo.ID, "user"); err != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	if err := DeleteUser(newUserInfo.ID); err != nil {
		t.Error(err)
	}
}
