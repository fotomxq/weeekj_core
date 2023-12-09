package MarketCore

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newBindData FieldsBind
)

func TestInitBind(t *testing.T) {
	TestInitConfig(t)
	TestCreateConfig(t)
}

func TestCreateBind(t *testing.T) {
	var err error
	newBindData, err = CreateBind(&ArgsCreateBind{
		SortID:     0,
		Tags:       []int64{},
		OrgID:      TestOrg.OrgData.ID,
		BindID:     TestOrg.BindData.ID,
		BindUserID: TestOrg.UserInfo.ID,
		Des:        "测试描述",
		Params:     nil,
	})
	ToolsTest.ReportData(t, err, newBindData)
}

func TestGetBindList(t *testing.T) {
	dataList, dataCount, err := GetBindList(&ArgsGetBindList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		SortID:     -1,
		Tags:       nil,
		OrgID:      TestOrg.OrgData.ID,
		BindID:     -1,
		BindUserID: -1,
		BindInfoID: -1,
		FromInfo:   CoreSQLFrom.FieldsFrom{},
		IsRemove:   false,
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetBindGroupList(t *testing.T) {
	dataList, dataCount, err := GetBindGroupList(&ArgsGetBindGroupList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "count",
			Desc: true,
		},
		SortID:   -1,
		Tags:     nil,
		OrgID:    -1,
		BindID:   -1,
		FromInfo: CoreSQLFrom.FieldsFrom{},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateBind(t *testing.T) {
	err := UpdateBind(&ArgsUpdateBind{
		ID:     newBindData.ID,
		SortID: newBindData.SortID,
		Tags:   newBindData.Tags,
		OrgID:  TestOrg.OrgData.ID,
		BindID: newBindData.BindID,
		Des:    newBindData.Des,
		Params: newBindData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteBind(t *testing.T) {
	err := DeleteBind(&ArgsDeleteBind{
		ID:    newBindData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearBind(t *testing.T) {
	TestDeleteConfig(t)
	TestClearConfig(t)
}
