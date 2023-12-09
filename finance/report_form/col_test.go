package FinanceReportForm

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newColData FieldsCol
)

func TestInitCol(t *testing.T) {
	TestInit(t)
}

func TestSetCol(t *testing.T) {
	data, err := SetCol(&ArgsSetCol{
		OrgID:   1,
		Mark:    "A1",
		Name:    "测试单元",
		Des:     "test des",
		ValType: 1,
		Params:  nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newColData = data
	}
}

func TestGetColList(t *testing.T) {
	dataList, dataCount, err := GetColList(&ArgsGetColList{
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

func TestDeleteCol(t *testing.T) {
	err := DeleteCol(&ArgsDeleteCol{
		ID:    newColData.ID,
		OrgID: newColData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
