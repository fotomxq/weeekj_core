package UserCore

import (
	"testing"
)

//验证测试
func TestGetUserData(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
	userData, err := GetUserData(&ArgsGetUserData{
		UserID: userData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(userData)
	}
}

func TestCheckUserPassword(t *testing.T) {
	if err := CheckUserPassword(&ArgsCheckUserPassword{
		UserInfo: &userData, Password: "abc123456789_text",
	}); err == nil {
		t.Error(err)
	}
	if err := CheckUserPassword(&ArgsCheckUserPassword{
		UserInfo: &userData, Password: "test_testXXX",
	}); err != nil {
		t.Error(err)
	}
}

func TestFilterUserData(t *testing.T) {
	data := FilterUserData(&ArgsFilterUserData{
		UserData: DataUserDataType{
			Info:        userData,
			Groups:      nil,
			Permissions: nil,
		},
	})
	t.Log(data)
}

func TestDeleteUserByID2(t *testing.T) {
	TestDeleteUserByID(t)
}
