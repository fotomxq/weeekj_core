package TestAPI

import "testing"

func TestInit2(t *testing.T){
	TestInit(t)
}

func TestGetConfig(t *testing.T) {
	data, err := GetConfig([]string{"AppName"})
	if err != nil{
		t.Error(err)
	}else{
		t.Log(data)
	}
}