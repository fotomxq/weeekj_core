package BaseConfigColumn

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestSet(t *testing.T) {
	err := Set(&ArgsSet{
		System: 0,
		Mark:   "test",
		BindID: 0,
		Data: FieldsChildList{
			{
				Mark: "a1",
				Name: "v1",
				Params: CoreSQLConfig.FieldsConfigsType{
					{
						Mark: "p1",
						Val:  "pv1",
					},
				},
			},
		},
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
		System: 0,
		BindID: 0,
		Mark:   "",
		Search: "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetMark(t *testing.T) {
	data, err := GetMark(&ArgsGetMark{
		Mark:   "test",
		OrgID:  0,
		UserID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestReturnOrg(t *testing.T) {
	err := Set(&ArgsSet{
		System: 2,
		Mark:   "test",
		BindID: 123,
		Data: FieldsChildList{
			{
				Mark: "a1o",
				Name: "v1o",
				Params: CoreSQLConfig.FieldsConfigsType{
					{
						Mark: "p1o",
						Val:  "pv1o",
					},
				},
			},
		},
	})
	ToolsTest.ReportError(t, err)
	err = ReturnOrg(&ArgsReturnOrg{
		Mark:  "test",
		OrgID: 123,
	})
	data, err := GetMark(&ArgsGetMark{
		Mark:   "test",
		OrgID:  123,
		UserID: 0,
	})
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		t.Error(err)
	} else {
		if data.System != 0 {
			t.Error("system not 0")
		}
		if data.BindID > 0 {
			t.Error("bind than 0")
		}
	}
}

func TestReturnUser(t *testing.T) {
	err := Set(&ArgsSet{
		System: 1,
		Mark:   "test",
		BindID: 234,
		Data: FieldsChildList{
			{
				Mark: "a1u",
				Name: "v1u",
				Params: CoreSQLConfig.FieldsConfigsType{
					{
						Mark: "p1u",
						Val:  "pv1u",
					},
				},
			},
		},
	})
	ToolsTest.ReportError(t, err)
	err = ReturnUser(&ArgsReturnUser{
		Mark:   "test",
		UserID: 234,
	})
	data, err := GetMark(&ArgsGetMark{
		Mark:   "test",
		OrgID:  0,
		UserID: 234,
	})
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		t.Error(err)
	} else {
		if data.System != 0 {
			t.Error("system not 0")
		}
		if data.BindID > 0 {
			t.Error("bind than 0")
		}
	}
}

func TestDeleteMark(t *testing.T) {
	err := DeleteMark(&ArgsDeleteMark{
		Mark: "test",
	})
	ToolsTest.ReportError(t, err)
}
