package BaseWeixinApp

import (
	"errors"
	"fmt"
	BaseQiniu "github.com/fotomxq/weeekj_core/v5/base/qiniu"
	CoreHttp "github.com/fotomxq/weeekj_core/v5/core/http"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserLogin2 "github.com/fotomxq/weeekj_core/v5/user/login2"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

// ArgsUserLogin 用户登陆和授权接口处理参数
type ArgsUserLogin struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//编码
	Code string `json:"code"`
	//推荐人手机号
	ReferrerNationCode string `db:"referrer_nation_code" json:"referrerNationCode" check:"nationCode" empty:"true"`
	ReferrerPhone      string `json:"referrerPhone" check:"phone" empty:"true"`
}

// UserLogin 用户登陆和授权接口处理
func UserLogin(args *ArgsUserLogin) (isNewUser bool, userData UserCore.FieldsUserType, errCode string, err error) {
	//获取openID等数据
	var openIDData DataGetOpenIDByCode
	openIDData, errCode, err = GetOpenIDByCode(args.OrgID, args.Code)
	if err != nil {
		return
	}
	//检查用户是否存在？
	var waitCheckLogins []UserCore.FieldsUserLoginType
	if openIDData.UnionID != "" {
		waitCheckLogins = append(waitCheckLogins, UserCore.FieldsUserLoginType{
			Mark:   "weixin_app_unionid",
			Val:    openIDData.UnionID,
			Config: "",
		})
	}
	if openIDData.OpenID != "" {
		waitCheckLogins = append(waitCheckLogins, UserCore.FieldsUserLoginType{
			Mark:   "weixin_app_openid",
			Val:    openIDData.OpenID,
			Config: "",
		})
	}
	userData, errCode, err = UserCore.MargeCheckLogins(args.OrgID, waitCheckLogins)
	if err != nil {
		if errCode == "err_user_no_exist" {
			err = nil
			errCode = ""
		} else {
			return
		}
	} else {
		//用户存在，注册成功反馈
		updateUserAvatarByURL(openIDData.Headimgurl, userData.ID)
		return
	}
	//注册新的用户
	var logins []UserCore.FieldsUserLoginType
	if openIDData.UnionID != "" {
		logins = append(logins, UserCore.FieldsUserLoginType{
			Mark:   "weixin_app_unionid",
			Val:    openIDData.UnionID,
			Config: "",
		})
	}
	if openIDData.OpenID != "" {
		logins = append(logins, UserCore.FieldsUserLoginType{
			Mark:   "weixin_app_openid",
			Val:    openIDData.OpenID,
			Config: "",
		})
	}
	userData, _, err = UserLogin2.CreateUser(&UserCore.ArgsCreateUser{
		OrgID:                args.OrgID,
		Name:                 openIDData.Nickname,
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
		Infos:                nil,
		Logins:               logins,
		SortID:               0,
		Tags:                 nil,
	}, &UserLogin2.ArgsCreateUser{
		RegFrom:            "weixin_app",
		ReferrerNationCode: args.ReferrerNationCode,
		ReferrerPhone:      args.ReferrerPhone,
	})
	if err != nil {
		errCode = "err_user_reg"
		return
	}
	//标记新用户
	isNewUser = true
	//保存用户头像
	updateUserAvatarByURL(openIDData.Headimgurl, userData.ID)
	//反馈
	return
}

// DataGetOpenIDByCode 获取微信OpenID数据
type DataGetOpenIDByCode struct {
	//OpenID
	OpenID string `json:"openid"`
	//UnionID
	UnionID string `json:"unionid"`
	//昵称
	Nickname string `json:"nickname"`
	//头像地址
	Headimgurl string `json:"headimgurl"`
}

// GetOpenIDByCode 获取微信OpenID
func GetOpenIDByCode(orgID int64, code string) (result DataGetOpenIDByCode, errCode string, err error) {
	//获取商户配置
	appID, appKey := getAppConfig(orgID)
	//构建URL地址
	postAuthURL := strings.ReplaceAll(userLoginAuthURL, "$1", appID)
	postAuthURL = strings.ReplaceAll(postAuthURL, "$2", appKey)
	postAuthURL = strings.ReplaceAll(postAuthURL, "$3", code)
	//向微信发起请求
	var bodyAuthByte []byte
	bodyAuthByte, err = CoreHttp.GetData(postAuthURL, nil, "", false)
	if err != nil {
		errCode = "err_http"
		return
	}
	//分配内存
	result = DataGetOpenIDByCode{}
	//解析数据
	accessToken := gjson.GetBytes(bodyAuthByte, "access_token").String()
	result.OpenID = gjson.GetBytes(bodyAuthByte, "openid").String()
	if accessToken == "" || result.OpenID == "" {
		errCode = "err_user_login_code"
		err = errors.New(fmt.Sprint("access token or open id empty, body data: ", string(bodyAuthByte)))
		return
	}
	//请求用户基本信息
	postInfoURL := strings.ReplaceAll(userLoginAuthInfoURL, "$1", accessToken)
	postInfoURL = strings.ReplaceAll(postInfoURL, "$2", result.OpenID)
	//向微信发起请求
	var bodyInfoByte []byte
	bodyInfoByte, err = CoreHttp.GetData(postInfoURL, nil, "", false)
	if err != nil {
		errCode = "err_http"
		return
	}
	//解析数据
	result.OpenID = gjson.GetBytes(bodyInfoByte, "openid").String()
	result.Nickname = gjson.GetBytes(bodyInfoByte, "nickname").String()
	result.Headimgurl = gjson.GetBytes(bodyInfoByte, "headimgurl").String()
	result.UnionID = gjson.GetBytes(bodyInfoByte, "unionid").String()
	if result.UnionID == "" && result.OpenID == "" {
		errCode = "err_user_login_code"
		err = errors.New(fmt.Sprint("open id or unionid is empty, body data: ", string(bodyInfoByte)))
		return
	}
	return
}

// 更新用户头像
func updateUserAvatarByURL(imgURL string, userID int64) {
	if imgURL == "" {
		return
	}
	//上传用户头像
	fileData, _, err := BaseQiniu.UploadByURL(&BaseQiniu.ArgsUploadByURL{
		URL:        imgURL,
		BucketName: "",
		FileType:   "",
		IP:         "0.0.0.0",
		CreateInfo: CoreSQLFrom.FieldsFrom{},
		UserID:     userID,
		OrgID:      0,
		IsPublic:   true,
		ExpireAt:   time.Time{},
		ClaimInfos: nil,
		Des:        "",
	})
	if err != nil {
		return
	}
	err = UserCore.UpdateUserAvatar(&UserCore.ArgsUpdateUserAvatar{
		ID:       userID,
		OrgID:    -1,
		AvatarID: fileData.ID,
	})
	if err != nil {
		return
	}
}
