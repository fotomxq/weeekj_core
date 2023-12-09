package IOTError

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit     = false
	newErrorID int64
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestCreate(t *testing.T) {
	err := Create(&ArgsCreate{
		SendEW:   false,
		OrgID:    1,
		GroupID:  1,
		DeviceID: 1,
		Code:     "001",
		Content:  "001test",
		Params:   []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:       0,
		GroupID:     0,
		DeviceID:    0,
		AllowDone:   false,
		Done:        false,
		AllowSendEW: false,
		SendEW:      false,
		TimeBetween: CoreSQLTime.FieldsCoreTime{},
		IsHistory:   false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if len(dataList) > 0 {
		newErrorID = dataList[0].ID
	}
}

func TestUpdateDone(t *testing.T) {
	err := UpdateDone(&ArgsUpdateDone{
		IDs:   []int64{newErrorID},
		OrgID: -1,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteByID(t *testing.T) {
	err := DeleteByID(&ArgsDeleteByID{
		IDs:   []int64{newErrorID},
		OrgID: -1,
	})
	ToolsTest.ReportError(t, err)
}
