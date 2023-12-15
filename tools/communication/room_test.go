package ToolsCommunication

import (
	"testing"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

var (
	newRoomData FieldsRoom
)

func TestInitRoom(t *testing.T) {
	TestInit(t)
}

func TestCreateRoom(t *testing.T) {
	var err error
	newRoomData, err = CreateRoom(&ArgsCreateRoom{
		ExpireAt:    CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
		ConnectType: 2,
		SortID:      0,
		Tags:        []int64{},
		OrgID:       0,
		Name:        "测试房间",
		Des:         "测试房间des",
		CoverFileID: 0,
		IsPublic:    false,
		Password:    "",
		MaxCount:    1,
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newRoomData)
}

func TestGetRoomList(t *testing.T) {
	dataList, dataCount, err := GetRoomList(&ArgsGetRoomList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		SortID:      0,
		Tags:        []int64{},
		ConnectType: 0,
		OrgID:       0,
		IsPublic:    false,
		IsRemove:    false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetRoom(t *testing.T) {
	data, err := GetRoom(&ArgsGetRoom{
		ID:    newRoomData.ID,
		OrgID: newRoomData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetRoomMore(t *testing.T) {
	dataList, err := GetRoomMore(&ArgsGetRoomMore{
		IDs:        []int64{newRoomData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestUpdateRoomExpire(t *testing.T) {
	err := UpdateRoomExpire(&ArgsUpdateRoomExpire{
		ID:     newRoomData.ID,
		OrgID:  newRoomData.OrgID,
		FromID: -1,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateRoom(t *testing.T) {
	err := UpdateRoom(&ArgsUpdateRoom{
		ID:          newRoomData.ID,
		OrgID:       newRoomData.OrgID,
		FromID:      -1,
		SortID:      newRoomData.SortID,
		Tags:        newRoomData.Tags,
		ExpireAt:    newRoomData.ExpireAt,
		Name:        newRoomData.Name,
		Des:         newRoomData.Des,
		CoverFileID: newRoomData.CoverFileID,
		IsPublic:    newRoomData.IsPublic,
		Password:    newRoomData.Password,
		MaxCount:    newRoomData.MaxCount,
		Params:      newRoomData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteRoom(t *testing.T) {
	err := DeleteRoom(&ArgsDeleteRoom{
		ID:     newRoomData.ID,
		OrgID:  newRoomData.OrgID,
		FromID: -1,
	})
	ToolsTest.ReportError(t, err)
}
