package AnalysisUserVisit

import (
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

func TestVisitInit(t *testing.T) {
	TestInit(t)
}

func TestVisitCreate(t *testing.T) {
	err := VisitCreate(&ArgsVisitCreate{
		OrgID:      1,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
		UserID:     1,
		Country:    0,
		Phone:      "",
		IP:         "",
		Mark:       "",
		Action:     "",
		Params:     nil,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetVisitList(t *testing.T) {
	dataList, dataCount, err := GetVisitList(&ArgsGetVisitList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      -1,
		CreateInfo: CoreSQLFrom.FieldsFrom{},
		UserID:     -1,
		Country:    -1,
		Phone:      "",
		IP:         "",
		Action:     "",
		Mark:       "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDeleteVisitByUser(t *testing.T) {
	err := DeleteVisitByUser(&ArgsDeleteVisitByUser{
		UserID: 1,
	})
	ToolsTest.ReportError(t, err)
}
