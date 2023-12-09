package OrgCert

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
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
		DefaultExpire: "10h",
		OrgID:         TestOrg.OrgData.ID,
		BindFrom:      "user",
		Mark:          "test_mark",
		Name:          "测试",
		Des:           "测试描述",
		CoverFileID:   0,
		DesFiles:      []int64{},
		AuditType:     "none",
		Currency:      86,
		Price:         30,
		SNLen:         10,
		TipType:       "none",
		Params:        nil,
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
		BindFrom: "",
		Mark:     "",
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
		ID:            newConfigData.ID,
		OrgID:         newConfigData.OrgID,
		DefaultExpire: newConfigData.DefaultExpire,
		Mark:          newConfigData.Mark,
		Name:          newConfigData.Name,
		Des:           newConfigData.Des,
		CoverFileID:   newConfigData.CoverFileID,
		DesFiles:      newConfigData.DesFiles,
		AuditType:     newConfigData.AuditType,
		Currency:      newConfigData.Currency,
		Price:         newConfigData.Price,
		SNLen:         newConfigData.SNLen,
		Params:        newConfigData.Params,
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

func TestClearConfig(t *testing.T) {
	TestClear(t)
}
