package OrgSubscription

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitConfig(t *testing.T) {
	TestInit(t)
}

func TestSetConfig(t *testing.T) {
	var err error
	newConfigData, err = SetConfig(&ArgsSetConfig{
		Mark:        "test",
		FuncList:    []string{"only", "all"},
		TimeType:    1,
		TimeN:       15,
		Currency:    86,
		Price:       15,
		PriceOld:    16,
		Title:       "测试标题",
		Des:         "测试描述",
		CoverFileID: 0,
		DesFiles:    []int64{},
		StyleID:     0,
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newConfigData)
}

func TestGetConfigByID(t *testing.T) {
	data, err := GetConfigByID(&ArgsGetConfigByID{
		ID: newConfigData.ID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetConfigByMark(t *testing.T) {
	data, err := GetConfigByMark(&ArgsGetConfigByMark{
		Mark: newConfigData.Mark,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetConfigList(t *testing.T) {
	dataList, dataCount, err := GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDeleteConfig(t *testing.T) {
	err := DeleteConfig(&ArgsDeleteConfig{
		ID: newConfigData.ID,
	})
	ToolsTest.ReportError(t, err)
	TestClear(t)
}
