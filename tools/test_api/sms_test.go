package TestAPI

import "testing"

func TestSendSMS(t *testing.T) {
	TestGetNewToken(t)
	if err := SendSMS("", ""); err != nil{
		t.Error(err)
	}
}