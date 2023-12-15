package MapRoom

import (
	"testing"
	"time"

	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ServiceUserInfo "github.com/fotomxq/weeekj_core/v5/service/user_info"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
)

var (
	isInit      = false
	newRoomData FieldsRoom
	newSortData ClassSort.FieldsSort
	newInfoData ServiceUserInfo.FieldsInfo
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
		Init()
		TestOrg.LocalCreateBind(t)
		var err error
		newSortData, err = Sort.Create(&ClassSort.ArgsCreate{
			BindID:      TestOrg.OrgData.ID,
			Mark:        "",
			ParentID:    0,
			CoverFileID: 0,
			DesFiles:    []int64{},
			Name:        "测试分类",
			Des:         "",
			Params:      []CoreSQLConfig.FieldsConfigType{},
		})
		ToolsTest.ReportData(t, err, newSortData)
		newInfoData, err = ServiceUserInfo.CreateInfo(&ServiceUserInfo.ArgsCreateInfo{
			OrgID:                 TestOrg.OrgData.ID,
			UserID:                TestOrg.UserInfo.ID,
			BindID:                0,
			BindType:              0,
			Name:                  "测试姓名",
			Country:               86,
			Gender:                0,
			IDCard:                "177754222542",
			Phone:                 "16544445541",
			CoverFileID:           0,
			DesFiles:              []int64{},
			Address:               "",
			DateOfBirth:           time.Time{},
			MaritalStatus:         false,
			EducationStatus:       0,
			Profession:            "",
			Level:                 0,
			EmergencyContact:      "",
			EmergencyContactPhone: "",
			SortID:                0,
			Tags:                  []int64{},
			DocID:                 0,
			Params:                []CoreSQLConfig.FieldsConfigType{},
		})
		ToolsTest.ReportData(t, err, newInfoData)
	}
	isInit = true
}

func TestCreateRoom(t *testing.T) {
	var err error
	newRoomData, err = CreateRoom(&ArgsCreateRoom{
		OrgID:       TestOrg.OrgData.ID,
		SortID:      newSortData.ID,
		Tags:        []int64{},
		Infos:       []int64{},
		Code:        "M-0301",
		Name:        "房间0301",
		Des:         "房间描述",
		CoverFileID: 0,
		DesFiles:    []int64{},
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newRoomData)
}

func TestGetRoomList(t *testing.T) {
	dataList, dataCount, err := GetRoomList(&ArgsGetRoomList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:  -1,
		SortID: -1,
		Tags:   []int64{},
		Status: -1,
		InfoID: -1,
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetRoomID(t *testing.T) {
	data, err := GetRoomID(&ArgsGetRoomID{
		ID:    newRoomData.ID,
		OrgID: newRoomData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetRooms(t *testing.T) {
	dataList, err := GetRooms(&ArgsGetRooms{
		IDs:        []int64{newRoomData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetRoomsName(t *testing.T) {
	data, err := GetRoomsName(&ArgsGetRooms{
		IDs:        []int64{newRoomData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateRoom(t *testing.T) {
	errCode, err := UpdateRoom(&ArgsUpdateRoom{
		ID:          newRoomData.ID,
		OrgID:       newRoomData.OrgID,
		SortID:      newRoomData.SortID,
		Tags:        newRoomData.Tags,
		Code:        newRoomData.Code,
		Name:        newRoomData.Name,
		Des:         newRoomData.Des,
		CoverFileID: newRoomData.CoverFileID,
		DesFiles:    newRoomData.DesFiles,
		Params:      newRoomData.Params,
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Error("errCode: ", errCode)
	}
}

func TestUpdateServiceSortBindGroup(t *testing.T) {
	err := UpdateServiceSortBindGroup(&ArgsUpdateServiceSortBindGroup{
		ID:      newSortData.ID,
		OrgID:   newSortData.BindID,
		GroupID: TestOrg.GroupData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateServiceSortMissionExpire(t *testing.T) {
	err := UpdateServiceSortMissionExpire(&ArgsUpdateServiceSortMissionExpire{
		ID:         newSortData.ID,
		OrgID:      newSortData.BindID,
		ExpireTime: 1800,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateStatus(t *testing.T) {
	//空闲
	errCode, err := UpdateStatus(&ArgsUpdateStatus{
		ID:     newRoomData.ID,
		OrgID:  newRoomData.OrgID,
		Status: 0,
		Infos:  []int64{},
	})
	if err != nil {
		t.Error("err: ", err, ", errCode: ", errCode)
		return
	} else {
		TestGetRoomID(t)
	}
	//有人
	errCode, err = UpdateStatus(&ArgsUpdateStatus{
		ID:     newRoomData.ID,
		OrgID:  newRoomData.OrgID,
		Status: 1,
		Infos:  []int64{newInfoData.ID},
	})
	if err != nil {
		t.Error("err: ", err, ", errCode: ", errCode)
		return
	} else {
		TestGetRoomID(t)
	}
	//退房
	errCode, err = UpdateStatus(&ArgsUpdateStatus{
		ID:     newRoomData.ID,
		OrgID:  newRoomData.OrgID,
		Status: 2,
		Infos:  []int64{},
	})
	if err != nil {
		t.Error("err: ", err, ", errCode: ", errCode)
		return
	} else {
		TestGetRoomID(t)
	}
	//禁用
	errCode, err = UpdateStatus(&ArgsUpdateStatus{
		ID:     newRoomData.ID,
		OrgID:  newRoomData.OrgID,
		Status: 3,
		Infos:  []int64{},
	})
	if err != nil {
		t.Error("err: ", err, ", errCode: ", errCode)
		return
	} else {
		TestGetRoomID(t)
	}
	//清理
	errCode, err = UpdateStatus(&ArgsUpdateStatus{
		ID:     newRoomData.ID,
		OrgID:  newRoomData.OrgID,
		Status: 4,
		Infos:  []int64{},
	})
	if err != nil {
		t.Error("err: ", err, ", errCode: ", errCode)
		return
	} else {
		TestGetRoomID(t)
	}
}

func TestUpdateServiceStatus(t *testing.T) {
	//空置
	errCode, err := UpdateServiceStatus(&ArgsUpdateServiceStatus{
		ID:               newRoomData.ID,
		OrgID:            newRoomData.OrgID,
		ServiceStatus:    0,
		ServiceBindID:    0,
		ServiceMissionID: 0,
	})
	if err != nil {
		t.Error("err: ", err, ", errCode: ", errCode)
		return
	} else {
		TestGetRoomID(t)
	}
	//呼叫
	errCode, err = UpdateServiceStatus(&ArgsUpdateServiceStatus{
		ID:               newRoomData.ID,
		OrgID:            newRoomData.OrgID,
		ServiceStatus:    1,
		ServiceBindID:    0,
		ServiceMissionID: 0,
	})
	if err != nil {
		t.Error("err: ", err, ", errCode: ", errCode)
		return
	} else {
		TestGetRoomID(t)
	}
	//反馈
	errCode, err = UpdateServiceStatus(&ArgsUpdateServiceStatus{
		ID:               newRoomData.ID,
		OrgID:            newRoomData.OrgID,
		ServiceStatus:    2,
		ServiceBindID:    0,
		ServiceMissionID: 0,
	})
	if err != nil {
		t.Error("err: ", err, ", errCode: ", errCode)
		return
	} else {
		TestGetRoomID(t)
	}
}

func TestGetLogList(t *testing.T) {
	dataList, dataCount, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:     -1,
		RoomID:    -1,
		Status:    -1,
		InfoID:    -1,
		MissionID: -1,
		Search:    "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDeleteRoom(t *testing.T) {
	err := DeleteRoom(&ArgsDeleteRoom{
		ID:    newRoomData.ID,
		OrgID: newRoomData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteClear(t *testing.T) {
	var err error
	TestOrg.LocalClear(t)
	err = Sort.DeleteByID(&ClassSort.ArgsDeleteByID{
		ID:     newSortData.ID,
		BindID: newSortData.BindID,
	})
	ToolsTest.ReportError(t, err)
	err = ServiceUserInfo.DeleteInfo(&ServiceUserInfo.ArgsDeleteInfo{
		ID:    newInfoData.ID,
		OrgID: newInfoData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
