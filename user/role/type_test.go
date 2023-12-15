package UserRole

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newTypeData FieldsType
)

func TestInitType(t *testing.T) {
	TestInit(t)
}

func TestCreateType(t *testing.T) {
	var err error
	newTypeData, err = CreateType(&ArgsCreateType{
		Mark:     "test",
		Name:     "测试",
		GroupIDs: []int64{},
		Params:   CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newTypeData)
}

func TestGetTypeID(t *testing.T) {
	data, err := GetTypeID(&ArgsGetTypeID{
		ID: newTypeData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetTypeList(t *testing.T) {
	dataList, dataCount, err := GetTypeList(&ArgsGetTypeList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetTypeMark(t *testing.T) {
	data, err := GetTypeMark(&ArgsGetTypeMark{
		Mark: newTypeData.Mark,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateType(t *testing.T) {
	err := UpdateType(&ArgsUpdateType{
		ID:       newTypeData.ID,
		Mark:     newTypeData.Mark,
		Name:     newTypeData.Name,
		GroupIDs: newTypeData.GroupIDs,
		Params:   newTypeData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteType(t *testing.T) {
	err := DeleteType(&ArgsDeleteType{
		ID: newTypeData.ID,
	})
	ToolsTest.ReportError(t, err)
}
