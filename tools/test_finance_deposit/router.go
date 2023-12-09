package TestFinanceDeposit

import (
	"encoding/json"
	FinanceDeposit "gitee.com/weeekj/weeekj_core/v5/finance/deposit"
	TestAPI "gitee.com/weeekj/weeekj_core/v5/tools/test_api"
	"testing"
)

//路由组件

var (
	RouterDeposit           FinanceDeposit.FieldsDepositType
	RouterDepositConfigData FinanceDeposit.FieldsConfigType
)

// 创建储蓄资金池配置
func RouterCreateConfig(t *testing.T, mark string, takeOut bool, name, des string) {
	type paramsType struct {
		Name    string `json:"name" filter:"Des" filterMin:"1" filterMax:"300"`
		Des     string `json:"des" filter:"Des" filterMin:"1" filterMax:"3000"`
		Takeout bool   `json:"takeout" filter:"Bool"`
		Mark    string `json:"mark" filter:"Mark"`
	}
	params := paramsType{
		Name:    name,
		Des:     des,
		Takeout: takeOut,
		Mark:    mark,
	}
	dataByte, err := TestAPI.Post("/v1/manager/finance/deposit/config/set", params)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data FinanceDeposit.FieldsConfigType `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
	} else {
		t.Log(data.Data, data.Count)
		RouterDepositConfigData = data.Data
	}
}

// 检查目标资金量
func RouterGetPriceByFrom(t *testing.T, fromSystem, fromID, fromMark string) {
	type paramsType struct {
		FromSystem string `json:"fromSystem" filter:"Mark"`
		FromID     string `json:"fromID" filter:"ID" empty:"true"`
		FromMark   string `json:"fromMark" filter:"Mark" empty:"true"`
	}
	params := paramsType{
		FromSystem: fromSystem,
		FromID:     fromID,
		FromMark:   fromMark,
	}
	dataByte, err := TestAPI.Post("/v1/manager/finance/deposit/from", params)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data FinanceDeposit.FieldsDepositType `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
	} else {
		t.Log(data.Data, data.Count)
		RouterDeposit = data.Data
	}
}

// 设置目标资金量
func RouterSetPrice(t *testing.T, updateHash, fromSystem, fromID, fromMark, fromName string, mark, currency string, price int64) {
	type paramsType struct {
		UpdateHash      string `json:"updateHash" filter:"Mark"`
		FromSystem      string `json:"fromSystem" filter:"Mark"`
		FromID          string `json:"fromID" filter:"ID" empty:"true"`
		FromMark        string `json:"fromMark" filter:"Mark" empty:"true"`
		FromName        string `json:"fromName" filter:"Des" filterMin:"1" filterMax:"300"`
		SaveMark        string `json:"saveMark" filter:"Mark"`
		SaveCurrency    string `json:"saveCurrency" filter:"Mark"`
		SaveAppendPrice int64  `json:"saveAppendPrice"`
	}
	params := paramsType{
		UpdateHash:      updateHash,
		FromSystem:      fromSystem,
		FromID:          fromID,
		FromMark:        fromMark,
		FromName:        fromName,
		SaveMark:        mark,
		SaveCurrency:    currency,
		SaveAppendPrice: price,
	}
	dataByte, err := TestAPI.Post("/v1/manager/finance/deposit/set", params)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data FinanceDeposit.FieldsDepositType `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
		return
	}
	t.Log(data.Data, data.Count)
	RouterDeposit = data.Data
}

// 删除配置
func RouterDeleteConfig(t *testing.T, mark string) {
	type paramsType struct {
		Mark string `json:"mark" filter:"Mark"`
	}
	params := paramsType{
		Mark: mark,
	}
	dataByte, err := TestAPI.Delete("/v1/manager/finance/deposit/config/mark", params)
	if err != nil {
		t.Error(err)
		return
	}
	//解析数据
	type dataType struct {
		//错误信息
		Status bool `json:"status"`
		//错误信息
		Code string `json:"code"`
		//错误描述
		Msg string `json:"msg"`
		//数据个数
		Count int64 `json:"count"`
		//数据集合
		Data interface{} `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		t.Error(err)
		return
	}
	if !data.Status {
		t.Error("status is false, code: " + data.Code + ", msg: " + data.Msg)
	}
}
