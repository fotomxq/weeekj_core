package ERPProductMall

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ERPProduct "github.com/fotomxq/weeekj_core/v5/erp/product"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newERPProductData  ERPProduct.FieldsProduct
	newProductMallData FieldsProductMall
)

func TestProductMallInit(t *testing.T) {
	TestInit(t)
	var errCode string
	var err error
	newERPProductData, errCode, err = ERPProduct.SetProduct2(&ERPProduct.ArgsSetProduct{
		OrgID:            TestOrg.OrgData.ID,
		CompanyID:        0,
		CompanyName:      "测试供应商A",
		SortID:           0,
		Tags:             []int64{},
		SN:               CoreFilter.GetRandStr4(10),
		Code:             CoreFilter.GetRandStr4(10),
		PinYin:           "ceshigongyingshanga",
		EnName:           "ceshigongyingshanga",
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
		RebatePrice:      ERPProduct.FieldsProductRebateList{},
		Params:           CoreSQLConfig.FieldsConfigsType{},
		SyncMallCore:     false,
	})
	if err != nil {
		t.Fatal(err, "...", errCode)
		return
	}
	ToolsTest.ReportData(t, err, newERPProductData)
}

func TestCreateProductMall(t *testing.T) {
	dataID, err := CreateProductMall(&ArgsCreateProductMall{
		OrgID:       TestOrg.OrgData.ID,
		ProductID:   newERPProductData.ID,
		ProductName: newERPProductData.Title,
		Price:       10,
		CategoryID:  0,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newProductMallData = getProductMallData(dataID)
}

func TestGetProductMall(t *testing.T) {
	data := GetProductMall(newProductMallData.ID, TestOrg.OrgData.ID)
	ToolsTest.ReportData(t, nil, data)
}

func TestGetProductMallList(t *testing.T) {
	dataList, dataCount, err := GetProductMallList(&ArgsGetProductMallList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      -1,
		ProductID:  -1,
		CategoryID: -1,
		IsRemove:   false,
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetProductMallByProductID(t *testing.T) {
	data := GetProductMallByProductID(newERPProductData.ID)
	ToolsTest.ReportData(t, nil, data)
}

func TestUpdateProductMall(t *testing.T) {
	err := UpdateProductMall(&ArgsUpdateProductMall{
		ID:          newProductMallData.ID,
		OrgID:       newProductMallData.OrgID,
		ProductID:   newProductMallData.ProductID,
		ProductName: newProductMallData.ProductName,
		Price:       newProductMallData.Price,
		CategoryID:  newProductMallData.CategoryID,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteProductMall(t *testing.T) {
	err := DeleteProductMall(&ArgsDeleteProductMall{
		ID:    newProductMallData.ID,
		OrgID: newProductMallData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestProductMallClear(t *testing.T) {
	err := ERPProduct.DeleteProduct(&ERPProduct.ArgsDeleteProduct{
		ID:    newERPProductData.ID,
		OrgID: newERPProductData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestClear(t)
}
