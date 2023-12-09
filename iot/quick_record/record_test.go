package IOTQuickRecord

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newRecordData FieldsRecord
	newRandCode   = ""
)

func TestInitRecord(t *testing.T) {
	TestInit(t)
	newRandCode, _ = CoreFilter.GetRandStr3(10)
}

func TestCreate(t *testing.T) {
	data, err := Create(&ArgsCreate{
		DeviceCode: newRandCode,
		Params:     nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newRecordData = data
	}
}

func TestAudit(t *testing.T) {
	newGroupData, err := IOTDevice.CreateGroup(&IOTDevice.ArgsCreateGroup{
		Mark:       newRandCode,
		Name:       "test",
		Des:        "test des",
		CoverFiles: []int64{},
		Action:     []int64{},
		ExpireTime: 60,
		UseType:    0,
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newGroupData)
	newDeviceData, err := IOTDevice.CreateDevice(&IOTDevice.ArgsCreateDevice{
		Status:     0,
		Name:       newRandCode,
		Des:        "test des",
		CoverFiles: []int64{},
		DesFiles:   []int64{},
		GroupID:    newGroupData.ID,
		Code:       newRandCode,
		Key:        newRandCode,
		Address:    "address",
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newDeviceData)
	err = Audit(&ArgsAudit{
		ID:       newRecordData.ID,
		DeviceID: newDeviceData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetResult(t *testing.T) {
	data, err := GetResult(&ArgsGetResult{
		ID: newRecordData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}
