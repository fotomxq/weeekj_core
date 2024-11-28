package AnalysisIndex

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"testing"
)

var (
	testIndexData FieldsIndex
)

func TestInitIndex(t *testing.T) {
	TestInit(t)
}

func TestGetIndexByCode(t *testing.T) {
	data, _ := GetIndexByCode("test1")
	testIndexData = data
	TestDeleteIndex(t)
}

func TestCreateIndex(t *testing.T) {
	err := CreateIndex(&ArgsCreateIndex{
		Code:        "test1",
		Name:        "测试指标",
		IsSystem:    true,
		Description: "测试指标描述",
		Decision:    "测试指标提示",
		Threshold:   0,
		IsEnable:    false,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("CreateIndex OK")
}

func TestGetIndexList(t *testing.T) {
	dataList, dataCount, err := GetIndexList(&ArgsGetIndexList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		NeedIsSystem: false,
		IsSystem:     false,
		IsRemove:     false,
		Search:       "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dataCount, dataList)
	testIndexData = dataList[0]
}

func TestDeleteIndex(t *testing.T) {
	err := DeleteIndex(testIndexData.ID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("DeleteIndex OK")
}

func TestClearIndex(t *testing.T) {
}
