package BaseStyle

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newComponentData FieldsComponent
)

func TestInitComponent(t *testing.T) {
	TestInit(t)
}

func TestCreateComponent(t *testing.T) {
	var err error
	newComponentData, err = CreateComponent(&ArgsCreateComponent{
		Mark:        "test01",
		Name:        "测试名称",
		Des:         "测试描述",
		CoverFileID: 1,
		DesFiles:    []int64{123},
		SortID:      2,
		Tags:        []int64{234},
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newComponentData)
}

func TestGetComponentList(t *testing.T) {
	dataList, dataCount, err := GetComponentList(&ArgsGetComponentList{
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

func TestGetComponentByID(t *testing.T) {
	data, err := GetComponentByID(&ArgsGetComponentByID{
		ID: newComponentData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetComponentByMark(t *testing.T) {
	data, err := GetComponentByMark(&ArgsGetComponentByMark{
		Mark: newComponentData.Mark,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetComponentMore(t *testing.T) {
	data, err := GetComponentMore(&ArgsGetComponentMore{
		IDs:        []int64{newComponentData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetComponentMoreMap(t *testing.T) {
	data, err := GetComponentMoreMap(&ArgsGetComponentMore{
		IDs:        []int64{newComponentData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateComponent(t *testing.T) {
	err := UpdateComponent(&ArgsUpdateComponent{
		ID:          newComponentData.ID,
		Mark:        newComponentData.Mark,
		Name:        newComponentData.Name,
		Des:         newComponentData.Des,
		CoverFileID: newComponentData.CoverFileID,
		DesFiles:    newComponentData.DesFiles,
		SortID:      newComponentData.SortID,
		Tags:        newComponentData.Tags,
		Params:      newComponentData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteComponent(t *testing.T) {
	err := DeleteComponent(&ArgsDeleteComponent{
		ID: newComponentData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearComponent(t *testing.T) {
	TestClear(t)
}
