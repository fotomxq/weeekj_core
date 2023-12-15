package BasePageStyle

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newPageData FieldsPage
)

func TestInitPage(t *testing.T) {
	TestInitTemplate(t)
	TestCreateTemplate(t)
	TestGetTemplateList(t)
}

func TestSetPage(t *testing.T) {
	err := SetPage(&ArgsSetPage{
		OrgID:  TestOrg.OrgData.ID,
		System: "test_system",
		Page:   "test_page",
		Title:  "test_title",
		Data:   "test_data",
		ComponentList: []FieldsPageComponent{
			{
				ComponentID:   newComponentData.ID,
				ComponentMark: newComponentData.Mark,
				Data:          "test_data_component",
				Params:        CoreSQLConfig.FieldsConfigsType{},
			},
		},
		Params: CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetPageList(t *testing.T) {
	dataList, dataCount, err := GetPageList(&ArgsGetPageList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		System:   "",
		Page:     "",
		SortID:   -1,
		Tags:     []int64{},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newPageData = dataList[0]
	}
}

func TestGetPageIDs(t *testing.T) {
	data, err := GetPageIDs(&ArgsGetPageIDs{
		IDs:        []int64{newPageData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetPageMark(t *testing.T) {
	data, err := GetPageMark(&ArgsGetPageMark{
		OrgID:  newPageData.OrgID,
		System: newPageData.System,
		Page:   newPageData.Page,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeletePage(t *testing.T) {
	err := DeletePage(&ArgsDeletePage{
		ID:    newPageData.ID,
		OrgID: newPageData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearPage(t *testing.T) {
	TestDeleteTemplate(t)
	TestClearTemplate(t)
}
