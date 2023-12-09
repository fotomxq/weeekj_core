package ToolsAppUpdate

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit        = false
	newUpdateData FieldsUpdate
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
}

func TestCreateUpdate(t *testing.T) {
	TestCreateApp(t)
	var err error
	newUpdateData, err = CreateUpdate(&ArgsCreateUpdate{
		OrgID:        0,
		Name:         "pda712",
		Des:          "",
		DesFiles:     []int64{},
		System:       "android_phone",
		SystemVerMin: []int64{7, 1, 0},
		SystemVerMax: []int64{7, 1, 2},
		SystemVer:    []string{},
		AppID:        newAppData.ID,
		Ver:          []int64{1, 0, 17},
		FileID:       8234,
		DownloadURL:  "https://dev-filex.weeekj.com/2022-01-20_14-50-41_ll07ul3r2P9baxy3cOvV6VcBf4Cf.apk",
		GrayscaleRes: false,
		IsSkip:       false,
		Params:       []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newUpdateData)
}

func TestGetUpdateID(t *testing.T) {
	data, err := GetUpdateID(&ArgsGetUpdateID{
		ID: newUpdateData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateUpdate(t *testing.T) {
	err := UpdateUpdate(&ArgsUpdateUpdate{
		ID:           newUpdateData.ID,
		OrgID:        newUpdateData.OrgID,
		Name:         newUpdateData.Name,
		Des:          newUpdateData.Des,
		DesFiles:     newUpdateData.DesFiles,
		System:       newUpdateData.System,
		SystemVerMin: newUpdateData.SystemVerMin,
		SystemVerMax: newUpdateData.SystemVerMax,
		SystemVer:    newUpdateData.SystemVer,
		Ver:          newUpdateData.Ver,
		FileID:       newUpdateData.FileID,
		DownloadURL:  newUpdateData.DownloadURL,
		GrayscaleRes: newUpdateData.GrayscaleRes,
		Params:       newUpdateData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestCheckUpdate(t *testing.T) {
	data, b := CheckUpdate(&ArgsCheckUpdate{
		System:        "android_phone",
		SystemVersion: "7.1.1",
		AppMark:       "transport_pda",
		Version:       "1.0.16",
	})
	if b {
		t.Log("update success, ", data)
	} else {
		t.Error("update failed: ", data)
	}
	data, b = CheckUpdate(&ArgsCheckUpdate{
		System:        "android_phone",
		SystemVersion: "7.1.1",
		AppMark:       "transport_pda",
		Version:       "1.0.17",
	})
	if b {
		t.Error("update failed: ", data)
		t.Error("update failed: check ver: 1.0.17, app update ver: ", data.Ver)
	} else {
		t.Log("update success, ", data)
	}
}

func TestGetCountList(t *testing.T) {
	data, err := GetCountList(&ArgsGetCountList{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTimeCarbon().SubHour().Time,
			MaxTime: CoreFilter.GetNowTimeCarbon().Time,
		},
		TimeType: "hour",
		OrgID:    newUpdateData.OrgID,
		AppID:    newUpdateData.AppID,
		UpdateID: newUpdateData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteUpdate(t *testing.T) {
	err := DeleteUpdate(&ArgsDeleteUpdate{
		ID:    newUpdateData.ID,
		OrgID: newUpdateData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteApp(t)
}
