package DataLakeSource

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"testing"
)

var (
	newTestFieldID int64
)

func TestInitFields(t *testing.T) {
	TestInitTable(t)
	TestCreateTable(t)
}

func TestCreateFields(t *testing.T) {
	newID, err := CreateFields(&ArgsCreateFields{
		TableID:       newTestTableID,
		InputName:     "测试字段",
		InputType:     "string",
		InputLength:   0,
		InputDefault:  "",
		InputRequired: false,
		InputPattern:  "",
		FieldName:     "test_field",
		FieldLabel:    "测试字段",
		IsPrimary:     false,
		IsIndex:       false,
		IsSystem:      false,
		IsSearch:      false,
		DataType:      "text",
		FieldDesc:     "测试字段",
	})
	if err != nil {
		t.Error(err)
		return
	}
	newTestFieldID = newID
	t.Log("newTestFieldID:", newTestFieldID)
}

func TestGetFieldsDetail(t *testing.T) {
	newData, err := GetFieldsDetail(newTestFieldID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(newData)
}

func TestGetFieldsDetailByTableIDAndFieldName(t *testing.T) {
	newData, err := GetFieldsDetailByTableIDAndFieldName(newTestTableID, "test_field")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(newData)
}

func TestGetFieldsList(t *testing.T) {
	dataList, dataCount, err := GetFieldsList(&ArgsGetFieldsList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		TableID:   -1,
		FieldName: "",
		InputType: "",
		DataType:  "",
		IsRemove:  false,
		Search:    "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dataList, dataCount)
}

func TestGetFieldsListByTableID(t *testing.T) {
	dataList, dataCount, err := GetFieldsListByTableID(newTestTableID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dataList, dataCount)
}

func TestUpdateFields(t *testing.T) {
	err := UpdateFields(&ArgsUpdateFields{
		ID:            newTestFieldID,
		InputName:     "测试字段",
		InputType:     "string",
		InputLength:   0,
		InputDefault:  "",
		InputRequired: false,
		InputPattern:  "",
		IsIndex:       false,
		IsSearch:      false,
		FieldDesc:     "测试字段",
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteFields(t *testing.T) {
	err := DeleteFields(newTestFieldID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClearFields(t *testing.T) {
	err := ClearFields(newTestTableID)
	if err != nil {
		t.Error(err)
		return
	}
	TestClearTable(t)
}
