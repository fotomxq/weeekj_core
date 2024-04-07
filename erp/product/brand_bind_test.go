package ERPProduct

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newBrandBindData FieldsBrandBind
)

func TestBrandBindInit(t *testing.T) {
	TestBrandInit(t)
	TestCreateBrand(t)
	TestGetBrand(t)
}

func TestCreateBrandBind(t *testing.T) {
	newBrandBindDataID, err := CreateBrandBind(&ArgsCreateBrandBind{
		OrgID:     newBrandData.OrgID,
		BrandID:   newBrandData.ID,
		CompanyID: 0,
		ProductID: newERPProductData.ID,
	})
	if err != nil {
		t.Fatal("TestCreateBrandBind: ", err)
		return
	}
	newBrandBindData.ID = newBrandBindDataID
	t.Log("new brand bind id: ", newBrandBindData.ID, ", newERPProductData id: ", newERPProductData.ID)
}

func TestGetBrandBindData(t *testing.T) {
	newBrandBindData = GetBrandBindData(&ArgsGetBrandBindData{
		OrgID:     newBrandData.OrgID,
		BrandID:   newBrandData.ID,
		CompanyID: 0,
		ProductID: newERPProductData.ID,
	})
	if newBrandBindData.ID < 1 {
		t.Fatal("GetBrandBindData fail")
		return
	}
	ToolsTest.ReportData(t, nil, newBrandBindData)
}

func TestGetBrandBindList(t *testing.T) {
	dataList, dataCount, err := GetBrandBindList(&ArgsGetBrandBindList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:     newBrandData.OrgID,
		BrandID:   -1,
		CompanyID: -1,
		ProductID: -1,
		IsRemove:  false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestCheckBrandBind(t *testing.T) {
	b := CheckBrandBind(&ArgsCheckBrandBind{
		OrgID:     newBrandData.OrgID,
		BrandID:   newBrandData.ID,
		CompanyID: 0,
		ProductID: newERPProductData.ID,
	})
	if !b {
		t.Fatal("CheckBrandBind fail")
		return
	}
}

func TestDeleteBrandBind(t *testing.T) {
	err := DeleteBrandBind(&ArgsDeleteBrandBind{
		OrgID:     newBrandData.OrgID,
		BrandID:   newBrandData.ID,
		CompanyID: 0,
		ProductID: newERPProductData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestBrandBindClear(t *testing.T) {
	TestDeleteBrand(t)
	TestBrandClear(t)
}
