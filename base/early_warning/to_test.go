package BaseEarlyWarning

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	"testing"
)

func TestInit4(t *testing.T) {
	TestInit(t)
}

func TestCreateTo(t *testing.T) {
	var err error
	toData, err = CreateTo(&ArgsCreateTo{
		UserID: 10, Name: "送达测试人", Des: "送达人描述", PhoneNationCode: "86", Phone: "17635705566", Email: "fotomxq@qq.com",
	})
	if err != nil {
		if err.Error() != "user id have bind to data" {
			t.Error(err)
		} else {
			TestGetToByUserID(t)
		}
	} else {
		t.Log("toData: ", toData)
	}
}

func TestGetToByID(t *testing.T) {
	getData, err := GetToByID(&ArgsGetToByID{
		ID: toData.ID,
	})
	if err != nil {
		t.Error(toData.ID, err)
	} else {
		t.Log(getData)
	}
}

func TestGetToByUserID(t *testing.T) {
	var err error
	toData, err = GetToByUserID(&ArgsGetToByUserID{
		UserID: 10,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(toData)
	}
}

func TestUpdateTo(t *testing.T) {
	if err := UpdateTo(&ArgsUpdateTo{
		ID: toData.ID, UserID: toData.ID, Name: "送达人修改1", Des: "送达人描述修改1", PhoneNationCode: "86", Phone: "17635705567", Email: "fotomxq@qq.com",
	}); err != nil {
		t.Error(err)
	} else {
		getData, err := GetBindByToID(&ArgsGetBindByToID{
			ToID: toData.ID,
		})
		if err != nil {
			t.Error(err)
		} else {
			t.Log(getData)
		}
	}
}

func TestGetToList(t *testing.T) {
	getDataList, count, err := GetToList(&ArgsGetToList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("count:", count, ", list: ", getDataList)
	}
}

func TestDeleteToByID(t *testing.T) {
	if err := DeleteToByID(&ArgsDeleteToByID{
		ID: toData.ID,
	}); err != nil {
		t.Error(err)
	}
}
