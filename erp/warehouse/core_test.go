package ERPWarehouse

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	ERPProduct "gitee.com/weeekj/weeekj_core/v5/erp/product"
	ServiceCompany "gitee.com/weeekj/weeekj_core/v5/service/company"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit            = false
	newERPProductData ERPProduct.FieldsProduct
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
		TestOrg.LocalInit()
		ServiceCompany.Init()
		ERPProduct.Init()
		Init()
	}
	isInit = true
	if TestOrg.OrgData.ID < 1 {
		TestOrg.LocalCreateOrg(t)
	}
}

func TestCreateProduct(t *testing.T) {
	if newERPProductData.ID > 0 {
		return
	}
	var errCode string
	var err error
	newERPProductData, errCode, err = ERPProduct.SetProduct(&ERPProduct.ArgsSetProduct{
		OrgID:            newWarehouseData.OrgID,
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
		t.Error("code: ", errCode, ", err: ", err)
		return
	}
	t.Log("new product data: ", newERPProductData.ID)
}

func TestClear(t *testing.T) {
	if newERPProductData.ID > 0 {
		err := ERPProduct.DeleteProduct(&ERPProduct.ArgsDeleteProduct{
			ID:    newERPProductData.ID,
			OrgID: newERPProductData.OrgID,
		})
		if err != nil {
			t.Error(err)
		} else {
			newERPProductData = ERPProduct.FieldsProduct{}
		}
	}
	TestOrg.LocalClear(t)
}
