package CoreReg

import (
	"testing"
)

var (
	localCode = ""
	localKey  = ""
)

func TestInit(t *testing.T) {
	Init("app1.0")
}

func TestGetCode(t *testing.T) {
	var err error
	localCode, err = GetCode()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(localCode)
	}
}

func TestGetKey(t *testing.T) {
	localKey = GetKey(localCode, "202001", "202012")
	t.Log(localKey)
}

func TestVerify(t *testing.T) {
	if !Verify(localKey) {
		t.Error("not this key: " + localKey)
	}
}
