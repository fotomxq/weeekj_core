package BaseWeixinPayPay

import (
	BaseWeixinPayClient "gitee.com/weeekj/weeekj_core/v5/base/weixin/pay/client"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit = false
	client BaseWeixinPayClient.ClientType
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	//ToolsTest.ConfigServiceFileSrc = "../../" + ToolsTest.ConfigServiceFileSrc
	ToolsTest.Init(t)
	//初始化基础配置
	//修正BaseDir
	/**
	baseDir, err := CoreFile.BaseWDDir()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	client.BaseDir = baseDir + "/../../../builds/"
	if err := client.Init(ctx); err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	*/
}

func TestPayCreate(t *testing.T) {

}

func TestPayRefund(t *testing.T) {

}
