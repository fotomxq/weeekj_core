package BaseWeixinWXXUser

import (
	"errors"
	"fmt"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	UserLogin2 "github.com/fotomxq/weeekj_core/v5/user/login2"
)

// ArgsLoginByPhone 手机号授权登录
type ArgsLoginByPhone struct {
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

func LoginByPhone(args *ArgsLoginByPhone) (data UserCore.FieldsUserType, errCode string, err error) {
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
		errCode = "err_weixin_phone"
		err = errors.New(fmt.Sprint("get weixin phone data, ", err))
		return
	}
	//查询openID是否被绑定过
	findUnionIDUserData, _ := UserCore.GetUserByLogin(&UserCore.ArgsGetUserByLogin{
		OrgID: args.OrgID,
		Mark:  "weixin-union-id",
		Val:   loginData.UnionID,
	})
	findOpenIDUserData, _ := UserCore.GetUserByLogin(&UserCore.ArgsGetUserByLogin{
		OrgID: args.OrgID,
		Mark:  "weixin-open-id",
		Val:   loginData.OpenID,
	})
	//查询用户手机号
	findPhoneUserData, _ := UserCore.GetUserByPhone(&UserCore.ArgsGetUserByPhone{
		OrgID:      args.OrgID,
		NationCode: "86",
		Phone:      phoneData.PhoneNumber,
	})
	if findPhoneUserData.ID > 0 {
		//找到该手机的用户
		data = findPhoneUserData
		// 检查openID用户是否存在，剥离openID旧的用户数据
		if findUnionIDUserData.ID > 0 {
			err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
				ID:       findUnionIDUserData.ID,
				OrgID:    findUnionIDUserData.OrgID,
				Mark:     "weixin-union-id",
				Val:      loginData.UnionID,
				Config:   "",
				IsRemove: true,
			})
			if err != nil {
				errCode = "err_update"
				return
			}
		}
		if findOpenIDUserData.ID > 0 {
			err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
				ID:       findOpenIDUserData.ID,
				OrgID:    findOpenIDUserData.OrgID,
				Mark:     "weixin-open-id",
				Val:      loginData.OpenID,
				Config:   "",
				IsRemove: true,
			})
			if err != nil {
				errCode = "err_update"
				return
			}
		}
	} else {
		//找不到该手机的用户
		//如果存在UnionID用户，则先赋值
		if findUnionIDUserData.ID > 0 {
			data = findUnionIDUserData
		}
		//如果两个openID不一致，则以UnionID用户优先，玻璃openID用户
		if findUnionIDUserData.ID > 0 && findOpenIDUserData.ID > 0 && findUnionIDUserData.ID != findOpenIDUserData.ID {
			err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
				ID:       findOpenIDUserData.ID,
				OrgID:    findOpenIDUserData.OrgID,
				Mark:     "weixin-open-id",
				Val:      loginData.OpenID,
				Config:   "",
				IsRemove: true,
			})
			if err != nil {
				errCode = "err_update"
				return
			}
			data = findUnionIDUserData
		}
		if data.ID < 1 {
			//如果openID用户存在，则赋值
			if findOpenIDUserData.ID > 0 {
				data = findOpenIDUserData
			}
		}
		if data.ID < 1 {
			//创建新的用户
			if len(phoneData.PhoneNumber) != 11 {
				errCode = "err_phone"
				err = errors.New("phone error")
				return
			}
			data, errCode, err = UserLogin2.CreateUser(&UserCore.ArgsCreateUser{
				OrgID:                args.OrgID,
				Name:                 "",
				Password:             "",
				NationCode:           "86",
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
				Logins:               nil,
				SortID:               0,
				Tags:                 nil,
			}, &UserLogin2.ArgsCreateUser{
				RegFrom:            "weixin_wxx_phone",
				ReferrerNationCode: args.ReferrerNationCode,
				ReferrerPhone:      args.ReferrerPhone,
			})
			if err != nil {
				return
			}
			// 检查openID用户是否存在，剥离openID旧的用户数据
			if findUnionIDUserData.ID > 0 || findOpenIDUserData.ID > 0 {
				//以openID为基准执行登录操作
				if findOpenIDUserData.ID > 0 {
					err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
						ID:       findOpenIDUserData.ID,
						OrgID:    findOpenIDUserData.OrgID,
						Mark:     "weixin-open-id",
						Val:      loginData.OpenID,
						Config:   "",
						IsRemove: true,
					})
					if err != nil {
						errCode = "err_update"
						return
					}
				}
				//如果两个openID不是同一个用户，则继续剥离第二个用户
				if findUnionIDUserData.ID > 0 && findUnionIDUserData.ID != findOpenIDUserData.ID {
					err = UserCore.UpdateUserLoginByID(&UserCore.ArgsUpdateUserLoginByID{
						ID:       findUnionIDUserData.ID,
						OrgID:    findUnionIDUserData.OrgID,
						Mark:     "weixin-union-id",
						Val:      loginData.UnionID,
						Config:   "",
						IsRemove: true,
					})
					if err != nil {
						errCode = "err_update"
						return
					}
				}
			}
		}
	}
	//修改用户扩展信息
	if loginData.UnionID != "" {
		data.Logins = UserCore.SetUserLogins(data.Logins, "weixin-union-id", loginData.UnionID, "")
	}
	if loginData.OpenID != "" {
		data.Logins = UserCore.SetUserLogins(data.Logins, "weixin-open-id", loginData.OpenID, "")
	}
	data.Infos = CoreSQLConfig.Set(data.Infos, "referrerNationCode", args.ReferrerNationCode)
	data.Infos = CoreSQLConfig.Set(data.Infos, "referrerPhone", args.ReferrerPhone)
	errCode, err = UserCore.MargeFinal(&UserCore.ArgsMargeFinal{
		ID:      data.ID,
		Name:    data.Name,
		Avatar:  data.Avatar,
		Parents: data.Parents,
		Groups:  data.Groups,
		Infos:   data.Infos,
		Logins:  data.Logins,
		SortID:  data.SortID,
		Tags:    data.Tags,
	})
	if err != nil {
		return
	}
	//反馈
	return
}
