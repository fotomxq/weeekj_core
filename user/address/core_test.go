package UserAddress

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newData FieldsAddress
	userID  int64 = 123
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestCreate(t *testing.T) {
	var err error
	newData, err = Create(&ArgsCreate{
		UserID:     userID,
		NiceName:   "测试昵称",
		Country:    86,
		Province:   10,
		City:       10010,
		Address:    "测试地址",
		MapType:    0,
		Longitude:  33.11,
		Latitude:   112.22,
		Name:       "测试姓名",
		NationCode: "86",
		Phone:      "1766666666",
		Email:      "adb@abc.com",
		Infos: []CoreSQLConfig.FieldsInfoType{
			{
				Mark: "a1",
				Name: "测试",
				Val:  "a1value",
			},
		},
	})
	ToolsTest.ReportData(t, err, newData)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		ParentID:    0,
		UserID:      0,
		Country:     0,
		Province:    0,
		City:        0,
		IsRemove:    false,
		SearchPhone: "",
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetCount(t *testing.T) {
	count, err := GetCount(&ArgsGetCount{
		UserID: userID,
	})
	ToolsTest.ReportData(t, err, count)
}

func TestGetID(t *testing.T) {
	data, err := GetID(&ArgsGetID{
		ID:       newData.ID,
		UserID:   0,
		IsRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetIDs(t *testing.T) {
	data, err := GetIDs(&ArgsGetIDs{
		IDs:      []int64{newData.ID},
		IsRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdate(t *testing.T) {
	err := Update(&ArgsUpdate{
		ID:         newData.ID,
		UserID:     userID,
		NiceName:   "测试昵称edit",
		Country:    86,
		Province:   10,
		City:       10010,
		Address:    "测试地址edit",
		MapType:    0,
		Longitude:  33.11,
		Latitude:   112.22,
		Name:       "测试姓名edit",
		NationCode: "86",
		Phone:      "1766666666",
		Email:      "adb@abcedit.com",
		Infos: []CoreSQLConfig.FieldsInfoType{
			{
				Mark: "a1edit",
				Name: "测试edit",
				Val:  "a1value_edit",
			},
		},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetIDTop(t *testing.T) {
	var err error
	newData, err = GetIDTop(&ArgsGetIDTop{
		ID:     newData.ID,
		UserID: 0,
	})
	ToolsTest.ReportData(t, err, newData)
}

func TestDelete(t *testing.T) {
	err := Delete(&ArgsDelete{
		ID:     newData.ID,
		UserID: 0,
	})
	t.Log("delete id: ", newData.ID)
	ToolsTest.ReportError(t, err)
}
