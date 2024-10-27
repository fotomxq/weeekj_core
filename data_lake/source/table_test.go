package DataLakeSource

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"testing"
)

var (
	newTestTableID int64
)

func TestInitTable(t *testing.T) {
	TestInit(t)
	//检查test_table表是否存在，如果已经存在说明测试异常，自动删除
	data, err := GetTableDetailByName("test_table")
	if err == nil {
		_ = DeleteTable(data.ID)
	}
}

func TestCreateTable(t *testing.T) {
	newID, err := CreateTable(&ArgsCreateTable{
		TableName:      "test_table",
		TableDesc:      "测试表",
		TipName:        "测试表",
		ChannelName:    "default",
		ChannelTipName: "默认渠道，主要用于标记渠道来源/数据源",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("new id: ", newID)
	newTestTableID = newID
}

func TestGetTableDetail(t *testing.T) {
	data, err := GetTableDetail(newTestTableID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestGetTableDetailByName(t *testing.T) {
	data, err := GetTableDetailByName("test_table")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestGetTableList(t *testing.T) {
	dataList, dataCount, err := GetTableList(&ArgsGetTableList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		IsRemove: false,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dataList, dataCount)
}

func TestUpdateTable(t *testing.T) {
	err := UpdateTable(&ArgsUpdateTable{
		ID:             newTestTableID,
		TableName:      "test_table",
		TableDesc:      "测试表",
		TipName:        "测试表",
		ChannelName:    "default",
		ChannelTipName: "默认渠道，主要用于标记渠道来源/数据源",
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteTable(t *testing.T) {
	TestClearFields(t)
	err := DeleteTable(newTestTableID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestClearTable(t *testing.T) {
	TestClear(t)
}
