package FinanceAssets

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newProductData FieldsProduct
)

func TestInit2(t *testing.T) {
	TestInit(t)
}

// 创建产品
func TestCreateProduct(t *testing.T) {
	var err error
	newProductData, err = CreateProduct(&ArgsCreateProduct{
		OrgID:               orgID,
		Name:                "测试商品名称",
		Des:                 "测试商品",
		CoverFiles:          []int64{},
		DesFiles:            []int64{},
		Code:                "123",
		Currency:            86,
		Price:               300,
		WarehouseProductIDs: []int64{},
		MallCommodityIDs:    []int64{},
		Params:              CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(newProductData)
	}
}

func TestGetProductList(t *testing.T) {
	data, count, err := GetProductList(&ArgsGetProductList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    0,
		Code:     "",
		IsRemove: false,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
	}
}

func TestGetProductByID(t *testing.T) {
	var err error
	newProductData, err = GetProductByID(&ArgsGetProductByID{
		ID:    newProductData.ID,
		OrgID: orgID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(newProductData)
	}
}

func TestGetProductsName(t *testing.T) {
	bindList, err := GetProductsName(&ArgsGetProducts{
		IDs: []int64{newProductData.ID},
	})
	ToolsTest.ReportData(t, err, bindList)
}

func TestUpdateProduct(t *testing.T) {
	var err error
	err = UpdateProduct(&ArgsUpdateProduct{
		ID:                  newProductData.ID,
		OrgID:               newProductData.OrgID,
		Name:                newProductData.Name,
		Des:                 newProductData.Des,
		CoverFiles:          []int64{},
		DesFiles:            []int64{},
		Code:                newProductData.Code,
		Currency:            86,
		Price:               newProductData.Price,
		WarehouseProductIDs: []int64{},
		MallCommodityIDs:    []int64{},
		Params:              CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteProduct(t *testing.T) {
	var err error
	err = DeleteProduct(&ArgsDeleteProduct{
		ID:    newProductData.ID,
		OrgID: newProductData.OrgID,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClear2(t *testing.T) {
	TestClear(t)
}
