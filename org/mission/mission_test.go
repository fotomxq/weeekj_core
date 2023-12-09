package OrgMission

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	newMissionData FieldsMission
)

func TestInitMission(t *testing.T) {
	TestInit(t)
}

func TestCreateMission(t *testing.T) {
	data, err := CreateMission(&ArgsCreateMission{
		OrgID:        1,
		CreateBindID: 0,
		BindID:       0,
		OtherBindIDs: []int64{},
		Title:        "测试任务",
		Des:          "test des",
		DesFiles:     []int64{},
		StartAt:      time.Time{},
		EndAt:        time.Time{},
		TipID:        0,
		Level:        0,
		SortID:       0,
		Tags:         []int64{},
		Params:       nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newMissionData = data
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
		OrgID:         -1,
		OperateBindID: -1,
		CreateBindID:  -1,
		BindID:        -1,
		OtherBindID:   -1,
		Level:         -1,
		SortID:        -1,
		Tags:          -1,
		IsRemove:      false,
		Search:        "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateMission(t *testing.T) {
	err := UpdateMission(&ArgsUpdateMission{
		ID:            newMissionData.ID,
		OrgID:         newMissionData.OrgID,
		OperateBindID: newMissionData.CreateBindID,
		Status:        1,
		BindID:        newMissionData.BindID,
		OtherBindIDs:  newMissionData.OtherBindIDs,
		Title:         newMissionData.Title,
		Des:           newMissionData.Des,
		DesFiles:      newMissionData.DesFiles,
		StartAt:       newMissionData.StartAt,
		EndAt:         newMissionData.EndAt,
		TipID:         newMissionData.TipID,
		Level:         newMissionData.Level,
		SortID:        newMissionData.SortID,
		Tags:          newMissionData.Tags,
		Params:        newMissionData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteMission(t *testing.T) {
	err := DeleteMission(&ArgsDeleteMission{
		ID:            newMissionData.ID,
		OrgID:         newMissionData.OrgID,
		OperateBindID: newMissionData.CreateBindID,
	})
	ToolsTest.ReportError(t, err)
}
