package TestBaseConfig

import (
	"encoding/json"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	TestAPI "gitee.com/weeekj/weeekj_core/v5/tools/test_api"
	"testing"
)

var (
	ConfigData BaseConfig.FieldsConfigType
)

// 获取某组配置项目
func GetConfig(t *testing.T, mark string) {
	dataByte, err := TestAPI.Get("/v1/manager/base/config/view/" + mark)
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
		Data BaseConfig.FieldsConfigType `json:"data"`
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
	ConfigData = data.Data
}

type SetParams2Type struct {
	Mark  string      `json:"mark"`
	Value interface{} `json:"value"`
}
type SetParamsType struct {
	Data []SetParams2Type `json:"data"`
}

// 修改某组配置项目
func SetConfig(t *testing.T, params SetParamsType) {
	dataByte, err := TestAPI.Post("/v1/manager/base/config/more", params)
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
		return
	}
}
