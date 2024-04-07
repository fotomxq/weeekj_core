package ERPProduct

import (
	"fmt"
	BaseBPM "github.com/fotomxq/weeekj_core/v5/base/bpm"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newBPMThemeCategoryData BaseBPM.FieldsThemeCategory
	newBPMThemeData         BaseBPM.FieldsTheme
	newBPMSlotData1         BaseBPM.FieldsSlot
	newBPMSlotData2         BaseBPM.FieldsSlot
	newTemplateData         FieldsTemplate
)

func TestTemplateInit(t *testing.T) {
	TestBrandBindInit(t)
	TestCreateBrandBind(t)
	TestGetBrandBindData(t)
}
func TestCreateBPMThemeSlotData(t *testing.T) {
	newBPMThemeCategoryDataID, err := BaseBPM.CreateThemeCategory(&BaseBPM.ArgsCreateThemeCategory{
		Name:        "测试主题域",
		Description: "测试主题域描述",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newBPMThemeCategoryData, err = BaseBPM.GetThemeCategoryByID(&BaseBPM.ArgsGetThemeCategoryByID{
		ID: newBPMThemeCategoryDataID,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newBPMThemeDataID, err := BaseBPM.CreateTheme(&BaseBPM.ArgsCreateTheme{
		CategoryID:  newBPMThemeCategoryData.ID,
		Name:        "测试主题",
		Description: "测试主题描述",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newBPMThemeData, err = BaseBPM.GetThemeByID(&BaseBPM.ArgsGetThemeByID{
		ID: newBPMThemeDataID,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newBPMSlotData1ID, err := BaseBPM.CreateSlot(&BaseBPM.ArgsCreateSlot{
		Name:            "测试插槽1",
		ThemeCategoryID: newBPMThemeCategoryData.ID,
		ThemeID:         newBPMThemeData.ID,
		ValueType:       "input",
		DefaultValue:    "default",
		Params:          "",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newBPMSlotData1, err = BaseBPM.GetSlotByID(&BaseBPM.ArgsGetSlotByID{
		ID: newBPMSlotData1ID,
	})
	newBPMSlotData2ID, err := BaseBPM.CreateSlot(&BaseBPM.ArgsCreateSlot{
		Name:            "测试插槽2",
		ThemeCategoryID: newBPMThemeCategoryData.ID,
		ThemeID:         newBPMThemeData.ID,
		ValueType:       "text",
		DefaultValue:    "default2text",
		Params:          "",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	newBPMSlotData2, err = BaseBPM.GetSlotByID(&BaseBPM.ArgsGetSlotByID{
		ID: newBPMSlotData2ID,
	})
}

func TestCreateTemplate(t *testing.T) {
	newTemplateID, err := CreateTemplate(&ArgsCreateTemplate{
		OrgID:      TestOrg.OrgData.ID,
		Name:       "测试模板",
		BPMThemeID: newBPMThemeData.ID,
	})
	if err != nil {
		t.Fatal("TestCreateTemplate: ", err)
		return
	}
	newTemplateData.ID = newTemplateID
}

func TestGetTemplate(t *testing.T) {
	newTemplateData = GetTemplate(newTemplateData.ID, TestOrg.OrgData.ID)
	if newTemplateData.ID < 1 {
		t.Fatal("GetTemplate Error")
		return
	}
}

func TestGetTemplateList(t *testing.T) {
	dataList, dataCount, err := GetTemplateList(&ArgsGetTemplateList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      TestOrg.OrgData.ID,
		BPMThemeID: -1,
		IsRemove:   false,
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateTemplate(t *testing.T) {
	err := UpdateTemplate(&ArgsUpdateTemplate{
		ID:         newTemplateData.ID,
		OrgID:      TestOrg.OrgData.ID,
		Name:       fmt.Sprint(newTemplateData.Name, "Update"),
		BPMThemeID: newTemplateData.BPMThemeID,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetTemplateBPMThemeSlotData(t *testing.T) {
	dataList, errCode, err := GetTemplateBPMThemeSlotData(newTemplateData.OrgID, newTemplateData.ID)
	if err != nil {
		t.Fatal(err, "...", errCode)
		return
	}
	ToolsTest.ReportData(t, nil, dataList)
}

func TestDeleteTemplate(t *testing.T) {
	err := DeleteTemplate(&ArgsDeleteTemplate{
		ID:    newTemplateData.ID,
		OrgID: newTemplateData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestTemplateClear(t *testing.T) {
	_ = BaseBPM.DeleteSlot(&BaseBPM.ArgsDeleteSlot{
		ID: newBPMSlotData1.ID,
	})
	_ = BaseBPM.DeleteSlot(&BaseBPM.ArgsDeleteSlot{
		ID: newBPMSlotData2.ID,
	})
	_ = BaseBPM.DeleteTheme(&BaseBPM.ArgsDeleteTheme{
		ID: newBPMThemeData.ID,
	})
	_ = BaseBPM.DeleteThemeCategory(&BaseBPM.ArgsDeleteThemeCategory{
		ID: newBPMThemeCategoryData.ID,
	})
	TestDeleteBrand(t)
	TestBrandClear(t)
}
