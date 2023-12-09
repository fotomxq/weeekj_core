package UserCore

import (
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateUser 插入新的用户参数
type ArgsCreateUser struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `json:"name" check:"name"`
	//密码
	Password string `json:"password" check:"password"`
	//电话联系
	NationCode string `json:"nationCode" check:"nationCode" empty:"true"`
	Phone      string `json:"phone" check:"phone" empty:"true"`
	//是否跳过手机号验证
	AllowSkipPhoneVerify bool `json:"allowSkipPhoneVerify" check:"bool"`
	//是否跳过验证验证？
	AllowSkipWaitEmail bool `json:"allowSkipWaitEmail"`
	//邮件
	Email string `json:"email" check:"email" empty:"true"`
	//用户名称
	Username string `json:"username" check:"username" empty:"true"`
	//用户头像
	Avatar int64 `db:"avatar" json:"avatar"`
	//状态
	// 0 -> ban后，用户可以正常登录，但一切内容无法使用
	// 1 -> audit后，用户登录后无法正常使用，但提示不太一样
	// 2 -> public 正常访问
	Status int `json:"status"`
	//上级
	Parents FieldsUserParents `json:"parents"`
	//分组
	Groups FieldsUserGroupsType `json:"groups"`
	//扩展信息
	Infos CoreSQLConfig.FieldsConfigsType `json:"infos"`
	//登陆信息
	Logins FieldsUserLoginsType `json:"logins"`
	//用户分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//用户标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
}

// CreateUser 插入新的用户
func CreateUser(args *ArgsCreateUser) (userInfo FieldsUserType, errCode string, err error) {
	//创建用户锁定机制
	createLock.Lock()
	defer createLock.Unlock()
	//优化参数设计
	if args.Parents == nil {
		args.Parents = FieldsUserParents{}
	}
	if args.Groups == nil {
		args.Groups = FieldsUserGroupsType{}
	}
	if args.Infos == nil {
		args.Infos = CoreSQLConfig.FieldsConfigsType{}
	}
	if args.Logins == nil {
		args.Logins = FieldsUserLoginsType{}
	}
	if args.Tags == nil {
		args.Tags = []int64{}
	}
	//检查status
	err = checkUserStatus(args.Status)
	if err != nil {
		errCode = "err_status"
		return
	}
	//检查密码
	var passwordSha string
	if args.Password != "" {
		if b := checkPassword(args.Password); !b {
			errCode = "err_password"
			err = errors.New("password type is error")
			return
		}
		//构建密码摘要
		passwordSha, err = getPasswordSha(args.Password)
		if err != nil {
			errCode = "err_sha"
			return
		}
	}
	//检查国家代码
	var phoneVerify time.Time
	if args.NationCode != "" || args.Phone != "" {
		if b := CoreFilter.CheckNationCode(args.NationCode); !b {
			errCode = "err_phone"
			err = errors.New("phone nation code is error")
			return
		}
		if b := CoreFilter.CheckPhone(args.Phone); !b {
			errCode = "err_phone"
			err = errors.New("phone is error")
			return
		}
		//检查手机号、邮箱、用户名唯一性
		_, err = GetUserByPhone(&ArgsGetUserByPhone{
			OrgID:      args.OrgID,
			NationCode: args.NationCode,
			Phone:      args.Phone,
		})
		if err == nil {
			errCode = "err_phone_exist"
			err = errors.New("user phone is exist")
			return
		}
		//检查是否跳过手机号验证
		if args.AllowSkipPhoneVerify {
			phoneVerify = CoreFilter.GetNowTime()
		}
	}
	//如果存在email，则检查email的唯一性
	if args.Email != "" {
		var data FieldsUserType
		data, err = GetUserByEmail(&ArgsGetUserByEmail{
			OrgID: args.OrgID,
			Email: args.Email,
		})
		if err == nil {
			if data.ID > 0 {
				errCode = "err_email_exist"
				err = errors.New("user email is exist")
				return
			}
		}
	}
	//如果存在用户名，则检查用户名的唯一性
	if args.Username != "" {
		_, err = GetUserByUsername(&ArgsGetUserByUsername{
			OrgID:    args.OrgID,
			Username: args.Username,
		})
		if err == nil {
			errCode = "err_username_exist"
			err = errors.New("user username is exist")
			return
		}
	}
	//处理用户默认用户组设置
	var userNewDefaultGroupID int64
	userNewDefaultGroupID, err = BaseConfig.GetDataInt64("UserNewDefaultGroupID")
	if err != nil {
		/**
		errCode = "new_user_group_config_not_exist"
		err = errors.New("cannot get config, " + err.Error())
		return
		*/
	}
	//获取默认的用户组，并授权组织相关配置关键数据包
	if args.OrgID > 0 {
		//TODO: 获取组织的配置项，用于授权
	}
	//生成过期时间
	var expireAt time.Time
	expireAt, err = CoreFilter.GetTimeByAdd("876000h")
	if err != nil {
		errCode = "err_time"
		return
	}
	//生成用户组信息结构体
	if userNewDefaultGroupID > 0 && len(args.Groups) < 1 {
		args.Groups = []FieldsUserGroupType{
			{
				GroupID:  userNewDefaultGroupID,
				CreateAt: CoreFilter.GetNowTime(),
				ExpireAt: expireAt,
			},
		}
	}
	//构建数据
	maps := map[string]interface{}{
		"status":       args.Status,
		"org_id":       args.OrgID,
		"name":         args.Name,
		"password":     passwordSha,
		"username":     args.Username,
		"nation_code":  args.NationCode,
		"phone":        args.Phone,
		"phone_verify": phoneVerify,
		"email":        args.Email,
		"avatar":       args.Avatar,
		"parents":      args.Parents,
		"groups":       args.Groups,
		"infos":        args.Infos,
		"logins":       args.Logins,
		"sort_id":      args.SortID,
		"tags":         args.Tags,
	}
	//检查是否允许跳过邮箱验证
	if args.AllowSkipWaitEmail && args.Email != "" {
		maps["email_verify"] = CoreFilter.GetNowTime()
	} else {
		maps["email_verify"] = time.Time{}
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_core", "INSERT INTO user_core (delete_at, status, org_id, name, password, nation_code, phone, phone_verify, email, email_verify, username, avatar, parents, groups, infos, logins, sort_id, tags) VALUES (to_timestamp(0), :status,:org_id,:name,:password,:nation_code,:phone,:phone_verify,:email,:email_verify,:username,:avatar,:parents,:groups,:infos,:logins,:sort_id,:tags)", maps, &userInfo)
	if err != nil || userInfo.ID < 1 {
		errCode = "err_insert"
		return
	}
	//通知新增了用户
	pushNatsCreateUser(userInfo)
	//反馈
	return
}
