package AnalysisIndex

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"testing"
)

var (
	testParamsData FieldsIndexParam
)

func TestInitParams(t *testing.T) {
	TestInit(t)
	TestCreateIndex(t)
	TestGetIndexList(t)
}

func TestCreateIndexParams(t *testing.T) {
	err := CreateIndexParam(&ArgsCreateIndexParam{
		IndexID:  1,
		Name:     "测试名称",
		Code:     "FH",
		ParamVal: "test",
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetIndexParamList(t *testing.T) {
	dataList, dataCount, err := GetIndexParamList(&ArgsGetIndexParamList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		IndexID:  testIndexData.ID,
		IsRemove: false,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dataCount, dataList)
	testParamsData = dataList[0]
}

func TestDeleteIndexParam(t *testing.T) {
	err := DeleteIndexParam(testParamsData.ID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClearParams(t *testing.T) {
	TestDeleteIndex(t)
	TestClearIndex(t)
}
