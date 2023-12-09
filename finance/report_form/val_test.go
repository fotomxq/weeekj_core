package FinanceReportForm

import (
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newValData FieldsVal
)

func TestInitVal(t *testing.T) {
	TestInit(t)
}

func TestSetVal(t *testing.T) {
	TestCreateFile(t)
	data, err := SetVal(&ArgsSetVal{
		OrgID:    1,
		FileID:   newFileData.ID,
		ColID:    newColData.ID,
		Mark:     "A1",
		Val:      "test val",
		ValFloat: 2.0,
		ValInt:   1,
		Params:   nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newValData = data
	}
}

func TestGetValByFile(t *testing.T) {
	data, err := GetValByFile(&ArgsGetValByFile{
		FileID: newFileData.ID,
		OrgID:  newFileData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteVal(t *testing.T) {
	err := DeleteVal(&ArgsDeleteVal{
		ID:    newValData.ID,
		OrgID: newValData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteFile(t)
}
