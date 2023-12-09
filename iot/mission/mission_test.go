package IOTMission

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newActionData  IOTDevice.FieldsAction
	newGroupData   IOTDevice.FieldsGroup
	newDeviceData  IOTDevice.FieldsDevice
	newMissionData FieldsMission
)

func TestInitMission(t *testing.T) {
	TestInit(t)
	var err error
	newActionData, err = IOTDevice.CreateAction(&IOTDevice.ArgsCreateAction{
		Mark:        "test_mark",
		Name:        "测试",
		Des:         "test mark",
		ExpireTime:  60,
		ConnectType: "none",
		Configs:     []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newActionData)
	newGroupData, err = IOTDevice.CreateGroup(&IOTDevice.ArgsCreateGroup{
		Mark:       "test_mark",
		Name:       "test",
		Des:        "test des",
		CoverFiles: []int64{},
		Action:     []int64{newActionData.ID},
		ExpireTime: 60,
		UseType:    0,
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newGroupData)
	newDeviceData, err = IOTDevice.CreateDevice(&IOTDevice.ArgsCreateDevice{
		Status:     0,
		Name:       "test_device",
		Des:        "test des",
		CoverFiles: []int64{},
		DesFiles:   []int64{},
		GroupID:    newGroupData.ID,
		Code:       "code_xxx1",
		Address:    "address",
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newDeviceData)
}

func TestCreateMission(t *testing.T) {
	data, errCode, err := CreateMission(&ArgsCreateMission{
		OrgID:      -1,
		DeviceID:   newDeviceData.ID,
		ParamsData: []byte(string("test adc")),
		Action:     newActionData.Mark,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newMissionData = data
	} else {
		t.Error("errCode: ", errCode)
	}
}

func TestGetMissionList(t *testing.T) {
	dataList, dataCount, err := GetMissionList(&ArgsGetMissionList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:       -1,
		GroupID:     -1,
		DeviceID:    -1,
		Status:      -1,
		Action:      "",
		ConnectType: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetWaitMissionByDevice(t *testing.T) {
	data, err := GetWaitMissionByDevice(&ArgsGetWaitMissionByDevice{
		DeviceID: newDeviceData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateMissionStatus(t *testing.T) {
	err := UpdateMissionStatus(&ArgsUpdateMissionStatus{
		ID:     newMissionData.ID,
		Status: 1,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateMissionFinish(t *testing.T) {
	err := UpdateMissionFinish(&ArgsUpdateMissionFinish{
		ID:         newMissionData.ID,
		ReportData: []byte("result..."),
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteMission(t *testing.T) {
	err := IOTDevice.DeleteGroup(&IOTDevice.ArgsDeleteGroup{
		ID: newGroupData.ID,
	})
	ToolsTest.ReportError(t, err)
	err = IOTDevice.DeleteAction(&IOTDevice.ArgsDeleteAction{
		ID: newActionData.ID,
	})
	ToolsTest.ReportError(t, err)
	err = IOTDevice.DeleteDevice(&IOTDevice.ArgsDeleteDevice{
		ID:    newDeviceData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportError(t, err)
}
