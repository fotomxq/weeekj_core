package FinanceReportForm

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newFileData FieldsFile
)

func TestInitFile(t *testing.T) {
	TestInit(t)
}

func TestCreateFile(t *testing.T) {
	TestCreateTemplate(t)
	data, err := CreateFile(&ArgsCreateFile{
		OrgID:      1,
		Name:       "测试文件",
		Des:        "test des",
		TemplateID: newTemplateData.ID,
		Params:     nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newFileData = data
		TestGetValByFile(t)
	}
}

func TestGetFileList(t *testing.T) {
	dataList, dataCount, err := GetFileList(&ArgsGetFileList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      -1,
		TemplateID: -1,
		IsRemove:   false,
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateFile(t *testing.T) {
	err := UpdateFile(&ArgsUpdateFile{
		ID:     newFileData.ID,
		OrgID:  newFileData.OrgID,
		Name:   newFileData.Name,
		Des:    newFileData.Des,
		Params: newFileData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteFile(t *testing.T) {
	err := DeleteFile(&ArgsDeleteFile{
		ID:    newFileData.ID,
		OrgID: newFileData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteTemplate(t)
}
