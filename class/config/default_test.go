package ClassConfig

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	"testing"
)

func TestInitDefault(t *testing.T) {
	TestInitConfig(t)
}

func TestSetConfigDefault(t *testing.T) {
	err := config.Default.SetConfigDefault(&ArgsSetConfigDefault{
		Mark:          "test_mark",
		Name:          "测试配置",
		AllowPublic:   true,
		AllowSelfView: false,
		AllowSelfSet:  false,
		ValueType:     0,
		ValueCheck:    "^[0-9a-zA-Z]{0,30}",
		ValueDefault:  "xx1",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetConfigDefaultList(t *testing.T) {
	dataList, dataCount, err := config.Default.GetConfigDefaultList(&ArgsGetConfigDefaultList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		NeedAllowPublic:   false,
		AllowPublic:       false,
		NeedAllowSelfView: false,
		AllowSelfView:     false,
		NeedAllowSelfSet:  false,
		AllowSelfSet:      false,
		Search:            "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataCount, dataList)
	}
}

func TestDeleteConfigDefault(t *testing.T) {
	err := config.Default.DeleteConfigDefault(&ArgsDeleteConfigDefault{
		Mark: "test_mark",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestClearConfigDefault(t *testing.T) {
}
