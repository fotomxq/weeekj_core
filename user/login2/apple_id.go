package UserLogin2

import (
	"context"
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreFile "github.com/fotomxq/weeekj_core/v5/core/file"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/jmind-systems/go-apple-signin"
)

// LoginByAppleID 苹果ID登录机制处理
// 参考：https://github.com/jmind-systems/go-apple-signin
// 参考流程设计：https://github.com/tptpp/sign-in-with-apple/blob/master/main.go
// 参考流程介绍：https://blog.csdn.net/tptpppp/article/details/99288426
// 参考java实现：https://www.albinwong.com/P7q15vDYQZDBbR0o.html
func LoginByAppleID(orgID int64, authCode string, referrerNationCode, referrerPhone string) (isNewUser bool, userInfo UserCore.FieldsUserType, errCode string, err error) {
	//获取配置
	teamID := BaseConfig.GetDataStringNoErr("LoginAppleIDTeamID")
	clientID := BaseConfig.GetDataStringNoErr("LoginAppleIDClientID")
	keyID := BaseConfig.GetDataStringNoErr("LoginAppleIDKeyID")
	keyP8 := BaseConfig.GetDataStringNoErr("LoginAppleIDKeyP8")
	// Pass credentials: team_id, client_id and key_id.
	opts := apple.WithCredentials(teamID, clientID, keyID)
	// Create the client.
	var client *apple.Client
	client, err = apple.NewClient(opts)
	if err != nil {
		errCode = "err_config"
		return
	}
	// Load your p8 key into the client.
	if keyP8 == "" {
		errCode = "err_key"
		err = errors.New("key p8 config is empty")
		return
	}
	//if err = client.LoadP8CertByByte([]byte(keyP8)); err != nil {
	//	errCode = "err_key"
	//	return
	//}
	if err = client.LoadP8CertByFile(fmt.Sprint(CoreFile.BaseDir(), CoreFile.Sep, "data", CoreFile.Sep, "apple_id_test.p8")); err != nil {
		errCode = "err_key"
		return
	}
	// Now client is ready.
	ctx := context.Background()
	resp, err := client.Authenticate(ctx, authCode)
	if err != nil {
		errCode = "err_user_login_code"
		return
	}
	//尝试登录或构建用户
	userEmail := ""
	if resp.UserIdentity.EmailVerified {
		userEmail = resp.UserIdentity.Email
	}
	isNewUser, userInfo, errCode, err = UserCore.LoginOrRegUser(&UserCore.ArgsLoginOrRegUser{
		OrgID:                orgID,
		Name:                 "",
		Password:             "",
		NationCode:           "",
		Phone:                "",
		AllowSkipPhoneVerify: false,
		AllowSkipWaitEmail:   true,
		Email:                userEmail,
		Username:             "",
		Avatar:               0,
		Parents:              UserCore.FieldsUserParents{},
		Groups:               UserCore.FieldsUserGroupsType{},
		Infos:                CoreSQLConfig.FieldsConfigsType{},
		Login: UserCore.FieldsUserLoginType{
			Mark:   "appleIDTokenID",
			Val:    resp.UserIdentity.ID,
			Config: resp.IDToken,
		},
		SortID: 0,
		Tags:   []int64{},
	})
	if err != nil {
		return
	}
	//尾巴处理
	if isNewUser {
		CreateAndFinal(userInfo.OrgID, userInfo.ID, &ArgsCreateUser{
			RegFrom:            "apple_auth",
			ReferrerNationCode: referrerNationCode,
			ReferrerPhone:      referrerPhone,
		})
	}
	//反馈
	return
}
