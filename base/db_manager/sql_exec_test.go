package BaseDBManager

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"testing"
	"time"
)

var (
	newSQLData FieldsSQL
)

func TestSQLExecInit(t *testing.T) {
	TestInit(t)
}

func TestCreateSQL(t *testing.T) {
	err := CreateSQL(&FieldsSQL{
		ID:         0,
		CreateAt:   time.Time{},
		UpdateAt:   time.Time{},
		DeleteAt:   time.Time{},
		FromSystem: "test",
		FromModule: "test",
		FromCode:   "test",
		CarbonCode: "",
		PostURL:    "",
		SQLData:    "",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetSQLList(t *testing.T) {
	dataList, dataCount, err := GetSQLList(&ArgsGetSQLList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		FromSystem: "",
		FromModule: "",
		FromCode:   "",
		IsRemove:   false,
		Search:     "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dataList)
	t.Log(dataCount)
	newSQLData = dataList[0]
}

func TestUpdateSQL(t *testing.T) {
	err := UpdateSQL(&FieldsSQL{
		ID:         newSQLData.ID,
		CreateAt:   time.Time{},
		UpdateAt:   time.Time{},
		DeleteAt:   time.Time{},
		FromSystem: "test",
		FromModule: "test",
		FromCode:   "test",
		CarbonCode: "",
		PostURL:    "",
		SQLData:    "",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetSQLByID(t *testing.T) {
	data, err := GetSQLByID(newSQLData.ID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestGetSQLByCode(t *testing.T) {
	data, err := GetSQLByCode(newSQLData.FromSystem, newSQLData.FromModule, newSQLData.FromCode)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestDeleteSQLByID(t *testing.T) {
	err := DeleteSQLByID(newSQLData.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestSQLExecClear(t *testing.T) {
	TestClear(t)
}
