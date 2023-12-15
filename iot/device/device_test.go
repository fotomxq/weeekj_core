package IOTDevice

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newDeviceData FieldsDevice
)

func TestInitDevice(t *testing.T) {
	TestInit(t)
}

func TestCreateDevice(t *testing.T) {
	TestCreateGroup(t)
	var err error
	newDeviceData, err = CreateDevice(&ArgsCreateDevice{
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

func TestGetDeviceList(t *testing.T) {
	dataList, dataCount, err := GetDeviceList(&ArgsGetDeviceList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		GroupID:  newGroupData.ID,
		SortID:   -1,
		Tags:     []int64{},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetDeviceByID(t *testing.T) {
	data, err := GetDeviceByID(&ArgsGetDeviceByID{
		ID:    newDeviceData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetDeviceByCode(t *testing.T) {
	data, err := GetDeviceByCode(&ArgsGetDeviceByCode{
		GroupMark: newGroupData.Mark,
		Code:      newDeviceData.Code,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetDeviceMore(t *testing.T) {
	data, err := GetDeviceMore(&ArgsGetDeviceMore{
		IDs:        []int64{newDeviceData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetDeviceMoreMap(t *testing.T) {
	data, err := GetDeviceMoreMap(&ArgsGetDeviceMore{
		IDs:        []int64{newDeviceData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateDevice(t *testing.T) {
	err := UpdateDevice(&ArgsUpdateDevice{
		ID:         newDeviceData.ID,
		OrgID:      -1,
		Status:     1,
		Name:       "test_device",
		Des:        "test des",
		CoverFiles: []int64{},
		DesFiles:   []int64{},
		GroupID:    newGroupData.ID,
		Code:       "code_xxx1",
		Address:    "address",
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteDevice(t *testing.T) {
	err := DeleteDevice(&ArgsDeleteDevice{
		ID:    newDeviceData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteGroup(t)
}
