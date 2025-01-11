package OrgSign

import "testing"

var (
	newSignData  FieldsSign
	newOrgID     int64 = 0
	newOrgBindID int64 = 0
	newUserID    int64 = 0
)

func TestGetInit(t *testing.T) {
	TestInit(t)
	TestCreateSign(t)
}

func TestGetSignAll(t *testing.T) {
	dataList, err := GetSignAll(&ArgsGetSignAll{
		OrgID:     newOrgID,
		OrgBindID: newOrgBindID,
		UserID:    newUserID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList)
		newSignData = dataList[0]
	}
}

func TestGetSignDefault(t *testing.T) {
	data, err := GetSignDefault(&ArgsGetSignDefault{
		OrgID:     newOrgID,
		OrgBindID: newOrgBindID,
		UserID:    newUserID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetSignByID(t *testing.T) {
	data, err := GetSignByID(newSignData.ID)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetClear(t *testing.T) {
	TestDeleteSignByID(t)
	TestClear(t)
}
