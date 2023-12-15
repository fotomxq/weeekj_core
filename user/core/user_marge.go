package UserCore

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//融合账户处理模块
/**
1. 自动查询和创建账户
2. 根据指定特征查询符合条件的账户，如果存在则反馈并修改其他信息；否则创建新的用户
*/

// MargeLoginRegByPhone 融合通过手机号登录和注册
func MargeLoginRegByPhone(orgID int64, nationCode string, phone string, allowSkipVerify bool) (userData FieldsUserType, errCode string, err error) {
	//查询用户手机号
	userData, err = GetUserByPhone(&ArgsGetUserByPhone{
		OrgID:      orgID,
		NationCode: nationCode,
		Phone:      phone,
	})
	if err == nil && userData.ID > 0 {
		return
	}
	if len(phone) != 11 {
		errCode = "err_phone"
		err = errors.New("phone error")
		return
	}
	//不存在则创建
	userData, errCode, err = CreateUser(&ArgsCreateUser{
		OrgID:                orgID,
		Name:                 fmt.Sprint(phone[0:3], "***", phone[7]),
		Password:             "",
		NationCode:           nationCode,
		Phone:                phone,
		AllowSkipPhoneVerify: allowSkipVerify,
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
	})
	if err != nil {
		return
	}
	return
}

// MargeLogin 通过手机号/用户名/邮箱任意一个路径登录
// 如果密码给空，将跳过验证
func MargeLogin(orgID int64, userName string, password string, noLogin bool) (userInfo FieldsUserType, errCode string, err error) {
	err = Router2SystemConfig.MainDB.Get(&userInfo, "SELECT id FROM user_core WHERE ((email = $1 AND email_verify > to_timestamp(1000000)) OR username = $1 OR (phone = $1 AND phone_verify > to_timestamp(1000000))) AND (org_id = $2 OR $2 < 0) AND delete_at < to_timestamp(1000000) AND status = 2 ORDER BY id LIMIT 1", userName, orgID)
	if err != nil || userInfo.ID < 1 {
		errCode = "err_user_no_exist"
		err = errors.New("no data")
		return
	}
	userInfo = getUserByID(userInfo.ID)
	if userInfo.ID < 1 {
		errCode = "err_user_no_exist"
		err = errors.New("no data")
		return
	}
	if noLogin {
		if password != "" {
			userPassword, _ := getPasswordSha(password)
			if userInfo.Password != userPassword {
				errCode = "err_user_password"
				err = errors.New("password error")
				return
			}
		}
	} else {
		if password == "" {
			errCode = "err_user_password"
			err = errors.New("password error")
			return
		}
		userPassword, _ := getPasswordSha(password)
		if userInfo.Password != userPassword {
			errCode = "err_user_password"
			err = errors.New("password error")
			return
		}
	}
	return
}

// ArgsMargeFinal 融合用户登录收尾处理模块参数
type ArgsMargeFinal struct {
	//ID
	ID int64 `db:"id" check:"id"`
	//名称
	Name string `json:"name" check:"name"`
	//用户头像
	Avatar int64 `db:"avatar" json:"avatar"`
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

// MargeFinal 融合用户登录收尾处理模块
func MargeFinal(args *ArgsMargeFinal) (errCode string, err error) {
	//检查上级
	if len(args.Parents) > 0 {
		for _, v := range args.Parents {
			if err = checkUserParentID(args.ID, v.ParentID, v.System, []int64{}); err != nil {
				return
			}
		}
	}
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), name = :name, avatar = :avatar, parents = :parents, infos = :infos, logins = :logins, sort_id = :sort_id, tags = :tags WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"id":      args.ID,
		"name":    args.Name,
		"avatar":  args.Avatar,
		"parents": args.Parents,
		"groups":  args.Groups,
		"infos":   args.Infos,
		"logins":  args.Logins,
		"sort_id": args.SortID,
		"tags":    args.Tags,
	})
	if err != nil {
		return
	}
	deleteUserCache(args.ID)
	return
}

// MargeCheckLogins 检查一组logins是否具备用户，并自动处理和反馈数据
// 可以判断errCode来识别是否不存在用户，errCode = "err_user_no_exist"
func MargeCheckLogins(orgID int64, logins FieldsUserLoginsType) (userData FieldsUserType, errCode string, err error) {
	if len(logins) < 1 {
		errCode = "err_user_logins_no_exist"
		err = errors.New("no logins")
		return
	}
	var findUserList []FieldsUserType
	var waitMargeList FieldsUserLoginsType
	for _, vLogin := range logins {
		//检查login方式是否存在
		vUserData, _ := GetUserByLogin(&ArgsGetUserByLogin{
			OrgID: orgID,
			Mark:  vLogin.Mark,
			Val:   vLogin.Val,
		})
		if vUserData.ID > 0 {
			findUserList = append(findUserList, vUserData)
		} else {
			waitMargeList = append(waitMargeList, vLogin)
		}
	}
	if len(findUserList) < 1 {
		errCode = "err_user_no_exist"
		err = errors.New("no data")
		return
	}
	var lockUserID int64
	for _, v := range findUserList {
		if lockUserID < 1 {
			lockUserID = v.ID
			continue
		}
		if lockUserID != v.ID {
			errCode = "err_user_logins_conflict"
			err = errors.New("user logins conflict")
			return
		}
	}
	for _, v := range waitMargeList {
		err = UpdateUserLoginByID(&ArgsUpdateUserLoginByID{
			ID:       lockUserID,
			OrgID:    orgID,
			Mark:     v.Mark,
			Val:      v.Val,
			Config:   v.Config,
			IsRemove: false,
		})
		if err != nil {
			errCode = "err_update"
			err = errors.New(fmt.Sprint("update logins failed, lockUserID: ", lockUserID, ", org id: ", orgID, ", ", err))
			return
		}
	}
	userData = getUserByID(lockUserID)
	if userData.ID < 1 {
		errCode = "err_user_no_exist"
		err = errors.New("no data")
		return
	}
	return
}
