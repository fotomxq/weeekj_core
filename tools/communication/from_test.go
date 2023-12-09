package ToolsCommunication

import (
	"testing"

	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
)

var (
	newFromData FieldsFrom
)

func TestInitFrom(t *testing.T) {
	TestInitRoom(t)
}

func TestAppendRoom(t *testing.T) {
	TestCreateRoom(t)
	var err error
	newFromData, err = AppendRoom(&ArgsAppendRoom{
		RoomID:     newRoomData.ID,
		FromSystem: 0,
		FromID:     1,
		Name:       "测试名称",
		Token:      "",
		AllowSend:  false,
		Role:       0,
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newFromData)
}

func TestGetFrom(t *testing.T) {
	data, err := GetFrom(&ArgsGetFrom{
		FromSystem: 0,
		FromID:     1,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetFromList(t *testing.T) {
	dataList, dataCount, err := GetFromList(&ArgsGetFromList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		RoomID:        -1,
		FromSystem:    -1,
		FromID:        -1,
		NeedAllowSend: false,
		AllowSend:     false,
		Role:          0,
		Search:        "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestCheckFromAndRoom(t *testing.T) {
	data, err := CheckFromAndRoom(&ArgsCheckFromAndRoom{
		RoomID:     newRoomData.ID,
		FromSystem: 0,
		FromID:     1,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateFromExpire(t *testing.T) {
	err := UpdateFromExpire(&ArgsUpdateFromExpire{
		FromSystem: 0,
		FromID:     1,
	})
	ToolsTest.ReportError(t, err)
}

func TestOutRoom(t *testing.T) {
	err := OutRoom(&ArgsOutRoom{
		ID:         newFromData.ID,
		FromSystem: 0,
		FromID:     1,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteRoom(t)
}

func TestAppendRoomTwo(t *testing.T) {
	data1, data2, data3, err := AppendRoomTwo(&ArgsAppendRoomTwo{
		OrgID:       1,
		ConnectType: 2,
		SortID:      0,
		Tags:        []int64{},
		Name:        "测试房间",
		Des:         "",
		CoverFileID: 0,
		Password:    "",
		FromSystem:  0,
		FromID:      1,
		FromName:    "测试用户1",
		ToSystem:    0,
		ToID:        2,
		ToName:      "测试用户2",
	})
	if err != nil {
		t.Error(err)
		return
	} else {
		t.Log("room: ", data1)
		t.Log("token 1: ", data2)
		t.Log("token 2: ", data3)
	}
	err = OutRoom(&ArgsOutRoom{
		ID:         newFromData.ID,
		FromSystem: data2.FromSystem,
		FromID:     data2.FromID,
	})
	ToolsTest.ReportError(t, err)
}
