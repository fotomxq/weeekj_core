package BaseConfig

import (
	"testing"
)

func TestInit2(t *testing.T) {
	TestInit(t)
}

func TestLoadGroupData(t *testing.T) {
	err := loadGroupData()
	if err != nil {
		t.Error(err)
	}
}

func TestGetGroupData(t *testing.T) {
	data := GetGroupData()
	t.Log(data)
}
