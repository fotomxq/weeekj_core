package FinanceReportForm

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newTemplateData FieldsTemplate
)

func TestInitTemplate(t *testing.T) {
	TestInit(t)
}

func TestCreateTemplate(t *testing.T) {
	TestSetCol(t)
	data, err := CreateTemplate(&ArgsCreateTemplate{
		OrgID:      1,
		Name:       "测试模版",
		Des:        "test des",
		CoverFiles: []int64{},
		ColIDs:     []int64{newColData.ID},
		Params:     nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newTemplateData = data
	}
}

func TestGetTemplateList(t *testing.T) {
	dataList, dataCount, err := GetTemplateList(&ArgsGetTemplateList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateTemplate(t *testing.T) {
	err := UpdateTemplate(&ArgsUpdateTemplate{
		ID:         newTemplateData.ID,
		OrgID:      newTemplateData.OrgID,
		Name:       newTemplateData.Name,
		Des:        newTemplateData.Des,
		CoverFiles: newTemplateData.CoverFiles,
		ColIDs:     newTemplateData.ColIDs,
		Params:     newTemplateData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteTemplate(t *testing.T) {
	err := DeleteTemplate(&ArgsDeleteTemplate{
		ID:    newTemplateData.ID,
		OrgID: newTemplateData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteCol(t)
}
