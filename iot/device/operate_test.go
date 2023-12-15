package IOTDevice

import (
	"testing"
	"time"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

var (
	newOperateData FieldsOperate
)

func TestInitOperate(t *testing.T) {
	TestInit(t)
	TestCreateDevice(t)
}

func TestSetOperate(t *testing.T) {
	err := SetOperate(&ArgsSetOperate{
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Second * 10),
		OrgID:       1,
		Permissions: []string{"all"},
		Action:      newGroupData.Action,
		DeviceID:    newDeviceData.ID,
		Address:     "test address",
		SortID:      -1,
		Tags:        []int64{},
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetOperateList(t *testing.T) {
	dataList, dataCount, err := GetOperateList(&ArgsGetOperateList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:  -1,
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetDeviceList2(t *testing.T) {
	dataList, dataCount, err := GetDeviceList(&ArgsGetDeviceList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    newOperateData.OrgID,
		GroupID:  -1,
		SortID:   -1,
		Tags:     []int64{},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetDeviceByID2(t *testing.T) {
	data, err := GetDeviceByID(&ArgsGetDeviceByID{
		ID:    newDeviceData.ID,
		OrgID: newOperateData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestCheckOperate(t *testing.T) {
	data, err := CheckOperate(&ArgsCheckOperate{
		DeviceID: newDeviceData.ID,
		OrgID:    1,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newOperateData = data
	}
}

func TestCheckOperates(t *testing.T) {
	err := CheckOperates(&ArgsCheckOperates{
		DeviceIDs: []int64{newDeviceData.ID},
		OrgID:     1,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetOperateActionByDeviceID(t *testing.T) {
	dataList, err := GetOperateActionByDeviceID(&ArgsGetOperateActionByDeviceID{
		DeviceID: newDeviceData.ID,
		OrgID:    -1,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetOperateActionByDeviceID2(t *testing.T) {
	dataList, err := GetOperateActionByDeviceID(&ArgsGetOperateActionByDeviceID{
		DeviceID: newDeviceData.ID,
		OrgID:    newOperateData.OrgID,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestDeleteOperate(t *testing.T) {
	err := DeleteOperate(&ArgsDeleteOperate{
		DeviceIDs: []int64{newDeviceData.ID},
		OrgID:     newOperateData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteDevice(t)
}
