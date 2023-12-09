package TestAPI

import (
	"encoding/json"
	"errors"
)

//测试发送短信验证码
// 注意，必须具备了前置token，才可以使用本接口
// 可用于部分需要验证码的环节
func SendSMS(nationCode, phone string) error{
	if nationCode == ""{
		nationCode = LoginNationCode
	}
	if phone == ""{
		phone = LoginPhone
	}
	//获取数据
	type DataType struct {
		CountryCode string `json:"countryCode" filter:"NationCode"`
		Phone string `json:"phone" filter:"Phone"`
		Vcode string `json:"vcode" filter:"Vcode"`
	}
	params := DataType{
		CountryCode: nationCode,
		Phone:       phone,
		Vcode:       "v12345",
	}
	dataByte, err := Put("/v1/base/verification_code/sms", params)
	if err != nil{
		return err
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
	if err := json.Unmarshal(dataByte, &data); err != nil{
		return err
	}
	if !data.Status{
		return errors.New("status is false, code: " + data.Code + ", msg: " + data.Msg)
	}
	//反馈
	return nil
}