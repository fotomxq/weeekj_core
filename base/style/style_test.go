package BaseStyle

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newStyleData FieldsStyle
)

func TestInitStyle(t *testing.T) {
	TestInit(t)
	TestCreateComponent(t)
}

func TestCreateStyle(t *testing.T) {
	var err error
	newStyleData, err = CreateStyle(&ArgsCreateStyle{
		Name:        "测试名称",
		Mark:        "test01",
		SystemMark:  "app",
		Components:  []int64{newComponentData.ID},
		Title:       "测试标题",
		Des:         "测试描述",
		CoverFileID: 0,
		DesFiles:    []int64{},
		SortID:      0,
		Tags:        []int64{},
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newStyleData)
}

func TestGetStyleList(t *testing.T) {
	dataList, dataCount, err := GetStyleList(&ArgsGetStyleList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		Mark:     "",
		Sort:     0,
		Tags:     nil,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetStyleMore(t *testing.T) {
	data, err := GetStyleMore(&ArgsGetStyleMore{
		IDs:        []int64{newStyleData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetStyleMoreMap(t *testing.T) {
	data, err := GetStyleMoreMap(&ArgsGetStyleMore{
		IDs:        []int64{newStyleData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetStyleByID(t *testing.T) {
	data, err := GetStyleByID(&ArgsGetStyleByID{
		ID: newStyleData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetStyleByMark(t *testing.T) {
	data, err := GetStyleByMark(&ArgsGetStyleByMark{
		Mark: newStyleData.Mark,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateStyle(t *testing.T) {
	err := UpdateStyle(&ArgsUpdateStyle{
		ID:          newStyleData.ID,
		Name:        newStyleData.Name,
		Mark:        newStyleData.Mark,
		SystemMark:  newStyleData.SystemMark,
		Components:  newStyleData.Components,
		Title:       newStyleData.Title,
		Des:         newStyleData.Des,
		CoverFileID: newStyleData.CoverFileID,
		DesFiles:    newStyleData.DesFiles,
		SortID:      newStyleData.SortID,
		Tags:        newStyleData.Tags,
		Params:      newStyleData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteStyle(t *testing.T) {
	err := DeleteStyle(&ArgsDeleteStyle{
		ID: newStyleData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearStyle(t *testing.T) {
	TestDeleteComponent(t)
	TestClearComponent(t)
}
