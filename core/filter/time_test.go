package CoreFilter

import "testing"

func TestGetNowTime(t *testing.T) {
	t1, err := GetTimeByISO("2022-08-30T16:00:00.000Z")
	if err != nil {
		t.Error(err)
	} else {
		t.Log("time: ", t1)
	}
}
