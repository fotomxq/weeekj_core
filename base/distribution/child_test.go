package BaseDistribution

import (
	"testing"
)

func TestInit2(t *testing.T) {
	TestInit(t)
}

func TestSetChild(t *testing.T) {
	if err := SetChild(&ArgsSetChild{
		Mark: "test", Name: "测试负载", IP: "1.1.1.1", Port: "9000",
	}); err == nil {
		t.Error("delete child failed")
	}
	TestSetService(t)
	if err := SetChild(&ArgsSetChild{
		Mark: "test", Name: "测试负载", IP: "1.1.1.1", Port: "9000",
	}); err != nil {
		t.Error(err)
	}
}

func TestGetChildAll(t *testing.T) {
	data, err := GetChildAll(&ArgsGetChildAll{
		Mark: "test",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetBalancing(t *testing.T) {
	data, err := GetBalancing(&ArgsGetBalancing{
		Mark: "test",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestDeleteChild(t *testing.T) {
	if err := DeleteChild(&ArgsDeleteChild{
		Mark: "test", IP: "1.1.1.1", Port: "9000",
	}); err != nil {
		t.Error(err)
	}
}
