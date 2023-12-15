package BaseConfig

import (
	"testing"

	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
)

//本测试模块将进行增删等操作测试

var (
	isInit   = false
	testMark = "test"
)

// 测试初始化入口
func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
	//ToolsTest.ConfigDirAppend = "/../../builds/test"
	if err := Init(); err != nil {
		//t.Error(err)
		//t.Fail()
		//return
	}
}

// 测试模块
func TestLoadAndSave(t *testing.T) {
	//testMark, true, "测试数据", "string", "test data", true, 1, 2, 3, "TEST", "测试数据描述"
	err := Create(&ArgsCreate{
		Mark:        testMark,
		AllowPublic: true,
		Name:        "测试数据",
		ValueType:   "string",
		Value:       "test data",
		GroupMark:   "TEST",
		Des:         "测试数据描述",
	})
	if err != nil {
		t.Error("cannot create new config data, ", err)
	} else {
		t.Log("finish create config")
	}
	//加载所有
	configs, err := GetAll()
	if err != nil {
		t.Error("cannot load all configs, ", err)
	} else {
		t.Log("finish configs all data, ", configs)
	}
	//单个测试
	Load(testMark, t)
	//多种获取测试
	configStr, err := GetData(&ArgsGetData{
		Mark: testMark,
	})
	if err != nil {
		t.Error("cannot load config data string, ", err)
	} else {
		t.Log("finish config data string: ", configStr)
	}
	//写入数据测试
	configData, err := GetByMark(&ArgsGetByMark{
		Mark: testMark,
	})
	if err != nil {
		t.Error("cannot load config data, ", err)
	}
	//configData.UpdateHash, testMark, "test write data"
	err = UpdateByMark(&ArgsUpdateByMark{
		UpdateHash: configData.UpdateHash,
		Mark:       testMark,
		Value:      "test write data",
	})
	if err != nil {
		t.Error("cannot set config data string, ", err)
	} else {
		t.Log("finish set config data string")
	}
	//二次修改后直接获取缓冲
	configData, err = GetByMark(&ArgsGetByMark{
		Mark: testMark,
	})
	if err != nil {
		t.Error("cannot load config data, ", err)
	}
	//configData.UpdateHash, testMark, "test write data and get cache"
	err = UpdateByMark(&ArgsUpdateByMark{
		UpdateHash: configData.UpdateHash,
		Mark:       testMark,
		Value:      "test write data and get cache",
	})
	if err != nil {
		t.Error("cannot set config data string, ", err)
	} else {
		t.Log("finish set config data string")
	}
	//删除测试数据
	configData, err = GetByMark(&ArgsGetByMark{
		Mark: testMark,
	})
	if err != nil {
		t.Error("cannot load config data, ", err)
	}
	err = DeleteByMark(&ArgsDeleteByMark{
		Mark: testMark,
	})
	if err != nil {
		t.Error("cannot delete config data, ", err)
	} else {
		t.Log("finish delete config")
	}
}

// 反复加载数据并输出
func Load(mark string, t *testing.T) {
	//加载单个
	configData, err := GetByMark(&ArgsGetByMark{
		Mark: mark,
	})
	if err != nil {
		t.Error("cannot load config data, ", err)
	} else {
		t.Log("finish config data: ", configData)
	}
}
