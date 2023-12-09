package TestAPI

import (
	"encoding/json"
	"errors"
)

//本模块用于解决登陆后问题

//获取匿名token
// 同时将自动写入token组
func GetNewToken() (int64, string, error) {
	//获取数据
	dataByte, err := Put("/v2/base/token/public", nil)
	if err != nil {
		return 0, "", err
	}
	//解析数据
	type dataTokenType struct {
		Token int64 `json:"token"`
		Key   string `json:"key"`
	}
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
		Data dataTokenType `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		return 0, "", err
	}
	if !data.Status {
		return 0, "", errors.New("status is false, code: " + data.Code + ", msg: " + data.Msg)
	}
	//写入数据
	SetToken(data.Data.Token, data.Data.Key)
	//反馈
	return data.Data.Token, data.Data.Key, nil
}
