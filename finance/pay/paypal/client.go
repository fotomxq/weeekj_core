package FinancePayPaypal

import (
	"errors"
	"fmt"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	"github.com/plutov/paypal"
)

// 获取请求前缀部分
func getClient(orgID int64) (client *paypal.Client, err error) {
	//获取配置
	var systemConfig OrgCore.FieldsSystem
	systemConfig, err = OrgCore.GetSystem(&OrgCore.ArgsGetSystem{
		OrgID:      orgID,
		SystemMark: "paypal",
	})
	if err != nil {
		err = errors.New(fmt.Sprint("get system config, ", err))
		return
	}
	//ID和密钥
	clientID, b := systemConfig.Params.GetVal("clientID")
	if !b {
		err = errors.New(fmt.Sprint("get config by clientID, ", err))
		return
	}
	secretID, b := systemConfig.Params.GetVal("secretID")
	if !b {
		err = errors.New(fmt.Sprint("get config by secretID, ", err))
		return
	}
	isSandBox, b := systemConfig.Params.GetValBool("isSandBox")
	if !b {
		err = errors.New(fmt.Sprint("get config by sandBox, ", err))
		return
	}
	//创建订单
	apiURL := paypal.APIBaseSandBox
	if !isSandBox {
		apiURL = paypal.APIBaseLive
	}
	client, err = paypal.NewClient(clientID, secretID, apiURL)
	if err != nil {
		err = errors.New(fmt.Sprint("create paypal client, ", err))
		return
	}
	//反馈
	return
}
