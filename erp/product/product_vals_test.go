package ERPProduct

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newProductValsSlotList []DataProductVal
)

func TestProductValsInit(t *testing.T) {
	TestTemplateBindInit(t)
	TestCreateTemplateBind(t)
	TestGetTemplateBindData(t)
}

func TestGetValsByBrandOrCategoryID(t *testing.T) {
	templateID, themeID, bpmSlotList, errCode, err := GetValsByBrandOrCategoryID(&ArgsGetValsByBrandOrCategoryID{
		OrgID:      TestOrg.OrgData.ID,
		BrandID:    newBrandData.ID,
		CategoryID: newSortData.ID,
	})
	if err != nil {
		t.Error(err, "...", errCode)
		return
	}
	t.Log(templateID, ",", themeID, ",", bpmSlotList)
}

func TestGetProductValsAndDefault(t *testing.T) {
	dataList, errCode, err := GetProductValsAndDefault(&ArgsGetProductValsAndDefault{
		OrgID:        TestOrg.OrgData.ID,
		ProductID:    newERPProductData.ID,
		CompanyID:    0,
		HaveMoreData: false,
	})
	if err != nil {
		t.Error(err, "...", errCode)
		return
	}
	t.Log(dataList)
	newProductValsSlotList = dataList
}

func TestGetProductValsTemplateID(t *testing.T) {
	templateBindData, errCode, err := GetProductValsTemplateID(&ArgsGetProductValsTemplateID{
		OrgID:     TestOrg.OrgData.ID,
		ProductID: newERPProductData.ID,
		CompanyID: 0,
	})
	if err != nil {
		t.Error(err, "...", errCode)
		return
	}
	t.Log(templateBindData)
}

func TestSetProductVals(t *testing.T) {
	if len(newProductValsSlotList) != 2 {
		t.Error("newProductValsSlotList length error")
		return
	}
	errCode, err := SetProductVals(&ArgsSetProductVals{
		OrgID:     TestOrg.OrgData.ID,
		ProductID: newERPProductData.ID,
		CompanyID: 0,
		Vals: []DataProductVal{
			{
				OrderNum:     1,
				SlotID:       newBPMSlotData1.ID,
				DataValue:    newBPMSlotData1.DefaultValue,
				DataValueNum: newProductValsSlotList[0].DataValueNum,
				DataValueInt: newProductValsSlotList[0].DataValueInt,
				Params:       newProductValsSlotList[0].Params,
			},
			{
				OrderNum:     2,
				SlotID:       newBPMSlotData2.ID,
				DataValue:    newBPMSlotData2.DefaultValue,
				DataValueNum: newProductValsSlotList[1].DataValueNum,
				DataValueInt: newProductValsSlotList[1].DataValueInt,
				Params:       newProductValsSlotList[1].Params,
			},
		},
		HaveMoreData: false,
	})
	if err != nil {
		t.Error(err, "...", errCode)
		return
	}
}

func TestGetProductVals(t *testing.T) {
	data, err := GetProductVals(&ArgsGetProductVals{
		OrgID:     TestOrg.OrgData.ID,
		ProductID: newERPProductData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestClearProductVals(t *testing.T) {
	err := ClearProductVals(&ArgsClearProductVals{
		OrgID:     TestOrg.OrgData.ID,
		ProductID: newERPProductData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestProductValsClear(t *testing.T) {
	TestDeleteTemplateBind(t)
	TestTemplateBindClear(t)
}
