package VCodeImageCore

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	lastToken int64 = 1
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
	if err := Init(); err != nil {
		t.Error(err)
		t.Fail()
	}
	//只是测试是否可用，会直接跳出
	go Run()
}

func TestRun(t *testing.T) {
	//跳过
}

// 生成新的验证码
func TestGenerate(t *testing.T) {
	newUpdateHash := CoreFilter.GetRandNumber(10, 999)
	lastToken = int64(newUpdateHash)
	captchaInterface, err := Generate(&ArgsGenerate{
		Token: lastToken,
	})
	if err != nil {
		t.Error(err)
	} else {
		//OK输出部分
		t.Log("new captchaInterface len, ", len(captchaInterface.BinaryEncoding()))
	}
}

// 验证验证码
func TestCheck(t *testing.T) {
	var data FieldsVCodeType
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, token, value FROM core_vcode_image WHERE token = $1", lastToken); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("data: ", data)
	}
	if !Check(&ArgsCheck{
		lastToken, "123",
	}) {
		t.Log("check pass")
	} else {
		t.Error("check is error, lastToken: ", lastToken)
	}
}

// 删除所有验证码
func TestClear(t *testing.T) {
	if err := Clear(); err != nil {
		t.Error(err)
	}
}
