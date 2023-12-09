package BaseWeixinWXXUser

import (
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	BaseFileSys "gitee.com/weeekj/weeekj_core/v5/base/filesys"
	BaseQiniu "gitee.com/weeekj/weeekj_core/v5/base/qiniu"
	BaseWeixinWXXClient "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	UserLogin2 "gitee.com/weeekj/weeekj_core/v5/user/login2"
	"time"
)

// ArgsLoginAndReg 自动登陆和注册处理接口参数
type ArgsLoginAndReg struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
	//商户ID
	// 可以留空，则走平台微信小程序主体
	Code          string             `json:"code"`
	UserData      DataWXUserInfoType `json:"userData"`
	EncryptedData string             `json:"encryptedData"`
	Signature     string             `json:"signature"`
	IV            string             `json:"iv"`
	//推荐人手机号
	ReferrerNationCode string `db:"referrer_nation_code" json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
}

// LoginAndReg 登陆验证模块
// API参考：https://developers.weixin.qq.com/miniprogram/dev/api/api-login.html
// 用于解决微信小程序登陆和匹配功能
// 完成后自动给与token
func LoginAndReg(args *ArgsLoginAndReg) (userInfo UserCore.FieldsUserType, isNewUser bool, errCode string, err error) {
	//获取微信openID
	var wxxData DataGetOpenIDByCode
	wxxData, errCode, err = GetOpenIDByCode(args.OrgID, args.Code)
	if err != nil {
		return
	}
	//检查登录方式是否一致
	var waitCheckLogins UserCore.FieldsUserLoginsType
	if wxxData.UnionID != "" {
		waitCheckLogins = append(waitCheckLogins, UserCore.FieldsUserLoginType{
			Mark:   "weixin-union-id",
			Val:    wxxData.UnionID,
			Config: "",
		})
	}
	if wxxData.OpenID != "" {
		waitCheckLogins = append(waitCheckLogins, UserCore.FieldsUserLoginType{
			Mark:   "weixin-open-id",
			Val:    wxxData.OpenID,
			Config: "",
		})
	}
	userInfo, errCode, err = UserCore.MargeCheckLogins(args.OrgID, waitCheckLogins)
	if err != nil {
		if errCode == "err_user_logins_no_exist" || errCode == "err_user_no_exist" {
			err = nil
			errCode = ""
		} else {
			return
		}
	} else {
		//自动修正名称开关，如果未启动，则跳出
		var LoginWeixinAutoFixName bool
		LoginWeixinAutoFixName, err = BaseConfig.GetDataBool("LoginWeixinAutoFixName")
		if err != nil {
			LoginWeixinAutoFixName = true
		}
		if !LoginWeixinAutoFixName {
			return
		}
		//检查如果用户不存在头像和昵称，再修改
		if userInfo.Avatar < 1 {
			//下载图片并推送到七牛云
			var fileData BaseFileSys.FieldsFileClaimType
			fileData, errCode, err = BaseQiniu.UploadByURL(&BaseQiniu.ArgsUploadByURL{
				URL:        args.UserData.AvatarUrl,
				BucketName: "",
				FileType:   "jpg",
				IP:         "0.0.0.0",
				CreateInfo: CoreSQLFrom.FieldsFrom{
					System: "user",
					ID:     userInfo.ID,
					Mark:   "",
					Name:   userInfo.Name,
				},
				UserID:     userInfo.ID,
				OrgID:      userInfo.OrgID,
				IsPublic:   true,
				ExpireAt:   time.Time{},
				ClaimInfos: []CoreSQLConfig.FieldsConfigType{},
				Des:        "",
			})
			if err != nil {
				return
			}
			userInfo.Avatar = fileData.ID
		}
		if userInfo.Name == "" {
			//修正数据
			if err = UserCore.UpdateUserInfoByID(&UserCore.ArgsUpdateUserInfoByID{
				ID:       userInfo.ID,
				OrgID:    userInfo.OrgID,
				NewOrgID: -1,
				Name:     args.UserData.NickName,
				Avatar:   userInfo.Avatar,
			}); err != nil {
				errCode = "update_info"
				err = errors.New("update user info by id, " + err.Error())
				return
			}
			userInfo.Name = args.UserData.NickName
		} else {
			if err = UserCore.UpdateUserInfoByID(&UserCore.ArgsUpdateUserInfoByID{
				ID:       userInfo.ID,
				OrgID:    userInfo.OrgID,
				NewOrgID: -1,
				Name:     userInfo.Name,
				Avatar:   userInfo.Avatar,
			}); err != nil {
				errCode = "update_info"
				err = errors.New("update user info by id, " + err.Error())
				return
			}
		}
		//反馈
		return
	}
	//组合用户infos
	infos := args.UserData.GetUserInfos()
	//构建登陆信息
	var logins []UserCore.FieldsUserLoginType
	if wxxData.UnionID != "" {
		logins = append(logins, UserCore.FieldsUserLoginType{
			Mark:   "weixin-union-id",
			Val:    wxxData.UnionID,
			Config: "",
		})
	}
	if wxxData.OpenID != "" {
		logins = append(logins, UserCore.FieldsUserLoginType{
			Mark:   "weixin-open-id",
			Val:    wxxData.OpenID,
			Config: "",
		})
	}
	//这里说明没有注册过的新用户，自动完成注册
	//用户信息组
	userInfo, errCode, err = UserLogin2.CreateUser(&UserCore.ArgsCreateUser{
		OrgID:                args.OrgID,
		Name:                 args.UserData.NickName,
		Password:             "",
		NationCode:           "",
		Phone:                "",
		AllowSkipPhoneVerify: false,
		AllowSkipWaitEmail:   false,
		Email:                "",
		Username:             "",
		Avatar:               0,
		Status:               2,
		Parents:              nil,
		Groups:               nil,
		Infos:                infos,
		Logins:               logins,
		SortID:               0,
		Tags:                 nil,
	}, &UserLogin2.ArgsCreateUser{
		RegFrom:            "weixin_wxx",
		ReferrerNationCode: args.ReferrerNationCode,
		ReferrerPhone:      args.ReferrerPhone,
	})
	if err != nil {
		return
	}
	isNewUser = true
	if err != nil {
		err = errors.New("create user, " + err.Error())
		return
	}
	return
}

// DataGetOpenIDByCode 获取微信OpenID数据
type DataGetOpenIDByCode struct {
	//OpenID
	OpenID string `json:"openid"`
	//UnionID
	UnionID string `json:"unionid"`
}

// GetOpenIDByCode 获取微信OpenID
func GetOpenIDByCode(orgID int64, code string) (result DataGetOpenIDByCode, errCode string, err error) {
	//获取操作对象
	var client BaseWeixinWXXClient.ClientType
	client, err = BaseWeixinWXXClient.GetMerchantClient(orgID)
	if err != nil {
		errCode = "no_merchant"
		return
	}
	//检查参数
	if code == "" {
		errCode = "no_code"
		err = errors.New("weixin xiaochengxu login code is empty")
		return
	}
	//开始登陆和注册
	var serverRes LoginResponseClient
	serverRes, err = loginWXX(&client, code)
	if err != nil {
		errCode = "weixin_wxx"
		err = errors.New("login weixin, " + err.Error())
		return
	}
	//检查是否已经注册，如果已经注册，则直接反馈用户信息
	//使用UnionID登陆
	result = DataGetOpenIDByCode{}
	if serverRes.UnionID != "" {
		result.UnionID = serverRes.UnionID
	}
	if serverRes.OpenID != "" {
		result.OpenID = serverRes.OpenID
	}
	if result.UnionID == "" && result.OpenID == "" {
		errCode = "err_user_login_code"
		err = errors.New(fmt.Sprint("open id or unionid is empty, body data: SessionKey: ", serverRes.SessionKey, ", OpenID: ", serverRes.OpenID, ", UnionID: ", serverRes.UnionID))
		return
	}
	return
}
