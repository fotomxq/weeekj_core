package BaseDistribution

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit = false
)

// 测试初始化入口
func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	//初始化配置路径
	ToolsTest.Init(t)
	if err := Init(); err != nil {
		t.Error(err)
	}
}

func TestSetService(t *testing.T) {
	if err := SetService(&ArgsSetService{
		Mark: "test", Name: "测试服务", ExpireInterval: 30,
	}); err != nil {
		t.Error(err)
	}
}

func TestGetServiceList(t *testing.T) {
	data, count, err := GetServiceList(&ArgsGetServiceList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "Mark",
			Desc: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
	}
}

func TestDeleteService(t *testing.T) {
	if err := DeleteService(&ArgsDeleteService{
		Mark: "test",
	}); err != nil {
		t.Error(err)
	}
}
