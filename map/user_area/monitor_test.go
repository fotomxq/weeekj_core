package MapUserArea

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newMonitorData FieldsMonitor
)

func TestInitMonitor(t *testing.T) {
	TestInit(t)
}

func TestCreateMonitor(t *testing.T) {
	data, err := CreateMonitor(&ArgsCreateMonitor{
		OrgID:      1,
		UserInfoID: 1,
		DeviceID:   1,
		AreaID:     1,
		OrgGroupID: 1,
		Params:     []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newMonitorData = data
	}
}

func TestGetMonitorList(t *testing.T) {
	dataList, dataCount, err := GetMonitorList(&ArgsGetMonitorList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		NeedIsInvalid: false,
		IsInvalid:     false,
		OrgID:         -1,
		UserInfoID:    -1,
		DeviceID:      -1,
		AreaID:        -1,
		InRange:       false,
		OrgGroupID:    -1,
		IsRemove:      false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateMonitor(t *testing.T) {
	err := UpdateMonitor(&ArgsUpdateMonitor{
		ID:         newMonitorData.ID,
		OrgID:      newMonitorData.OrgID,
		UserInfoID: newMonitorData.UserInfoID,
		DeviceID:   newMonitorData.DeviceID,
		AreaID:     newMonitorData.AreaID,
		OrgGroupID: newMonitorData.OrgGroupID,
		Params:     newMonitorData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteMonitor(t *testing.T) {
	err := DeleteMonitor(&ArgsDeleteMonitor{
		ID:    newMonitorData.ID,
		OrgID: newMonitorData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
