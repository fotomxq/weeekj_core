package OrgMission

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	newAutoData FieldsAuto
)

func TestInitAuto(t *testing.T) {
	TestInit(t)
}

func TestCreateAuto(t *testing.T) {
	data, err := CreateAuto(&ArgsCreateAuto{
		TimeType:     0,
		TimeN:        []int64{},
		SkipHoliday:  false,
		StartHour:    0,
		StartMinute:  0,
		EndHour:      0,
		EndMinute:    0,
		OrgID:        1,
		CreateBindID: 0,
		BindID:       0,
		OtherBindIDs: []int64{},
		Title:        "",
		Des:          "test des",
		DesFiles:     []int64{},
		StartAt:      time.Time{},
		EndAt:        time.Time{},
		TipID:        0,
		Level:        0,
		SortID:       0,
		Tags:         []int64{},
		Params:       []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newAutoData = data
	}
}

func TestGetAutoList(t *testing.T) {
	dataList, dataCount, err := GetAutoList(&ArgsGetAutoList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:        -1,
		CreateBindID: -1,
		BindID:       -1,
		OtherBindID:  -1,
		Level:        -1,
		SortID:       -1,
		Tags:         -1,
		IsRemove:     false,
		Search:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateAuto(t *testing.T) {
	err := UpdateAuto(&ArgsUpdateAuto{
		ID:            newAutoData.ID,
		OrgID:         newAutoData.OrgID,
		OperateBindID: newAutoData.CreateBindID,
		TimeType:      newAutoData.TimeType,
		TimeN:         newAutoData.TimeN,
		SkipHoliday:   newAutoData.SkipHoliday,
		StartHour:     newAutoData.StartHour,
		StartMinute:   newAutoData.StartMinute,
		EndHour:       newAutoData.EndHour,
		EndMinute:     newAutoData.EndMinute,
		BindID:        newAutoData.BindID,
		OtherBindIDs:  newAutoData.OtherBindIDs,
		Title:         newAutoData.Title,
		Des:           newAutoData.Des,
		DesFiles:      newAutoData.DesFiles,
		StartAt:       newAutoData.StartAt,
		EndAt:         newAutoData.EndAt,
		TipID:         newAutoData.TipID,
		Level:         newAutoData.Level,
		SortID:        newAutoData.SortID,
		Tags:          newAutoData.Tags,
		Params:        newAutoData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteAuto(t *testing.T) {
	err := DeleteAuto(&ArgsDeleteAuto{
		ID:            newAutoData.ID,
		OrgID:         newAutoData.OrgID,
		OperateBindID: newAutoData.CreateBindID,
	})
	ToolsTest.ReportError(t, err)
}
