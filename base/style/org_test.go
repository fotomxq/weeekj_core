package BaseStyle

import (
	"testing"

	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

var (
	newOrg FieldsOrg
)

func TestInitOrg(t *testing.T) {
	TestInit(t)
	TestCreateComponent(t)
	TestCreateStyle(t)
}

func TestCreateOrgStyle(t *testing.T) {
	var err error
	newOrg, err = SetOrgStyle(&ArgsSetOrgStyle{
		OrgID:   123,
		StyleID: newStyleData.ID,
		Components: []FieldsOrgComponent{
			{
				ComponentID: newComponentData.ID,
				Params:      []CoreSQLConfig.FieldsConfigType{},
			},
		},
		Title:       "测试标题",
		Des:         "测试描述",
		CoverFileID: 0,
		DesFiles:    []int64{},
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newOrg)
}

func TestGetOrgList(t *testing.T) {
	dataList, dataCount, err := GetOrgList(&ArgsGetOrgList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    0,
		StyleID:  0,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetOrgByID(t *testing.T) {
	data, err := GetOrgByID(&ArgsGetOrgByID{
		ID:    newOrg.ID,
		OrgID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetOrgByStyleMark(t *testing.T) {
	data1, err := GetOrgByStyleMark(&ArgsGetOrgByStyleMark{
		Mark:  newStyleData.Mark,
		OrgID: newOrg.OrgID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("data1: ", data1)
	}
}

func TestGetOrgByStyleID(t *testing.T) {
	data1, err := GetOrgByStyleID(&ArgsGetOrgByStyleID{
		StyleID: newStyleData.ID,
		OrgID:   newOrg.OrgID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("data1: ", data1)
	}
}

func TestDeleteOrgStyle(t *testing.T) {
	err := DeleteOrgStyle(&ArgsDeleteOrgStyle{
		ID:    newOrg.ID,
		OrgID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearOrg(t *testing.T) {
	TestDeleteStyle(t)
	TestClearStyle(t)
}
