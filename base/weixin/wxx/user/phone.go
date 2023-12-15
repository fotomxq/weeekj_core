package BaseWeixinWXXUser

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	BaseWeixinWXXClient "github.com/fotomxq/weeekj_core/v5/base/weixin/wxx/client"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserLogin2 "github.com/fotomxq/weeekj_core/v5/user/login2"
)

// ArgsGetPhone 通过用户加密摘要，解密用户手机号参数
type ArgsGetPhone struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
	//解码数据，通过login获取
	Code string `json:"code"`
	//包括敏感数据在内的完整用户信息的加密数据
	EncryptedData string `json:"encryptedData"`
	//加密算法的初始向量
	IV string `json:"iv"`
}

//GetPhone 通过用户加密摘要，解密用户手机号
/** 反馈结构体
{
    "phoneNumber": "13580006666",
    "purePhoneNumber": "13580006666",
    "countryCode": "86",
    "watermark":
    {
        "appid":"APPID",
        "timestamp": TIMESTAMP
    }
}
*/
func GetPhone(args *ArgsGetPhone) (loginData LoginResponseClient, data LoginPhoneNumber, err error) {
	//获取操作对象
	var client BaseWeixinWXXClient.ClientType
	client, err = BaseWeixinWXXClient.GetMerchantClient(args.OrgID)
	if err != nil {
		return
	}
	//使用登陆接口获取关键数据
	loginData, err = loginWXX(&client, args.Code)
	if err != nil {
		err = errors.New(fmt.Sprint("login wxx, ", err, ", org id: ", args.OrgID, ", code: ", args.Code, ", encryptedData: ", args.EncryptedData, ", iv: ", args.IV, ", client app id: ", client.ConfigData.AppID))
		return
	}
	data, err = decryptPhoneNumber(loginData.SessionKey, args.EncryptedData, args.IV)
	if err != nil {
		err = errors.New("decrypt phone number, " + err.Error())
	}
	CoreLog.Info("login wxx and get phone, login data: ", loginData, ", phone data: ", data)
	return
}

// ArgsLoginOrRegByPhone 使用手机号快速注册或登陆
// 该方法不会对手机号进行二次验证，会直接信任微信提供的资料
// 注意，必须启动系统开关才能执行，否则将拒绝创建
// 必须启动LoginNewOnlyPhone、LoginWeixinQuickPhone开关
// Deprecated
type ArgsLoginOrRegByPhone struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
	//解码数据，通过login获取
	Code string `json:"code"`
	//包括敏感数据在内的完整用户信息的加密数据
	EncryptedData string `json:"encryptedData"`
	//加密算法的初始向量
	IV string `json:"iv"`
	//推荐人手机号
	ReferrerNationCode string `db:"referrer_nation_code" json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
}

// LoginOrRegByPhone
// Deprecated
func LoginOrRegByPhone(args *ArgsLoginOrRegByPhone) (data UserCore.FieldsUserType, isNewUser bool, errCode string, err error) {
	//获取用户手机号
	var loginData LoginResponseClient
	var phoneData LoginPhoneNumber
	loginData, phoneData, err = GetPhone(&ArgsGetPhone{
		OrgID:         args.OrgID,
		Code:          args.Code,
		EncryptedData: args.EncryptedData,
		IV:            args.IV,
	})
	if err != nil {
		errCode = "weixin_phone"
		err = errors.New(fmt.Sprint("get weixin phone data, ", err))
		return
	}
	if len(phoneData.CountryCode) < 2 || len(phoneData.PhoneNumber) < 11 {
		errCode = "phone_type"
		err = errors.New("user phone is error")
		return
	}
	//生成扩展信息部分
	loginMark := ""
	loginValue := ""
	//检查是否已经注册，如果已经注册，则直接反馈用户信息
	//使用UnionID登陆
	if loginData.UnionID != "" {
		loginMark = "weixin-union-id"
		loginValue = loginData.UnionID
	} else {
		if loginData.OpenID != "" {
			loginMark = "weixin-open-id"
			loginValue = loginData.OpenID
		} else {
			//两个关键数据都为空，则拒绝登陆和注册
			errCode = "no_type"
			err = errors.New("login weixin, union id and open id is empty")
			return
		}
	}
	//尝试获取该用户
	data, err = UserCore.GetUserByPhone(&UserCore.ArgsGetUserByPhone{
		OrgID:      args.OrgID,
		NationCode: phoneData.CountryCode,
		Phone:      phoneData.PhoneNumber,
	})
	if err == nil {
		//获取配置项
		var LoginWeixinAutoMerge, LoginWeixinAutoMergePhone bool
		LoginWeixinAutoMerge, err = BaseConfig.GetDataBool("LoginWeixinAutoMerge")
		if err != nil {
			return
		}
		LoginWeixinAutoMergePhone, err = BaseConfig.GetDataBool("LoginWeixinAutoMergePhone")
		if err != nil {
			return
		}
		//第二次获取，检查是否微信已经绑定过账户？
		var userWxxInfo UserCore.FieldsUserType
		userWxxInfo, err = UserCore.GetUserByLogin(&UserCore.ArgsGetUserByLogin{
			OrgID: args.OrgID,
			Mark:  loginMark,
			Val:   loginValue,
		})
		//如果用户存在
		if err == nil && userWxxInfo.ID > 0 {
			//如果关闭了融合模式，则不进行融合尝试，直接退出
			if !LoginWeixinAutoMerge {
				return
			}
			//如果以手机号的账户优先，则优先保留手机号账户，将用户wxx数据迁移过来
			if LoginWeixinAutoMergePhone {
				//去掉微信账户的wxx信息
				if loginData.UnionID != "" {
					err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
						ID:       userWxxInfo.ID,
						OrgID:    args.OrgID,
						Mark:     "weixin-union-id",
						Val:      loginData.UnionID,
						Config:   "",
						IsRemove: true,
					})
					if err != nil {
						errCode = "marge_failed"
						err = errors.New("remove user login data by id, " + err.Error())
						return
					}
				}
				if loginData.OpenID != "" {
					err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
						ID:       userWxxInfo.ID,
						OrgID:    args.OrgID,
						Mark:     "weixin-open-id",
						Val:      loginData.OpenID,
						Config:   "",
						IsRemove: true,
					})
					if err != nil {
						errCode = "marge_failed"
						err = errors.New("remove user login data by id, " + err.Error())
						return
					}
				}
			} else {
				//无法完成融合，自动退出
				return
			}
		} else {
			//用户不存在，则直接给手机用户添加数据
			//为手机账户添加wxx信息
			if loginData.UnionID != "" {
				err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
					ID:       data.ID,
					OrgID:    args.OrgID,
					Mark:     "weixin-union-id",
					Val:      loginData.UnionID,
					Config:   "",
					IsRemove: false,
				})
				if err != nil {
					errCode = "update_info"
					err = errors.New("update user login data by id, " + err.Error())
					return
				}
			}
			if loginData.OpenID != "" {
				err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
					ID:       data.ID,
					OrgID:    args.OrgID,
					Mark:     "weixin-open-id",
					Val:      loginData.OpenID,
					Config:   "",
					IsRemove: false,
				})
				if err != nil {
					errCode = "update_info"
					err = errors.New("update user login data by id, " + err.Error())
					return
				}
			}
		}
		//反馈
		return
	}
	//检查是否启动了相关设定
	var LoginNewOnlyPhone, LoginWeixinQuickPhone bool
	LoginNewOnlyPhone, err = BaseConfig.GetDataBool("LoginNewOnlyPhone")
	if err != nil {
		LoginNewOnlyPhone = false
	}
	LoginWeixinQuickPhone, err = BaseConfig.GetDataBool("LoginWeixinQuickPhone")
	if err != nil {
		LoginWeixinQuickPhone = false
	}
	if !LoginNewOnlyPhone || !LoginWeixinQuickPhone {
		errCode = "config_phone"
		err = errors.New("user new only phone or quick off")
		return
	}
	//查询账户
	data, err = UserCore.GetUserByLogin(&UserCore.ArgsGetUserByLogin{
		OrgID: args.OrgID,
		Mark:  loginMark,
		Val:   loginValue,
	})
	if err == nil {
		//更新用户手机号码
		err = UserCore.UpdateUserPhoneByID(&UserCore.ArgsUpdateUserPhoneByID{
			ID:                   data.ID,
			OrgID:                args.OrgID,
			NationCode:           phoneData.CountryCode,
			Phone:                phoneData.PhoneNumber,
			AllowSkipPhoneVerify: true,
		})
		if err != nil {
			errCode = "update_info"
		}
		//反馈
		return
	}
	//生成随机昵称
	var nickName string
	if data.Name != "" {
		nickName = data.Name
	} else {
		if phoneData.PhoneNumber != "" {
			if len(phoneData.PhoneNumber) == 11 {
				nickName = fmt.Sprint(phoneData.PhoneNumber[0:3], "***", phoneData.PhoneNumber[11-4:11])
			}
		}
		if nickName == "" {
			nickName, err = CoreFilter.GetRandStr3(10)
			if err != nil {
				errCode = "rand"
				return
			}
		}
	}
	//构建登陆信息
	var logins []UserCore.FieldsUserLoginType
	if loginData.UnionID != "" {
		logins = append(logins, UserCore.FieldsUserLoginType{
			Mark:   "weixin-union-id",
			Val:    loginData.UnionID,
			Config: "",
		})
	}
	if loginData.OpenID != "" {
		logins = append(logins, UserCore.FieldsUserLoginType{
			Mark:   "weixin-open-id",
			Val:    loginData.OpenID,
			Config: "",
		})
	}
	//注册新的账户，查询是否已经存在？
	data, errCode, err = UserLogin2.CreateUser(&UserCore.ArgsCreateUser{
		OrgID:                args.OrgID,
		Name:                 nickName,
		Password:             "",
		NationCode:           phoneData.CountryCode,
		Phone:                phoneData.PhoneNumber,
		AllowSkipPhoneVerify: true,
		AllowSkipWaitEmail:   false,
		Email:                "",
		Username:             "",
		Avatar:               0,
		Status:               2,
		Parents:              nil,
		Groups:               nil,
		Infos:                nil,
		Logins:               logins,
		SortID:               0,
		Tags:                 nil,
	}, &UserLogin2.ArgsCreateUser{
		RegFrom:            "weixin_wxx_phone",
		ReferrerNationCode: args.ReferrerNationCode,
		ReferrerPhone:      args.ReferrerPhone,
	})
	isNewUser = true
	return
}
