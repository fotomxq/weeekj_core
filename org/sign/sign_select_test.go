package OrgSign

import "testing"

func TestSelectInit(t *testing.T) {
	TestGetInit(t)
}

func TestSelectSignDefault(t *testing.T) {
	err := SelectSignDefault(&ArgsSelectSignDefault{
		ID:        newSignData.ID,
		OrgID:     newSignData.OrgID,
		OrgBindID: newSignData.OrgBindID,
		UserID:    newSignData.UserID,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetSignLastTemp(t *testing.T) {
	TestDeleteSignByID(t)
	errCode, err := CreateSign(&ArgsCreateSign{
		OrgID:     newOrgID,
		OrgBindID: newOrgBindID,
		UserID:    newUserID,
		IsDefault: false,
		SignType:  "base64",
		IsTemp:    true,
		SignData:  "abc",
		FileID:    0,
	})
	if err != nil {
		t.Error(errCode, ", ", err)
	}
	TestGetSignAll(t)
	data, err := GetSignLastTemp(&ArgsGetSignLastTemp{
		OrgID:     newSignData.OrgID,
		OrgBindID: newSignData.OrgBindID,
		UserID:    newSignData.UserID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetSignLastTempAndDefault(t *testing.T) {
	errCode, err := CreateSign(&ArgsCreateSign{
		OrgID:     newOrgID,
		OrgBindID: newOrgBindID,
		UserID:    newUserID,
		IsDefault: false,
		SignType:  "base64",
		IsTemp:    true,
		SignData:  "abc",
		FileID:    0,
	})
	if err != nil {
		t.Error(errCode, ", ", err)
	}
	TestGetSignAll(t)
	data, err := GetSignLastTempAndDefault(&ArgsGetSignLastTempAndDefault{
		OrgID:     newSignData.OrgID,
		OrgBindID: newSignData.OrgBindID,
		UserID:    newSignData.UserID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestSelectClear(t *testing.T) {
	TestGetClear(t)
}
