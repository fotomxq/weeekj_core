package BasePageStyle

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newComponentData FieldsComponent
)

func TestInitComponent(t *testing.T) {
	TestInit(t)
}

func TestCreateComponent(t *testing.T) {
	err := CreateComponent(&ArgsCreateComponent{
		System:         "test_system1",
		Mark:           "test_mark1",
		Name:           "test_name1",
		Des:            "test_des1",
		CoverFileID:    0,
		SortID:         0,
		Tags:           []int64{},
		OrgSubConfigID: []int64{},
		OrgFuncList:    []string{},
		Data:           "test_data1",
		Params:         CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetComponentList(t *testing.T) {
	dataList, dataCount, err := GetComponentList(&ArgsGetComponentList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		System:   "",
		Mark:     "",
		SortID:   -1,
		Tags:     []int64{},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newComponentData = dataList[0]
	}
}

func TestGetComponentIDs(t *testing.T) {
	data, err := GetComponentIDs(&ArgsGetComponentIDs{
		IDs:        []int64{newComponentData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetComponentMarks(t *testing.T) {
	data, err := GetComponentMarks(&ArgsGetComponentMarks{
		System:     newComponentData.System,
		Marks:      []string{newComponentData.Mark},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateComponent(t *testing.T) {
	err := UpdateComponent(&ArgsUpdateComponent{
		ID:             newComponentData.ID,
		Name:           newComponentData.Name,
		Des:            newComponentData.Data,
		CoverFileID:    newComponentData.CoverFileID,
		SortID:         newComponentData.SortID,
		Tags:           newComponentData.Tags,
		OrgSubConfigID: newComponentData.OrgSubConfigID,
		OrgFuncList:    newComponentData.OrgFuncList,
		Data:           newComponentData.Data,
		Params:         newComponentData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestOutputComponent(t *testing.T) {
	data, err := OutputComponent()
	ToolsTest.ReportData(t, err, data)
}

func TestImportComponent(t *testing.T) {
	data, err := OutputComponent()
	ToolsTest.ReportData(t, err, data)
	if err != nil {
		return
	}
	err = ImportComponent(&ArgsImportComponent{
		Data: data,
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
