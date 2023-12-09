package UserFocus

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	userID   int64 = 123
	orgID    int64 = 234
	fromInfo       = CoreSQLFrom.FieldsFrom{
		System: "mall",
		ID:     345,
		Mark:   "",
		Name:   "测试内容",
	}
	newData FieldsFocus
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestCreate(t *testing.T) {
	var err error
	newData, err = Create(&ArgsCreate{
		UserID:   userID,
		OrgID:    orgID,
		FromInfo: fromInfo,
	})
	ToolsTest.ReportData(t, err, newData)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:   userID,
		OrgID:    0,
		FromInfo: fromInfo,
		IsRemove: false,
		Search:   "内容",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDelete(t *testing.T) {
	err := Delete(&ArgsDelete{
		ID:     newData.ID,
		UserID: userID,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteByFrom(t *testing.T) {
	TestCreate(t)
	err := DeleteByFrom(&ArgsDeleteByFrom{
		FromInfo: fromInfo,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteByOrg(t *testing.T) {
	TestCreate(t)
	err := DeleteByOrg(&ArgsDeleteByOrg{
		OrgID: orgID,
	})
	ToolsTest.ReportError(t, err)
}
