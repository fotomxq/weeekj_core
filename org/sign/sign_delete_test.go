package OrgSign

import "testing"

func TestDeleteInit(t *testing.T) {
	TestInsertInit(t)
	TestCreateSign(t)
}

func TestDeleteSignByID(t *testing.T) {
	err := DeleteSignByID(&ArgsDeleteSignByID{
		ID:        newSignData.ID,
		OrgID:     newSignData.OrgID,
		OrgBindID: newSignData.OrgBindID,
		UserID:    newSignData.UserID,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteClear(t *testing.T) {
	TestInsertClear(t)
}
