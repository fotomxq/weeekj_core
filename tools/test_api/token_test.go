package TestAPI

import "testing"

func TestInit3(t *testing.T){
	TestInit(t)
}

func TestGetNewToken(t *testing.T) {
	newToken, newKey, err := GetNewToken()
	if err != nil{
		t.Error(err)
	}else{
		t.Log("new token: ", newToken, " | ", newKey)
	}
}
