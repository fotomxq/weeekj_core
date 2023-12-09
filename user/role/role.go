package UserRole

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	"github.com/lib/pq"
)

// ArgsGetRoleList 获取角色列表参数
type ArgsGetRoleList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//角色类型
	RoleType int64 `db:"role_type" json:"roleType" check:"id" empty:"true"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRoleList 获取角色列表
func GetRoleList(args *ArgsGetRoleList) (dataList []FieldsRole, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.RoleType > -1 {
		where = where + " AND role_type = :role_type"
		maps["role_type"] = args.RoleType
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "user_role"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, role_type, apply_id, user_id, name, country, city, gender, phone, cover_file_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetRoleID 获取指定角色参数
type ArgsGetRoleID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetRoleID 获取指定角色
func GetRoleID(args *ArgsGetRoleID) (data FieldsRole, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, role_type, apply_id, user_id, name, country, city, gender, phone, cover_file_id, cert_files, params FROM user_role WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.ID)
	if err == nil && data.ID < 1 {
		err = errors.New("data not exist")
		return
	}
	return
}

// ArgsGetRoleUserID 获取指定角色参数
type ArgsGetRoleUserID struct {
	//角色类型
	RoleType int64 `db:"role_type" json:"roleType" check:"id"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetRoleUserID 获取指定角色
func GetRoleUserID(args *ArgsGetRoleUserID) (data FieldsRole, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, role_type, apply_id, user_id, name, country, city, gender, phone, cover_file_id, cert_files, params FROM user_role WHERE user_id = $1 AND role_type = $2 AND delete_at < to_timestamp(1000000)", args.UserID, args.RoleType)
	if err == nil && data.ID < 1 {
		err = errors.New("data not exist")
		return
	}
	return
}

// ArgsGetRoleIDByUserID 通过用户获取角色ID参数
type ArgsGetRoleIDByUserID struct {
	//角色类型
	RoleType int64 `db:"role_type" json:"roleType" check:"id"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetRoleIDByUserID 通过用户获取角色ID
func GetRoleIDByUserID(args *ArgsGetRoleIDByUserID) (roleID int64, err error) {
	err = Router2SystemConfig.MainDB.Get(&roleID, "SELECT id FROM user_role WHERE role_type = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000)", args.RoleType, args.UserID)
	if err == nil && roleID > 0 {
		err = errors.New("no data")
		return
	}
	return
}

// CheckRoleByUserID 通过角色类型检查用户是否绑定关系
func CheckRoleByUserID(roleMark string, userID int64) (roleID int64, b bool) {
	typeData := GetTypeMarkNoErr(roleMark)
	if typeData.ID < 1 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&roleID, "SELECT id FROM user_role WHERE role_type = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000)", typeData.ID, userID)
	if err != nil {
		return
	}
	b = roleID > 0
	return
}

// ArgsSetRole 设置用户为指定角色参数
type ArgsSetRole struct {
	//角色类型
	RoleType int64 `db:"role_type" json:"roleType" check:"id"`
	//申请ID
	ApplyID int64 `db:"apply_id" json:"applyID" check:"id"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//姓名
	Name string `db:"name" json:"name" check:"name"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//城市编码
	City string `db:"city" json:"city" check:"cityCode"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender" check:"gender"`
	//联系电话
	Phone string `db:"phone" json:"phone" check:"phone"`
	//个人照片
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//证件列
	CertFiles pq.Int64Array `db:"cert_files" json:"certFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// SetRole 设置用户为指定角色
func SetRole(args *ArgsSetRole) (data FieldsRole, err error) {
	//找到对应的角色类型
	var dataType FieldsType
	err = Router2SystemConfig.MainDB.Get(&dataType, "SELECT id, group_ids FROM user_role_type WHERE id = $1", args.RoleType)
	if err != nil || dataType.ID < 1 {
		err = errors.New("role type not exist")
		return
	}
	//检查用户是否存在绑定关系，如果存在则修改
	var roleID int64
	err = Router2SystemConfig.MainDB.Get(&roleID, "SELECT id FROM user_role WHERE role_type = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000)", args.RoleType, args.UserID)
	if err == nil && roleID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_role SET update_at = NOW(), name = :name, country = :country, city = :city, gender = :gender, phone = :phone, cover_file_id = :cover_file_id, cert_files = :cert_files, params = :params WHERE id = :id", map[string]interface{}{
			"id":            roleID,
			"name":          args.Name,
			"country":       args.Country,
			"city":          args.City,
			"gender":        args.Gender,
			"phone":         args.Phone,
			"cover_file_id": args.CoverFileID,
			"cert_files":    args.CertFiles,
			"params":        args.Params,
		})
		if err != nil {
			return
		}
		data, err = GetRoleID(&ArgsGetRoleID{
			ID: roleID,
		})
		if err != nil {
			return
		}
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_role", "INSERT INTO user_role (role_type, apply_id, user_id, name, country, city, gender, phone, cover_file_id, cert_files, params) VALUES (:role_type,:apply_id,:user_id,:name,:country,:city,:gender,:phone,:cover_file_id,:cert_files,:params)", args, &data)
		if err != nil {
			return
		}
	}
	//为用户建立用户组关系
	// 获取用户数据
	var userData UserCore.FieldsUserType
	userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    args.UserID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	// 遍历要绑定的分组
	for _, v := range dataType.GroupIDs {
		err = UserCore.UpdateUserGroupByID(&UserCore.ArgsUpdateUserGroupByID{
			ID:       userData.ID,
			OrgID:    userData.OrgID,
			GroupID:  v,
			ExpireAt: CoreFilter.GetNowTimeCarbon().AddYears(100).Time,
			IsRemove: false,
		})
		if err != nil {
			CoreLog.Warn("user role update user group failed, user id: ", userData.ID, ", group id: ", v, ", err: ", err)
			err = nil
			continue
		}
	}
	//反馈
	return
}

// ArgsDeleteRole 删除角色参数
type ArgsDeleteRole struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteRole 删除角色
func DeleteRole(args *ArgsDeleteRole) (err error) {
	//获取数据包
	var data FieldsRole
	data, err = GetRoleID(&ArgsGetRoleID{
		ID: args.ID,
	})
	if err != nil {
		return
	}
	//获取角色类型
	var dataType FieldsType
	err = Router2SystemConfig.MainDB.Get(&dataType, "SELECT id, group_ids FROM user_role_type WHERE id = $1", data.RoleType)
	if err != nil || dataType.ID < 1 {
		err = errors.New("role type not exist")
		return
	}
	//剥离相关用户组
	for _, v := range dataType.GroupIDs {
		err = UserCore.UpdateUserGroupByID(&UserCore.ArgsUpdateUserGroupByID{
			ID:       data.ID,
			OrgID:    -1,
			GroupID:  v,
			ExpireAt: CoreFilter.GetNowTimeCarbon().SubSeconds(1).Time,
			IsRemove: true,
		})
		if err != nil {
			CoreLog.Warn("user role delete user group failed, user id: ", data.UserID, ", group id: ", v, ", err: ", err)
			err = nil
			continue
		}
	}
	//删除数据
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "user_role", "id", args)
	return
}
