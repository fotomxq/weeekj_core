package ERPProduct

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newERPProductData FieldsProduct
)

func TestProductInit(t *testing.T) {
	TestInit(t)
	TestSortCreate(t)
	TestCreateModelType(t)
	TestGetModelType(t)
}

func TestSetProduct(t *testing.T) {
	var errCode string
	var err error
	newERPProductData, errCode, err = SetProduct(&ArgsSetProduct{
		OrgID:            TestOrg.OrgData.ID,
		CompanyID:        0,
		CompanyName:      "测试供应商A",
		SortID:           newSortData.ID,
		Tags:             []int64{},
		SN:               CoreFilter.GetRandStr4(10),
		Code:             CoreFilter.GetRandStr4(10),
		PinYin:           "ceshigongyingshanga",
		EnName:           "ceshigongyingshanga",
		ModelTypeID:      newModelTypeData.ID,
		ManufacturerName: "测试生产厂商A",
		Title:            "产品测试名称",
		TitleDes:         "产品测试描述",
		Des:              "产品描述信息",
		CoverFileIDs:     []int64{},
		ExpireHour:       0,
		Weight:           0,
		SizeW:            0,
		SizeH:            0,
		SizeZ:            0,
		PackType:         0,
		PackUnitName:     "个",
		PackUnit:         1,
		TipPrice:         100,
		TipTaxPrice:      150,
		IsDiscount:       false,
		Currency:         86,
		CostPrice:        50,
		Tax:              60,
		TaxCostPrice:     70,
		RebatePrice:      FieldsProductRebateList{},
		Params:           CoreSQLConfig.FieldsConfigsType{},
		SyncMallCore:     false,
	})
	if err != nil {
		t.Fatal(err, "...", errCode)
		return
	}
	ToolsTest.ReportData(t, err, newERPProductData)
}

func TestSetProduct2(t *testing.T) {
	newData, errCode, err := SetProduct2(&ArgsSetProduct{
		OrgID:            TestOrg.OrgData.ID,
		CompanyID:        0,
		CompanyName:      "测试供应商A",
		SortID:           newSortData.ID,
		Tags:             []int64{},
		SN:               CoreFilter.GetRandStr4(10),
		Code:             CoreFilter.GetRandStr4(10),
		PinYin:           "ceshigongyingshanga",
		EnName:           "ceshigongyingshanga",
		ModelTypeID:      newModelTypeData.ID,
		ManufacturerName: "测试生产厂商A",
		Title:            "产品测试名称",
		TitleDes:         "产品测试描述",
		Des:              "产品描述信息",
		CoverFileIDs:     []int64{},
		ExpireHour:       0,
		Weight:           0,
		SizeW:            0,
		SizeH:            0,
		SizeZ:            0,
		PackType:         0,
		PackUnitName:     "个",
		PackUnit:         1,
		TipPrice:         100,
		TipTaxPrice:      150,
		IsDiscount:       false,
		Currency:         86,
		CostPrice:        50,
		Tax:              60,
		TaxCostPrice:     70,
		RebatePrice:      FieldsProductRebateList{},
		Params:           CoreSQLConfig.FieldsConfigsType{},
		SyncMallCore:     false,
	})
	if err != nil {
		t.Fatal(err, "...", errCode)
		return
	}
	ToolsTest.ReportData(t, err, newData)
}

func TestGetProductByID(t *testing.T) {
	data, err := GetProductByID(&ArgsGetProductByID{
		ID:    newERPProductData.ID,
		OrgID: newERPProductData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetProductByCode(t *testing.T) {
	data := GetProductBySN(newERPProductData.OrgID, newERPProductData.SN)
	if data.ID < 1 {
		t.Error("no data")
		return
	}
	ToolsTest.ReportData(t, nil, data)
}

func TestGetProductByIDNoErr(t *testing.T) {
	data := GetProductByIDNoErr(newERPProductData.ID)
	if data.ID < 1 {
		t.Error("no data")
		return
	}
	ToolsTest.ReportData(t, nil, data)
}

func TestGetProductBySN(t *testing.T) {
	data := GetProductBySN(newERPProductData.OrgID, newERPProductData.SN)
	if data.ID < 1 {
		t.Error("no data")
		return
	}
	ToolsTest.ReportData(t, nil, data)
}

func TestGetProductName(t *testing.T) {
	data := GetProductName(newERPProductData.ID)
	if data == "" {
		t.Error("no data")
		return
	}
	ToolsTest.ReportData(t, nil, data)
}

func TestGetProductMore(t *testing.T) {
	data := GetProductMore(&ArgsGetProductMore{
		IDs:   []int64{newERPProductData.ID},
		OrgID: newERPProductData.OrgID,
	})
	if len(data) < 1 {
		t.Error("no data")
		return
	}
	ToolsTest.ReportData(t, nil, data)
}

func TestGetProductList(t *testing.T) {
	dataList, dataCount, err := GetProductList(&ArgsGetProductList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      newERPProductData.OrgID,
		SortID:     -1,
		Tags:       nil,
		PackType:   -1,
		IsRemove:   false,
		SearchCode: "",
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDeleteProduct(t *testing.T) {
	err := DeleteProduct(&ArgsDeleteProduct{
		ID:    newERPProductData.ID,
		OrgID: newERPProductData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestProductClear(t *testing.T) {
	TestSortDelete(t)
	TestDeleteModelType(t)
	TestClear(t)
}
