package TestAPI

import (
	"encoding/json"
	"errors"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
)

// 获取全局配置
func GetConfig(markList []string) ([]BaseConfig.FieldsConfigType, error) {
	//获取数据
	type paramsType struct {
		Mark []string `json:"mark"`
	}
	dataByte, err := Post("/v1/base/config", paramsType{Mark: markList})
	if err != nil {
		return []BaseConfig.FieldsConfigType{}, err
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
		Data []BaseConfig.FieldsConfigType `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		return []BaseConfig.FieldsConfigType{}, err
	}
	if !data.Status {
		return data.Data, errors.New("status is false, code: " + data.Code + ", msg: " + data.Msg)
	}
	//反馈
	return data.Data, nil
}
