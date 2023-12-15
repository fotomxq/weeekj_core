package BasePageStyle

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newTemplateData FieldsTemplate
)

func TestInitTemplate(t *testing.T) {
	TestInit(t)
	TestCreateComponent(t)
	TestGetComponentList(t)
}

func TestCreateTemplate(t *testing.T) {
	err := CreateTemplate(&ArgsCreateTemplate{
		System:              "test_system",
		Page:                "test_page",
		Name:                "test_name",
		Des:                 "test_des",
		CoverFileID:         0,
		SortID:              0,
		Tags:                []int64{},
		OrgSubConfigID:      []int64{},
		OrgFuncList:         []string{},
		ComponentIDs:        []int64{newTemplateData.ID},
		DefaultComponentIDs: []int64{newTemplateData.ID},
		Data:                "test_data",
		DefaultData:         "test_default_data",
		Params:              CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetTemplateList(t *testing.T) {
	dataList, dataCount, err := GetTemplateList(&ArgsGetTemplateList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		System:   "",
		Page:     "",
		SortID:   -1,
		Tags:     []int64{},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newTemplateData = dataList[0]
	}
}

func TestGetTemplateID(t *testing.T) {
	data, err := GetTemplateID(&ArgsGetTemplateID{
		ID: newTemplateData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateTemplate(t *testing.T) {
	err := UpdateTemplate(&ArgsUpdateTemplate{
		ID:                  newTemplateData.ID,
		System:              newTemplateData.System,
		Page:                newTemplateData.Page,
		Name:                newTemplateData.Name,
		Des:                 newTemplateData.Des,
		CoverFileID:         newTemplateData.CoverFileID,
		SortID:              newTemplateData.SortID,
		Tags:                newTemplateData.Tags,
		OrgSubConfigID:      newTemplateData.OrgSubConfigID,
		OrgFuncList:         newTemplateData.OrgFuncList,
		ComponentIDs:        newTemplateData.ComponentIDs,
		DefaultComponentIDs: newTemplateData.DefaultComponentIDs,
		Data:                newTemplateData.Data,
		DefaultData:         newTemplateData.DefaultData,
		Params:              newTemplateData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestOutputTemplate(t *testing.T) {
	data, err := OutputTemplate()
	ToolsTest.ReportData(t, err, data)
}

func TestImportTemplate(t *testing.T) {
	data, err := OutputTemplate()
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		return
	}
	err = ImportTemplate(&ArgsImportTemplate{
		Data: data,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteTemplate(t *testing.T) {
	err := DeleteTemplate(&ArgsDeleteTemplate{
		ID: newTemplateData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearTemplate(t *testing.T) {
	TestDeleteComponent(t)
	TestClear(t)
}
