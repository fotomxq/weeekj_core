package EAMCore

import (
	BaseBPM "github.com/fotomxq/weeekj_core/v5/base/bpm"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ERPProduct "github.com/fotomxq/weeekj_core/v5/erp/product"
	ServiceCompany "github.com/fotomxq/weeekj_core/v5/service/company"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit            = false
	newSortData       ClassSort.FieldsSort
	newERPProductData ERPProduct.FieldsProduct
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
		TestOrg.LocalInit()
		ServiceCompany.Init()
		BaseBPM.Init()
	}
	isInit = true
	Init()
	var err error
	//创建分类
	newSortData, err = ERPProduct.Sort.Create(&ClassSort.ArgsCreate{
		BindID:      TestOrg.OrgData.ID,
		Mark:        CoreFilter.GetRandStr4(10),
		ParentID:    0,
		CoverFileID: 0,
		DesFiles:    []int64{},
		Name:        "测试分类",
		Des:         "测试分类描述",
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newSortData)
}
func TestSetProduct(t *testing.T) {
	var errCode string
	var err error
	newERPProductData, errCode, err = ERPProduct.SetProduct(&ERPProduct.ArgsSetProduct{
		OrgID:            TestOrg.OrgData.ID,
		CompanyID:        0,
		CompanyName:      "测试供应商A",
		SortID:           newSortData.ID,
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

func TestDeleteProduct(t *testing.T) {
	err := ERPProduct.DeleteProduct(&ERPProduct.ArgsDeleteProduct{
		ID:    newERPProductData.ID,
		OrgID: newERPProductData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClear(t *testing.T) {
	//删除分类
	err := ERPProduct.Sort.DeleteByID(&ClassSort.ArgsDeleteByID{
		ID:     newSortData.ID,
		BindID: TestOrg.OrgData.ID,
	})
	//清理程序
	ToolsTest.ReportError(t, err)
	TestOrg.LocalClear(t)
}
