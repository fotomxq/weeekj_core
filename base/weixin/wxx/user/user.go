package BaseWeixinWXXUser

import (
	"errors"
	"fmt"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
)

// ArgsGetOpenIDByUserInfo 从user数据中获取openid数据参数
type ArgsGetOpenIDByUserInfo struct {
	//用户数据
	UserInfo UserCore.FieldsUserType
}

// GetOpenIDByUserInfo 从user数据中获取openid数据
func GetOpenIDByUserInfo(args *ArgsGetOpenIDByUserInfo) (string, error) {
	for _, v := range args.UserInfo.Logins {
		if v.Mark == "weixin-open-id" {
			return v.Val, nil
		}
	}
	return "", errors.New(fmt.Sprint("user is not have weixin openid, user id: ", args.UserInfo.ID, ", userData: ", args.UserInfo))
}

// DataWXUserInfoType 提交来的定义组参数
type DataWXUserInfoType struct {
	//昵称
	NickName string `json:"nickName"`
	//性别 1男 2女 0未知
	Gender int `json:"gender"`
	//语言 zh_CN
	Language string `json:"language"`
	//城市
	City string `json:"city"`
	//省份
	Province string `json:"province"`
	//国家编号
	Country string `json:"country"`
	//头像URL地址
	AvatarUrl string `json:"avatarUrl"`
}

// GetUserInfos 提交来的定义组
// 通过提交数据，解码用户信息结构体
func (t *DataWXUserInfoType) GetUserInfos() (data CoreSQLConfig.FieldsConfigsType) {
	data = append(data, CoreSQLConfig.FieldsConfigType{
		Mark: "WXX_NickName",
		Val:  t.NickName,
	})
	data = append(data, CoreSQLConfig.FieldsConfigType{
		Mark: "WXX_Gender",
		Val:  fmt.Sprint(t.Gender),
	})
	data = append(data, CoreSQLConfig.FieldsConfigType{
		Mark: "WXX_Language",
		Val:  t.Language,
	})
	data = append(data, CoreSQLConfig.FieldsConfigType{
		Mark: "WXX_City",
		Val:  t.City,
	})
	data = append(data, CoreSQLConfig.FieldsConfigType{
		Mark: "WXX_Province",
		Val:  t.Province,
	})
	data = append(data, CoreSQLConfig.FieldsConfigType{
		Mark: "WXX_Country",
		Val:  t.Country,
	})
	data = append(data, CoreSQLConfig.FieldsConfigType{
		Mark: "WXX_AvatarUrl",
		Val:  t.AvatarUrl,
	})
	return
}
