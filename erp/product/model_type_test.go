package ERPProduct

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newModelTypeData FieldsModelType
)

func TestModelTypeInit(t *testing.T) {
	TestProductInit(t)
	TestSetProduct(t)
}

func TestCreateModelType(t *testing.T) {
	newModelTypeDataID, err := CreateModelType(&ArgsCreateModelType{
		OrgID: newERPProductData.OrgID,
		Code:  CoreFilter.GetRandStr4(10),
		Name:  "测试品牌",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newModelTypeData.ID = newModelTypeDataID
	t.Log("new ModelType id: ", newModelTypeDataID)
}

func TestGetModelType(t *testing.T) {
	var err error
	newModelTypeData = GetModelType(newModelTypeData.ID, newERPProductData.OrgID)
	if err != nil {
		t.Fatal(err)
		return
	}
	ToolsTest.ReportData(t, err, newModelTypeData)
}

func TestGetModelTypeByCode(t *testing.T) {
	data := GetModelTypeByCode(newModelTypeData.Code, newModelTypeData.OrgID)
	if data.ID < 1 {
		t.Fatal("GetModelTypeByCode fail")
		return
	}
	ToolsTest.ReportData(t, nil, data)
}

func TestGetModelTypeList(t *testing.T) {
	dataList, dataCount, err := GetModelTypeList(&ArgsGetModelTypeList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    newModelTypeData.OrgID,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateModelType(t *testing.T) {
	err := UpdateModelType(&ArgsUpdateModelType{
		ID:    newModelTypeData.ID,
		OrgID: newModelTypeData.OrgID,
		Name:  fmt.Sprint(newModelTypeData.Name, "Update"),
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteModelType(t *testing.T) {
	err := DeleteModelType(&ArgsDeleteModelType{
		ID:    newModelTypeData.ID,
		OrgID: newModelTypeData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestModelTypeClear(t *testing.T) {
	TestDeleteProduct(t)
	TestProductClear(t)
}
