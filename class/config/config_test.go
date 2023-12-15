package ClassConfig

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	config Config
	bindID int64 = 123
)

func TestInitConfig(t *testing.T) {
	ToolsTest.Init(t)
	config.Init("test_config", "test_config_default")
	TestSetConfigDefault(t)
}

func TestSetConfig(t *testing.T) {
	err := config.SetConfig(&ArgsSetConfig{
		BindID:    bindID,
		Mark:      "test_mark",
		VisitType: "admin",
		Val:       "value_data",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetConfig(t *testing.T) {
	configDefault, configData, err := config.GetConfig(&ArgsGetConfig{
		BindID:    bindID,
		Mark:      "test_mark",
		VisitType: "admin",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(configDefault, configData)
	}
}

func TestGetConfigList(t *testing.T) {
	dataList, dataList2, dataCount, err := config.GetConfigList(&ArgsGetConfigList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		BindID:    bindID,
		VisitType: "admin",
		Search:    "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList, dataList2, dataCount)
	}
}

func TestGetConfigValue(t *testing.T) {
	data, err := config.GetConfigVal(&ArgsGetConfig{
		BindID:    bindID,
		Mark:      "test_mark",
		VisitType: "admin",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestDeleteConfig(t *testing.T) {
	if err := config.DeleteConfig(&ArgsDeleteConfig{
		Mark:   "test_mark",
		BindID: bindID,
	}); err != nil {
		t.Error(err)
	}
}

func TestClearConfig(t *testing.T) {
	TestDeleteConfigDefault(t)
}
