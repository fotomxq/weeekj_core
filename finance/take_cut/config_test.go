package FinanceTakeCut

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newConfigData FieldsConfig
)

func TestInitConfig(t *testing.T) {
	TestInit(t)
	OrgCoreCore.Init(true, true)
}

func TestSetConfig(t *testing.T) {
	err := SetConfig(&ArgsSetConfig{
		SortID:             0,
		OrgID:              TestOrg.OrgData.ID,
		OrderSystem:        "mall",
		CutPriceProportion: 500000,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetConfigList(t *testing.T) {
	dataList, dataCount, err := GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:  -1,
		SortID: -1,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newConfigData = dataList[0]
	}
}

func TestDeleteConfig(t *testing.T) {
	err := DeleteConfig(&ArgsDeleteConfig{
		ID: newConfigData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearConfig(t *testing.T) {
	TestClear(t)
}
