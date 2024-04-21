package ERPProduct

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newBrandData FieldsBrand
)

func TestBrandInit(t *testing.T) {
	TestProductInit(t)
	TestSetProduct(t)
}

func TestCreateBrand(t *testing.T) {
	newBrandDataID, err := CreateBrand(&ArgsCreateBrand{
		OrgID:      newERPProductData.OrgID,
		Code:       CoreFilter.GetRandStr4(10),
		Name:       "测试品牌",
		CategoryID: 0,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newBrandData.ID = newBrandDataID
	t.Log("new brand id: ", newBrandDataID)
}

func TestGetBrand(t *testing.T) {
	var err error
	newBrandData = GetBrand(newBrandData.ID, newERPProductData.OrgID)
	if err != nil {
		t.Fatal(err)
		return
	}
	ToolsTest.ReportData(t, err, newBrandData)
}

func TestGetBrandByCode(t *testing.T) {
	data := GetBrandByCode(newBrandData.Code, newBrandData.OrgID)
	if data.ID < 1 {
		t.Fatal("GetBrandByCode fail")
		return
	}
	ToolsTest.ReportData(t, nil, data)
}

func TestGetBrandList(t *testing.T) {
	dataList, dataCount, err := GetBrandList(&ArgsGetBrandList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    newBrandData.OrgID,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateBrand(t *testing.T) {
	err := UpdateBrand(&ArgsUpdateBrand{
		ID:    newBrandData.ID,
		OrgID: newBrandData.OrgID,
		Name:  fmt.Sprint(newBrandData.Name, "Update"),
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteBrand(t *testing.T) {
	err := DeleteBrand(&ArgsDeleteBrand{
		ID:    newBrandData.ID,
		OrgID: newBrandData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestBrandClear(t *testing.T) {
	TestDeleteProduct(t)
	TestProductClear(t)
}
