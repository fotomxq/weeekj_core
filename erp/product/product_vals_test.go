package ERPProduct

import (
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

func TestProductValsInit(t *testing.T) {
	TestTemplateBindInit(t)
	TestCreateTemplateBind(t)
	TestGetTemplateBindData(t)
}

func TestGetValsByBrandOrCategoryID(t *testing.T) {

}

func TestGetProductValsAndDefault(t *testing.T) {

}

func TestGetProductValsTemplateID(t *testing.T) {

}

func TestSetProductVals(t *testing.T) {

}

func TestGetProductVals(t *testing.T) {

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
