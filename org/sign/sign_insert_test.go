package OrgSign

import "testing"

func TestInsertInit(t *testing.T) {
	TestInit(t)
}

func TestCreateSign(t *testing.T) {
	errCode, err := CreateSign(&ArgsCreateSign{
		OrgID:     newOrgID,
		OrgBindID: newOrgBindID,
		UserID:    newUserID,
		IsDefault: false,
		SignType:  "base64",
		IsTemp:    false,
		SignData:  "abc",
		FileID:    0,
	})
	if err != nil {
		t.Error(errCode, ", ", err)
	}
	TestGetSignAll(t)
}

func TestInsertClear(t *testing.T) {
	TestDeleteSignByID(t)
	TestClear(t)
}
