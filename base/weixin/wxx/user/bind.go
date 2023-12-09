package BaseWeixinWXXUser

import (
	"errors"
	BaseWeixinWXXClient "gitee.com/weeekj/weeekj_core/v5/base/weixin/wxx/client"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
)

// ArgsBindExistUser 绑定存在用户
type ArgsBindExistUser struct {
	//用户ID
	UserID int64 `json:"userID" check:"id"`
	//是否覆盖不一致的用户
	// 如果发现微信被绑定到其他用户，则强制解绑
	NeedForceReplace bool `json:"needForceReplace"`
	//微信小程序数据
	Code          string             `json:"code"`
	UserData      DataWXUserInfoType `json:"userData"`
	EncryptedData string             `json:"encryptedData"`
	Signature     string             `json:"signature"`
	IV            string             `json:"iv"`
}

func BindExistUser(args *ArgsBindExistUser) (errCode string, err error) {
	//查询用户数据
	var userInfo UserCore.FieldsUserType
	userInfo, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    args.UserID,
		OrgID: -1,
	})
	if err != nil {
		errCode = "user_not_exist"
		return
	}
	//获取操作对象
	var client BaseWeixinWXXClient.ClientType
	client, err = BaseWeixinWXXClient.GetMerchantClient(userInfo.OrgID)
	if err != nil {
		errCode = "no_merchant"
		return
	}
	//检查参数
	if args.Code == "" {
		errCode = "no_code"
		err = errors.New("weixin xiaochengxu login code is empty")
		return
	}
	//开始登陆和注册
	var serverRes LoginResponseClient
	serverRes, err = loginWXX(&client, args.Code)
	if err != nil {
		errCode = "weixin_wxx"
		err = errors.New("login weixin, " + err.Error())
		return
	}
	loginMark := ""
	loginValue := ""
	//检查是否已经注册，如果已经注册，则直接反馈用户信息
	//使用UnionID登陆
	if serverRes.UnionID != "" {
		loginMark = "weixin-union-id"
		loginValue = serverRes.UnionID
	} else {
		if serverRes.OpenID != "" {
			loginMark = "weixin-open-id"
			loginValue = serverRes.OpenID
		} else {
			//两个关键数据都为空，则拒绝登陆和注册
			errCode = "no_type"
			err = errors.New("login weixin, union id and open id is empty")
			return
		}
	}
	//查询账户
	var oldUserInfo UserCore.FieldsUserType
	oldUserInfo, err = UserCore.GetUserByLogin(&UserCore.ArgsGetUserByLogin{
		OrgID: userInfo.OrgID,
		Mark:  loginMark,
		Val:   loginValue,
	})
	if err == nil {
		//如果存在，则检查是否可解绑
		if args.NeedForceReplace {
			if serverRes.UnionID != "" {
				loginMark = "weixin-union-id"
				loginValue = serverRes.UnionID
				err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
					ID:       oldUserInfo.ID,
					OrgID:    -1,
					Mark:     loginMark,
					Val:      loginValue,
					Config:   "",
					IsRemove: true,
				})
				if err != nil {
					errCode = "replace_old_user"
					return
				}
			}
			if serverRes.OpenID != "" {
				loginMark = "weixin-open-id"
				loginValue = serverRes.OpenID
				err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
					ID:       oldUserInfo.ID,
					OrgID:    -1,
					Mark:     loginMark,
					Val:      loginValue,
					Config:   "",
					IsRemove: true,
				})
				if err != nil {
					errCode = "replace_old_user"
					return
				}
			}
		} else {
			errCode = "wxx_replace"
			err = errors.New("user have weixin wxx")
			return
		}
	}
	//更新用户登录项
	if serverRes.UnionID != "" {
		loginMark = "weixin-union-id"
		loginValue = serverRes.UnionID
		err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
			ID:       userInfo.ID,
			OrgID:    -1,
			Mark:     loginMark,
			Val:      loginValue,
			Config:   "",
			IsRemove: false,
		})
		if err != nil {
			errCode = "update"
			return
		}
	}
	if serverRes.OpenID != "" {
		loginMark = "weixin-open-id"
		loginValue = serverRes.OpenID
		err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
			ID:       userInfo.ID,
			OrgID:    -1,
			Mark:     loginMark,
			Val:      loginValue,
			Config:   "",
			IsRemove: false,
		})
		if err != nil {
			errCode = "update"
			return
		}
	}
	//反馈成功
	return
}
