package ERPProduct

import (
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newTemplateBindData FieldsTemplateBind
)

func TestTemplateBindInit(t *testing.T) {
	TestTemplateInit(t)
	TestCreateBPMThemeSlotData(t)
	TestCreateTemplate(t)
}

func TestCreateTemplateBind(t *testing.T) {
	newTemplateBindDataID, err := CreateTemplateBind(&ArgsCreateTemplateBind{
		OrgID:      TestOrg.OrgData.ID,
		TemplateID: newTemplateData.ID,
		CategoryID: newSortData.ID,
		BrandID:    newBrandData.ID,
	})
	ToolsTest.ReportData(t, err, newTemplateBindDataID)
	newTemplateBindData.ID = newTemplateBindDataID
}

func TestGetTemplateBindData(t *testing.T) {
	newTemplateBindData = GetTemplateBindData(&ArgsGetTemplateBindData{
		OrgID:      TestOrg.OrgData.ID,
		TemplateID: newTemplateData.ID,
		CategoryID: newSortData.ID,
		BrandID:    newBrandData.ID,
	})
	if newTemplateBindData.ID < 1 {
		t.Fatal("get template bind data fail")
		return
	}
	ToolsTest.ReportData(t, nil, newTemplateBindData)
}

func TestGetTemplateBindList(t *testing.T) {
	dataList, dataCount, err := GetTemplateBindList(&ArgsGetTemplateBindList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      TestOrg.OrgData.ID,
		TemplateID: -1,
		CategoryID: -1,
		BrandID:    -1,
		IsRemove:   false,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestCheckTemplateBind(t *testing.T) {
	b := CheckTemplateBind(&ArgsCheckTemplateBind{
		OrgID:      TestOrg.OrgData.ID,
		TemplateID: newTemplateData.ID,
		CategoryID: newSortData.ID,
		BrandID:    newBrandData.ID,
	})
	if !b {
		t.Fatal("check template bind fail")
		return
	}
}

func TestDeleteTemplateBind(t *testing.T) {
	err := DeleteTemplateBind(&ArgsDeleteTemplateBind{
		OrgID:      TestOrg.OrgData.ID,
		TemplateID: newTemplateData.ID,
		CategoryID: newSortData.ID,
		BrandID:    newBrandData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestTemplateBindClear(t *testing.T) {
	TestDeleteTemplate(t)
	TestTemplateClear(t)
}
