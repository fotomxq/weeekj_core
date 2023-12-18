package TestAPI

import (
	"encoding/json"
	"errors"
	"fmt"

	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

//本模块用于获取用户基本信息

var (
	LoginUsername   = "weeekj_admin"
	LoginPassword   = "weeekj_admin_weeekj_admin"
	LoginNationCode = "86"
	LoginPhone      = "17000000001"
)

// 使用标准接口登陆
// 注意应关闭验证码，否则将失败
// url: /v2/login/login/header/user/password
func LoginUserByPassword() (token int64, key string, userData UserCore.DataUserDataType, err error) {
	//修正token
	_, _, err = GetNewToken()
	if err != nil {
		return
	}
	//获取数据
	type paramsType struct {
		Username   string `json:"username" check:"username"`
		Password   string `json:"password" check:"password"`
		RememberMe bool   `json:"remember" check:"bool"`
		Vcode      string `json:"vcode" check:"vcode" empty:"true"`
	}
	params := paramsType{
		Username:   LoginUsername,
		Password:   LoginPassword,
		RememberMe: false,
		Vcode:      "v12345",
	}
	var dataByte []byte
	dataByte, err = Post("/v2/login/login/header/user/password", params)
	if err != nil {
		return
	}
	//解析数据
	type ReportDataToken struct {
		Token int64  `json:"token"`
		Key   string `json:"key"`
	}
	type ReportData struct {
		//会话
		Token ReportDataToken `json:"token"`
		//用户脱敏数据
		UserData UserCore.DataUserDataType `json:"userData"`
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
		Data ReportData `json:"data"`
	}
	var data dataType
	if err = json.Unmarshal(dataByte, &data); err != nil {
		return
	}
	if !data.Status {
		err = errors.New("status is false, code: " + data.Code + ", msg: " + data.Msg)
		return
	}
	//写入数据
	SetToken(data.Data.Token.Token, data.Data.Token.Key)
	//反馈
	token = data.Data.Token.Token
	key = data.Data.Token.Key
	userData = data.Data.UserData
	return
}

// 获取登陆用户信息
func GetUserData() (UserCore.DataUserDataType, error) {
	dataByte, err := Get("/v2/user/base/user/info")
	if err != nil {
		return UserCore.DataUserDataType{}, err
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
		Data UserCore.DataUserDataType `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		return UserCore.DataUserDataType{}, err
	}
	if !data.Status {
		return UserCore.DataUserDataType{}, errors.New("status is false, code: " + data.Code + ", msg: " + data.Msg)
	}
	//反馈
	return data.Data, nil
}

// 创建临时使用的用户
// 自动授权基础用户组权限，可调整其他用户组，用于权限封闭测试
func CreateUser() (UserCore.FieldsUserType, error) {
	type paramsType struct {
		Name       string `json:"name" filter:"Name"`
		Password   string `json:"password" filter:"Password"`
		NationCode string `json:"nationCode" filter:"NationCode"`
		Phone      string `json:"phone" filter:"Phone"`
		Email      string `json:"email" filter:"Email"`
		Username   string `json:"username" filter:"Username"`
		Status     string `json:"status" filter:"Mark"`
	}
	params := paramsType{
		Name:       "测试用户用例",
		Password:   "password_test",
		NationCode: "86",
		Phone:      "17789998887",
		Email:      "abc321@ac.com",
		Username:   "abc123abc132",
		Status:     "public",
	}
	dataByte, err := Put("/v1/manager/user/user", params)
	if err != nil {
		return UserCore.FieldsUserType{}, err
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
		Data UserCore.FieldsUserType `json:"data"`
	}
	var data dataType
	if err := json.Unmarshal(dataByte, &data); err != nil {
		return UserCore.FieldsUserType{}, err
	}
	if !data.Status {
		return UserCore.FieldsUserType{}, errors.New("status is false, code: " + data.Code + ", msg: " + data.Msg)
	}
	//反馈
	return data.Data, nil
}

// 修改用户的用户组
func UpdateUserGroups(id int64, groupMark string) error {
	type paramsType struct {
		ID         int64  `json:"id" filter:"ID"`
		Mark       string `json:"mark" filter:"Mark"`
		ExpireTime string `json:"expireTime" filter:"ExpireTime"`
		IsRemove   bool   `json:"isRemove" filter:"Bool"`
	}
	params := paramsType{
		ID:         id,
		Mark:       groupMark,
		ExpireTime: "",
		IsRemove:   false,
	}
	dataByte, err := Post("/v1/manager/user/user/info/group", params)
	if err != nil {
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
	if err := json.Unmarshal(dataByte, &data); err != nil {
		return err
	}
	if !data.Status {
		return errors.New("status is false, code: " + data.Code + ", msg: " + data.Msg)
	}
	//反馈
	return nil
}

// 删除临时创建的用户
func DeleteUser(id int64) error {
	dataByte, err := Delete(fmt.Sprint("/v1/manager/user/user/id/", id), nil)
	if err != nil {
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
	if err := json.Unmarshal(dataByte, &data); err != nil {
		return err
	}
	if !data.Status {
		return errors.New("status is false, code: " + data.Code + ", msg: " + data.Msg)
	}
	//反馈
	return nil
}
