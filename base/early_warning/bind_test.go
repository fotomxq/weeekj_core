package BaseEarlyWarning

import (
	"testing"
)

func TestInit3(t *testing.T) {
	TestInit(t)
	TestCreateTemplate(t)
	TestCreateTo(t)
}

//绑定关系
func TestSetBind(t *testing.T) {
	var err error
	bindData, err = SetBind(&ArgsSetBind{
		ToID: toData.ID, TemplateID: templateData.ID, Level: 1, LevelMode: "none", NextWaitTime: "31s", NeedSMS: true, NeedEmail: true,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(bindData)
	}
}

func TestGetBindID(t *testing.T) {
	var err error
	bindData, err = GetBindID(&ArgsGetBindID{
		ID: bindData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(bindData)
	}
}

func TestGetBindByToID(t *testing.T) {
	getData, err := GetBindByToID(&ArgsGetBindByToID{
		ToID: toData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(getData)
	}
}

func TestGetBindByTemplateID(t *testing.T) {
	getData, err := GetBindByTemplateID(&ArgsGetBindByTemplateID{
		TemplateID: templateData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(getData)
	}
}

func TestSetUnBind(t *testing.T) {
	if err := SetUnBind(&ArgsSetUnBind{
		ToID: toData.ID, TemplateID: templateData.ID,
	}); err != nil {
		t.Error(err)
	}
}

func TestClear(t *testing.T) {
	TestDeleteTemplate(t)
	TestDeleteToByID(t)
}
