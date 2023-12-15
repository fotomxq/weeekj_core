package UserCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateUserParent 修改上级ID参数
type ArgsUpdateUserParent struct {
	//ID
	ID int64 `db:"id" check:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//上级ID组
	Parents FieldsUserParents `db:"parents"`
}

// UpdateUserParent 修改上级ID
func UpdateUserParent(args *ArgsUpdateUserParent) (err error) {
	if len(args.Parents) > 0 {
		for _, v := range args.Parents {
			if err = checkUserParentID(args.ID, v.ParentID, v.System, []int64{}); err != nil {
				return
			}
		}
	}
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET  update_at = NOW(), parents = :parents WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", args); err != nil {
		return
	}
	if err != nil {
		return
	}
	deleteUserCache(args.ID)
	return
}

// ArgsUpdateUserStatus 修改用户信息参数
type ArgsUpdateUserStatus struct {
	//ID
	ID int64 `db:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//状态
	// 0 -> ban后，用户可以正常登录，但一切内容无法使用
	// 1 -> audit后，用户登录后无法正常使用，但提示不太一样
	// 2 -> public 正常访问
	Status int `db:"status" json:"status"`
}

// UpdateUserStatus 修改用户信息
func UpdateUserStatus(args *ArgsUpdateUserStatus) (err error) {
	//检查status
	if err = checkUserStatus(args.Status); err != nil {
		return
	}
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET  update_at = NOW(), status = :status WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", args); err != nil {
		return
	}
	if err != nil {
		return
	}
	deleteUserCache(args.ID)
	return
}

// ArgsUpdateUserInfoByID 修改用户信息参数
type ArgsUpdateUserInfoByID struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织ID
	NewOrgID int64 `db:"new_org_id" json:"newOrgID"`
	//名称
	Name string `db:"name" json:"name"`
	//用户头像
	Avatar int64 `db:"avatar" json:"avatar" check:"id" empty:"true"`
}

// UpdateUserInfoByID 修改用户信息
func UpdateUserInfoByID(args *ArgsUpdateUserInfoByID) (err error) {
	if args.NewOrgID < 0 {
		args.NewOrgID = args.OrgID
	}
	if args.OrgID > 0 && args.OrgID != args.NewOrgID {
		//获取旧的账户信息，开始验证
		var oldData FieldsUserType
		oldData, err = GetUserByID(&ArgsGetUserByID{
			ID:    args.ID,
			OrgID: args.OrgID,
		})
		if err != nil {
			return
		}
		var data FieldsUserType
		if oldData.NationCode != "" && oldData.Phone != "" {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND phone = $2 AND nation_code = $3", args.NewOrgID, oldData.Phone, oldData.NationCode)
			if err == nil && data.ID > 0 {
				err = errors.New("phone replace")
				return
			}
		}
		if oldData.Username != "" {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND username = $2", args.NewOrgID, oldData.Username)
			if err == nil && data.ID > 0 {
				err = errors.New("username replace")
				return
			}
		}
		if oldData.Email != "" {
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND email = $2", args.NewOrgID, oldData.Email)
			if err == nil && data.ID > 0 {
				err = errors.New("email replace")
				return
			}
		}
		for _, v := range oldData.Logins {
			if v.Mark == "" {
				continue
			}
			data, err = GetUserByLogin(&ArgsGetUserByLogin{
				OrgID: args.NewOrgID,
				Mark:  v.Mark,
				Val:   v.Val,
			})
			if err == nil && data.ID > 0 {
				err = errors.New("login replace")
				return
			}
		}
	}
	if args.Avatar < 1 {
		if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), name = :name WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
			"name":   args.Name,
		}); err != nil {
			return
		}
	} else {
		if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), name = :name, avatar = :avatar WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", args); err != nil {
			return
		}
	}
	if args.NewOrgID > -1 {
		if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), name = :name, org_id = :new_org_id WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
			"id":         args.ID,
			"new_org_id": args.NewOrgID,
			"org_id":     args.OrgID,
			"name":       args.Name,
		}); err != nil {
			return
		}
	}
	deleteUserCache(args.ID)
	if args.Avatar == 0 {
		pushNatsCreateAvatar(args.ID)
	}
	return
}

// ArgsUpdateUserAvatar 修改用户头像参数
type ArgsUpdateUserAvatar struct {
	//ID
	ID int64 `db:"id" check:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//头像ID
	AvatarID int64 `db:"avatar" json:"avatar" check:"id"`
}

// UpdateUserAvatar 修改用户头像
func UpdateUserAvatar(args *ArgsUpdateUserAvatar) (err error) {
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), avatar = :avatar WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", args); err != nil {
		return
	}
	if err != nil {
		return
	}
	deleteUserCache(args.ID)
	return
}

// ArgsUpdateUserSort 修改用户分类和标签参数
type ArgsUpdateUserSort struct {
	//ID
	ID int64 `db:"id" check:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//用户标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
}

// UpdateUserSort 修改用户分类和标签
func UpdateUserSort(args *ArgsUpdateUserSort) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), sort_id = :sort_id, tags = :tags WHERE id = :id AND (:org_id < 0 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteUserCache(args.ID)
	return
}

// ArgsUpdateUserPasswordByID 修改用户密码参数
type ArgsUpdateUserPasswordByID struct {
	//ID
	ID int64 `db:"id" check:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//密码
	Password string `db:"password"`
}

// UpdateUserPasswordByID 修改用户密码
func UpdateUserPasswordByID(args *ArgsUpdateUserPasswordByID) (err error) {
	//构建密码摘要
	args.Password, err = getPasswordSha(args.Password)
	if err != nil {
		return
	}
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), password = :password WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", args); err != nil {
		return
	}
	deleteUserCache(args.ID)
	return
}

// ArgsUpdateUserPhoneByID 修改用户电话参数
type ArgsUpdateUserPhoneByID struct {
	//ID
	ID int64 `db:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//电话联系
	// 允许留空，将抹除数据
	NationCode string `db:"nation_code"`
	Phone      string `db:"phone"`
	//是否跳过手机号验证
	AllowSkipPhoneVerify bool `json:"allowSkipPhoneVerify" check:"bool"`
}

// UpdateUserPhoneByID 修改用户电话
func UpdateUserPhoneByID(args *ArgsUpdateUserPhoneByID) (err error) {
	//检查不重复，注意如果是自己，则直接成功
	if args.NationCode == "" && args.Phone == "" {
		var data FieldsUserType
		data, err = GetUserByPhone(&ArgsGetUserByPhone{
			OrgID:      args.OrgID,
			NationCode: args.NationCode,
			Phone:      args.Phone,
		})
		if err == nil {
			if data.ID == args.ID {
				return
			}
			return errors.New("phone is exist")
		}
	}
	//验证手机号
	var phoneVerify time.Time
	if args.AllowSkipPhoneVerify {
		phoneVerify = CoreFilter.GetNowTime()
	}
	//修改
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), nation_code = :nation_code, phone = :phone, phone_verify = :phone_verify WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"id":           args.ID,
		"org_id":       args.OrgID,
		"nation_code":  args.NationCode,
		"phone":        args.Phone,
		"phone_verify": phoneVerify,
	}); err != nil {
		return
	}
	//删除缓冲
	deleteUserCache(args.ID)
	//通知新的手机号
	if args.AllowSkipPhoneVerify {
		if args.Phone != "" {
			pushNatsNewPhone(args.ID, args.NationCode, args.Phone)
		}
	}
	//如果需要同步修改用户名，则触发
	if Router2SystemConfig.GlobConfig.User.SyncUserPhoneUsername {
		_ = UpdateUserUsernameByID(&ArgsUpdateUserUsernameByID{
			ID:       args.ID,
			OrgID:    args.OrgID,
			Username: args.Phone,
		})
	}
	//反馈
	return
}

// ArgsUpdateUserEmailByID 修改用户邮箱参数
type ArgsUpdateUserEmailByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//邮件
	// 可以留空，将抹除该数据
	Email string `db:"email" json:"email" check:"email"`
	//是否跳过验证？
	AllowSkip bool `json:"allowSkip" check:"bool" empty:"true"`
}

// UpdateUserEmailByID 修改用户邮箱
func UpdateUserEmailByID(args *ArgsUpdateUserEmailByID) (err error) {
	//检查不重复，注意如果是自己，则直接成功
	if args.Email != "" {
		var data FieldsUserType
		data, err = GetUserByEmail(&ArgsGetUserByEmail{
			OrgID: args.OrgID,
			Email: args.Email,
		})
		if err == nil {
			if data.ID == args.ID {
				return nil
			}
			return errors.New("email is exist")
		}
	}
	//修改
	if args.AllowSkip {
		if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), email = :email, email_verify = :email_verify WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
			"id":           args.ID,
			"org_id":       args.OrgID,
			"email":        args.Email,
			"email_verify": CoreFilter.GetNowTime(),
		}); err != nil {
			return
		}
	} else {
		if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), email = :email, email_verify = :email_verify WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
			"id":           args.ID,
			"org_id":       args.OrgID,
			"email":        args.Email,
			"email_verify": time.Time{},
		}); err != nil {
			return
		}
		pushNatsUserEmailWait(args.ID)
	}
	deleteUserCache(args.ID)
	return
}

// ArgsUpdateUserUsernameByID 修改用户用户名参数
type ArgsUpdateUserUsernameByID struct {
	//ID
	ID int64 `db:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户名
	Username string `db:"username"`
}

// UpdateUserUsernameByID 修改用户用户名
func UpdateUserUsernameByID(args *ArgsUpdateUserUsernameByID) (err error) {
	//检查用户名是否重复
	if err = checkUsername(args.ID, args.OrgID, args.Username); err != nil {
		return
	}
	//修改
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), username = :username WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", args); err != nil {
		return
	}
	deleteUserCache(args.ID)
	return nil
}

// 检查用户名是否重复
func checkUsername(checkID int64, orgID int64, username string) (err error) {
	var data FieldsUserType
	data, err = GetUserByUsername(&ArgsGetUserByUsername{
		OrgID:    orgID,
		Username: username,
	})
	if err == nil {
		if data.ID != checkID {
			return errors.New("username is exist")
		}
	} else {
		err = nil
	}
	return
}

type ArgsUpdateUserGroupByID struct {
	//ID
	ID int64 `db:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户组ID
	GroupID int64
	//过期时间
	ExpireAt time.Time
	//是否删除
	IsRemove bool
}

func UpdateUserGroupByID(args *ArgsUpdateUserGroupByID) (err error) {
	//确保mark存在
	var groupInfo FieldsGroupType
	groupInfo, err = GetGroup(&ArgsGetGroup{
		ID: args.GroupID,
	})
	if err != nil {
		return errors.New("group is not exist, " + err.Error())
	}
	//获取数据
	var data FieldsUserType
	data, err = GetUserByID(&ArgsGetUserByID{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return err
	}
	var newGroups FieldsUserGroupsType
	//寻找groups
	isFind := false
	for _, v := range data.Groups {
		if v.GroupID == groupInfo.ID {
			isFind = true
			if args.IsRemove {
				continue
			}
			v.ExpireAt = args.ExpireAt
		}
		newGroups = append(newGroups, v)
	}
	if !isFind && !args.IsRemove {
		newGroups = append(newGroups, FieldsUserGroupType{
			GroupID:  groupInfo.ID,
			CreateAt: CoreFilter.GetNowTime(),
			ExpireAt: args.ExpireAt,
		})
	}
	//修改
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), groups = :groups WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
		"groups": newGroups,
	}); err != nil {
		return
	}
	deleteUserCache(args.ID)
	return nil
}

// ArgsUpdateUserInfosByID 修改用户的信息扩展参数
type ArgsUpdateUserInfosByID struct {
	//ID
	ID int64 `db:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//参数
	Mark string
	//值
	Val string
	//是否删除
	IsRemove bool
}

// UpdateUserInfosByID 修改用户的信息扩展
func UpdateUserInfosByID(args *ArgsUpdateUserInfosByID) (err error) {
	//获取数据
	var data FieldsUserType
	data, err = GetUserByID(&ArgsGetUserByID{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return err
	}
	var newInfos CoreSQLConfig.FieldsConfigsType
	//寻找groups
	isFind := false
	for _, v := range data.Infos {
		if v.Mark == args.Mark {
			isFind = true
			if args.IsRemove {
				continue
			}
			v.Val = args.Val
		}
		newInfos = append(newInfos, v)
	}
	if !isFind {
		newInfos = append(newInfos, CoreSQLConfig.FieldsConfigType{
			Mark: args.Mark,
			Val:  args.Val,
		})
	}
	//修改
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET update_at = NOW(), infos = :infos WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
		"infos":  newInfos,
	}); err != nil {
		return
	}
	deleteUserCache(args.ID)
	return
}

// ArgsUpdateUserLoginByID 修改用户登录项参数
type ArgsUpdateUserLoginByID struct {
	//ID
	ID int64
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//参数
	Mark string `db:"mark" json:"mark" check:"mark"`
	//值
	Val string
	//配置
	Config string
	//是否删除
	IsRemove bool
}

// UpdateUserLoginByID 修改用户登录项
func UpdateUserLoginByID(args *ArgsUpdateUserLoginByID) (err error) {
	//获取数据
	var data FieldsUserType
	data, err = GetUserByID(&ArgsGetUserByID{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return err
	}
	//检查该mark/value是否已经存在？
	var findData FieldsUserType
	findData, err = GetUserByLogin(&ArgsGetUserByLogin{
		OrgID: args.OrgID,
		Mark:  args.Mark,
		Val:   args.Val,
	})
	if err == nil {
		//如果不是本用户，则报错
		if findData.ID != args.ID {
			err = errors.New(fmt.Sprint("login mark and value is replace, "))
			return
		}
	}
	var newLogin FieldsUserLoginsType
	//寻找login
	isFind := false
	for _, v := range data.Logins {
		if v.Mark == args.Mark {
			isFind = true
			if args.IsRemove {
				continue
			}
			v.Val = args.Val
			v.Config = args.Config
		}
		newLogin = append(newLogin, v)
	}
	if !isFind {
		newLogin = append(newLogin, FieldsUserLoginType{
			Mark:   args.Mark,
			Val:    args.Val,
			Config: args.Config,
		})
	}
	//修改
	if _, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_core SET  update_at = NOW(), logins = :logins WHERE id = :id AND (:org_id < 0 OR org_id = :org_id) AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"id":     args.ID,
		"org_id": args.OrgID,
		"logins": newLogin,
	}); err != nil {
		return
	}
	deleteUserCache(args.ID)
	return nil
}
