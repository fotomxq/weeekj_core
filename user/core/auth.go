package UserCore

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
)

// ArgsLoginOrRegUser 授权方式处理登录或注册参数
type ArgsLoginOrRegUser struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `json:"name" check:"name" empty:"true"`
	//密码
	Password string `json:"password" check:"password" empty:"true"`
	//电话联系
	NationCode string `json:"nationCode" check:"nationCode" empty:"true"`
	Phone      string `json:"phone" check:"phone" empty:"true"`
	//是否跳过手机号验证
	AllowSkipPhoneVerify bool `json:"allowSkipPhoneVerify" check:"bool"`
	//是否跳过验证验证？
	AllowSkipWaitEmail bool `json:"allowSkipWaitEmail"`
	//邮件
	Email string `json:"email"`
	//用户名称
	Username string `json:"username"`
	//用户头像
	Avatar int64 `db:"avatar" json:"avatar"`
	//上级
	Parents FieldsUserParents `json:"parents"`
	//分组
	Groups FieldsUserGroupsType `json:"groups"`
	//扩展信息
	Infos CoreSQLConfig.FieldsConfigsType `json:"infos"`
	//登陆信息
	Login FieldsUserLoginType `json:"login"`
	//用户分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//用户标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
}

// LoginOrRegUser 授权方式处理登录或注册
// 用于第三方登录方式，统一方法并处理融合账户的问题
func LoginOrRegUser(args *ArgsLoginOrRegUser) (isNewUser bool, userInfo FieldsUserType, errCode string, err error) {
	//通过登录信息获取账户信息
	userInfo, err = GetUserByLogin(&ArgsGetUserByLogin{
		OrgID: args.OrgID,
		Mark:  args.Login.Mark,
		Val:   args.Login.Val,
	})
	if err == nil && userInfo.ID > 0 {
		return
	}
	//创建新的用户
	userInfo, errCode, err = CreateUser(&ArgsCreateUser{
		OrgID:                args.OrgID,
		Name:                 args.Name,
		Password:             args.Password,
		NationCode:           args.NationCode,
		Phone:                args.Phone,
		AllowSkipPhoneVerify: args.AllowSkipPhoneVerify,
		AllowSkipWaitEmail:   args.AllowSkipWaitEmail,
		Email:                args.Email,
		Username:             args.Username,
		Avatar:               args.Avatar,
		Status:               2,
		Parents:              args.Parents,
		Groups:               args.Groups,
		Infos:                CoreSQLConfig.FieldsConfigsType{},
		Logins: FieldsUserLoginsType{
			args.Login,
		},
		SortID: args.SortID,
		Tags:   args.Tags,
	})
	if err != nil {
		return
	}
	isNewUser = true
	//反馈
	return
}
