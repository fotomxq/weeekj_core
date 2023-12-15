package UserSubscription

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newConfigData FieldsConfig
)

func TestInitConfig(t *testing.T) {
	TestInit(t)
}

func TestCreateConfig(t *testing.T) {
	var err error
	newConfigData, err = CreateConfig(&ArgsCreateConfig{
		OrgID:             TestOrg.OrgData.ID,
		Mark:              "test_mark",
		TimeType:          0,
		TimeN:             1,
		Currency:          86,
		Price:             30,
		PriceOld:          56,
		Title:             "测试标题",
		Des:               "测试描述",
		CoverFileID:       0,
		DesFiles:          []int64{},
		UserGroups:        []int64{},
		ExemptionPrice:    15,
		ExemptionDiscount: 16,
		ExemptionMinPrice: 30,
		Limits:            nil,
		ExemptionTime:     nil,
		StyleID:           0,
		Params:            nil,
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestGetConfigByID(t *testing.T) {
	var err error
	newConfigData, err = GetConfigByID(&ArgsGetConfigByID{
		ID:    newConfigData.ID,
		OrgID: newConfigData.OrgID,
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestGetConfigList(t *testing.T) {
	dataList, dataCount, err := GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetConfigMore(t *testing.T) {
	data, err := GetConfigMore(&ArgsGetConfigMore{
		IDs:        []int64{newConfigData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetConfigMoreMap(t *testing.T) {
	data, err := GetConfigMoreMap(&ArgsGetConfigMore{
		IDs:        []int64{newConfigData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateConfig(t *testing.T) {
	err := UpdateConfig(&ArgsUpdateConfig{
		ID:                newConfigData.ID,
		OrgID:             newConfigData.OrgID,
		Mark:              newConfigData.Mark,
		TimeType:          newConfigData.TimeType,
		TimeN:             newConfigData.TimeN,
		Currency:          newConfigData.Currency,
		Price:             newConfigData.Price,
		PriceOld:          newConfigData.PriceOld,
		Title:             newConfigData.Title,
		Des:               newConfigData.Des,
		CoverFileID:       newConfigData.CoverFileID,
		DesFiles:          newConfigData.DesFiles,
		UserGroups:        newConfigData.UserGroups,
		ExemptionPrice:    newConfigData.ExemptionPrice,
		ExemptionDiscount: newConfigData.ExemptionDiscount,
		ExemptionMinPrice: newConfigData.ExemptionMinPrice,
		Limits:            newConfigData.Limits,
		StyleID:           newConfigData.StyleID,
		Params:            newConfigData.Params,
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Error("config id: ", newConfigData.ID, ", org id: ", newConfigData.OrgID, ", delete at: ", newConfigData.DeleteAt)
	}
}

func TestDeleteConfig(t *testing.T) {
	err := DeleteConfig(&ArgsDeleteConfig{
		ID:    newConfigData.ID,
		OrgID: newConfigData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
