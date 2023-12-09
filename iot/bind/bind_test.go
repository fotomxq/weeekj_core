package IOTBind

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newBindData FieldsBind
)

func TestInitBind(t *testing.T) {
	TestInit(t)
}

func TestSetBind(t *testing.T) {
	var err error
	newBindData, err = SetBind(&ArgsSetBind{
		OrgID:    1,
		DeviceID: 1,
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "room",
			ID:     1,
			Mark:   "",
			Name:   "测试房间",
		},
		Params: CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newBindData)
}

func TestGetBindList(t *testing.T) {
	dataList, dataCount, err := GetBindList(&ArgsGetBindList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		DeviceID: -1,
		FromInfo: CoreSQLFrom.FieldsFrom{},
		IsRemove: false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetBindFrom(t *testing.T) {
	data, err := GetBindFrom(&ArgsGetBindFrom{
		OrgID: 1,
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "room",
			ID:     1,
			Mark:   "",
			Name:   "",
		},
	})
	ToolsTest.ReportData(t, err, data)
	data, err = GetBindFrom(&ArgsGetBindFrom{
		OrgID: 1,
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "room",
			ID:     1,
			Mark:   "*",
			Name:   "",
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get * binds: ", data)
	}
}

func TestGetBindDevice(t *testing.T) {
	data, err := GetBindDevice(1, 1)
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteBind(t *testing.T) {
	err := DeleteBind(&ArgsDeleteBind{
		ID:    newBindData.ID,
		OrgID: newBindData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
