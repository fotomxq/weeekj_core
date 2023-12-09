package TMSUserRunning

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	ToolsTestUserRole "gitee.com/weeekj/weeekj_core/v5/tools/test_user_role"
	"testing"
)

var (
	newMissionData FieldsMission
)

func TestMissionInit(t *testing.T) {
	TestInit(t)
}

func TestCreateMission(t *testing.T) {
	var err error
	newMissionData, err = CreateMission(&ArgsCreateMission{
		RunType:      0,
		WaitAt:       CoreFilter.GetISOByTime(CoreFilter.GetNowTime()),
		GoodType:     "none",
		OrgID:        0,
		UserID:       TestOrg.UserInfo.ID,
		OrderID:      0,
		RunWaitPrice: 10,
		RunPayAfter:  false,
		Des:          "测试跑腿",
		GoodWidget:   10,
		FromAddress: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   10010,
			City:       10010,
			Address:    "测试地址",
			MapType:    0,
			Longitude:  1,
			Latitude:   1,
			Name:       "测试姓名",
			NationCode: "86",
			Phone:      "17777777777",
		},
		ToAddress: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   10010,
			City:       10010,
			Address:    "测试地址",
			MapType:    0,
			Longitude:  1,
			Latitude:   1,
			Name:       "测试姓名",
			NationCode: "86",
			Phone:      "17777777777",
		},
		Params: nil,
	})
	ToolsTest.ReportData(t, err, newMissionData)
}

func TestGetMissionID(t *testing.T) {
	var err error
	newMissionData, err = GetMissionID(&ArgsGetMissionID{
		ID:     newMissionData.ID,
		OrgID:  -1,
		UserID: -1,
		RoleID: -1,
	})
	ToolsTest.ReportData(t, err, newMissionData)
}

func TestGetMissionAllInfoID(t *testing.T) {
	var err error
	newMissionData, err = GetMissionAllInfoID(&ArgsGetMissionID{
		ID:     newMissionData.ID,
		OrgID:  -1,
		UserID: -1,
		RoleID: -1,
	})
	ToolsTest.ReportData(t, err, newMissionData)
}

func TestGetMissionList(t *testing.T) {
	dataList, dataCount, err := GetMissionList(&ArgsGetMissionList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		RunType:        -1,
		OrgID:          -1,
		UserID:         -1,
		NeedIsTake:     false,
		IsTake:         false,
		NeedIsFinish:   false,
		IsFinish:       false,
		OrderID:        -1,
		RoleID:         -1,
		NeedIsRunPay:   false,
		IsRunPay:       false,
		NeedHaveRunPay: false,
		HaveRunPay:     false,
		IsRemove:       false,
		Search:         "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetMissionTakeCodeByID(t *testing.T) {
	data, err := GetMissionTakeCodeByID(&ArgsGetMissionID{
		ID:     newMissionData.ID,
		OrgID:  -1,
		UserID: -1,
		RoleID: -1,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateMission(t *testing.T) {
	err := UpdateMission(&ArgsUpdateMission{
		ID:            newMissionData.ID,
		WaitAt:        CoreFilter.GetISOByTime(newMissionData.WaitAt),
		GoodType:      newMissionData.GoodType,
		TakeAt:        CoreFilter.GetISOByTime(newMissionData.TakeAt),
		FinishAt:      CoreFilter.GetISOByTime(newMissionData.FinishAt),
		TakeCode:      newMissionData.TakeCode,
		RunType:       newMissionData.RunType,
		OrgID:         newMissionData.OrgID,
		UserID:        newMissionData.UserID,
		OrderID:       newMissionData.OrderID,
		RoleID:        newMissionData.RoleID,
		RunPayAt:      CoreFilter.GetISOByTime(newMissionData.RunPayAt),
		RunPayID:      newMissionData.RunPayID,
		RunPrice:      newMissionData.RunPrice,
		RunWaitPrice:  newMissionData.RunWaitPrice,
		RunPayAfter:   newMissionData.RunPayAfter,
		OrderPayAfter: newMissionData.OrderPayAfter,
		OrderPrice:    newMissionData.OrderPrice,
		OrderPayAt:    newMissionData.OrderDes,
		OrderPayID:    newMissionData.OrderPayID,
		Des:           newMissionData.Des,
		OrderDesFiles: newMissionData.OrderDesFiles,
		OrderDes:      newMissionData.OrderDes,
		GoodWidget:    newMissionData.GoodWidget,
		FromAddress:   newMissionData.FromAddress,
		ToAddress:     newMissionData.ToAddress,
		Params:        newMissionData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateMissionRunner(t *testing.T) {
	err := UpdateMissionRunner(&ArgsUpdateMissionRunner{
		ID:     newMissionData.ID,
		RoleID: ToolsTestUserRole.RoleData.ID,
		Des:    "测试更换跑腿员",
	})
	ToolsTest.ReportError(t, err)
}

func TestTakeMission(t *testing.T) {
	err := TakeMission(&ArgsTakeMission{
		ID:     newMissionData.ID,
		RoleID: newMissionData.RoleID,
	})
	ToolsTest.ReportError(t, err)
}

func TestFinishMission(t *testing.T) {
	err := FinishMission(&ArgsFinishMission{
		ID:       newMissionData.ID,
		RoleID:   newMissionData.RoleID,
		TakeCode: newMissionData.TakeCode,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteMission(t *testing.T) {
	err := DeleteMission(&ArgsDeleteMission{
		ID:     newMissionData.ID,
		UserID: newMissionData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestMissionClear(t *testing.T) {
	TestClear(t)
}
