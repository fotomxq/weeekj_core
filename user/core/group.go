package UserCore

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//用户组

// ArgsGetAllGroup 查看所有用户组参数
type ArgsGetAllGroup struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
}

// GetAllGroup 查看所有用户组
func GetAllGroup(args *ArgsGetAllGroup) (dataList []FieldsGroupType, err error) {
	var rawList []FieldsGroupType
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM user_groups WHERE org_id = $1 OR $1 < 0;", args.OrgID)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getGroupByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetGroup 获取某个用户组参数
type ArgsGetGroup struct {
	//用户组ID
	ID int64 `db:"id" json:"id"`
}

// GetGroup 获取某个用户组
func GetGroup(args *ArgsGetGroup) (data FieldsGroupType, err error) {
	data = getGroupByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

func GetGroupName(id int64) string {
	data := getGroupByID(id)
	return data.Name
}

// getGroupByID 获取指定的用户组ID
func getGroupByID(id int64) (data FieldsGroupType) {
	cacheMark := getGroupCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, name, des, permissions FROM user_groups WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheGroupTime)
	return
}

// GetGroupPermissionList 获取指定用户组的所有权限
func GetGroupPermissionList(groupIDs pq.Int64Array) []string {
	var permissionList []string
	for _, v := range groupIDs {
		vData := getGroupByID(v)
		if vData.ID < 1 {
			continue
		}
		permissionList = CoreFilter.MargeArrayString(permissionList, vData.Permissions)
	}
	return permissionList
}

// ArgsCreateGroup 创建新的用户组参数
type ArgsCreateGroup struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name"`
	//备注
	Des string `db:"des"`
	//权限
	Permissions pq.StringArray `db:"permissions"`
}

// CreateGroup 创建新的用户组
func CreateGroup(args *ArgsCreateGroup) (data FieldsGroupType, err error) {
	//检查权限是否存在
	for _, v := range args.Permissions {
		if _, err = GetPermissionByMark(&ArgsGetPermissionByMark{
			Mark: v,
		}); err != nil {
			err = errors.New("permission is not exist, mark: " + v)
			return
		}
	}
	//创建新的组
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_groups", "INSERT INTO user_groups (org_id, name, des, permissions) VALUES (:org_id,:name,:des,:permissions)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateGroup 修改用户组参数
type ArgsUpdateGroup struct {
	//用户组ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name"`
	//备注
	Des string `db:"des"`
	//权限
	Permissions pq.StringArray `db:"permissions"`
}

// UpdateGroup 修改用户组
func UpdateGroup(args *ArgsUpdateGroup) (err error) {
	//检查权限是否存在
	for _, v := range args.Permissions {
		if _, err = GetPermissionByMark(&ArgsGetPermissionByMark{
			Mark: v,
		}); err != nil {
			err = errors.New("permission is not exist, mark is " + v)
			return
		}
	}
	//执行操作
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_groups SET  name = :name, org_id = :org_id, des = :des, permissions = :permissions WHERE id = :id", args)
	if err != nil {
		return
	}
	deleteGroupCache(args.ID)
	return
}

// ArgsDeleteGroup 删除用户组参数
type ArgsDeleteGroup struct {
	//用户组ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
}

// DeleteGroup 删除用户组
func DeleteGroup(args *ArgsDeleteGroup) (err error) {
	if getUserCountByGroup(args.ID, args.OrgID) > 0 {
		return errors.New("group is used")
	}
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_groups", "id = :id AND (org_id = :org_id OR :org_id < 1)", args)
	if err != nil {
		return
	}
	deleteGroupCache(args.ID)
	return
}

// 检查某个权限存在于用户组的个数
func getGroupCountByPermission(permissionMark string) (count int64) {
	var err error
	count, err = CoreSQL.GetAllCount(Router2SystemConfig.MainDB.DB, "user_groups", "id", "$1 = ANY(permissions)", permissionMark)
	if err != nil {
		return
	}
	return
}

// CheckGroupIsOrg 检查用户组是否属于商户？
func CheckGroupIsOrg(orgID int64, groupIDs []int64) (err error) {
	if len(groupIDs) < 1 {
		return
	}
	var groupList []FieldsGroupType
	for _, v := range groupIDs {
		vData := getGroupByID(v)
		if vData.ID < 1 {
			continue
		}
		groupList = append(groupList, vData)
	}
	if len(groupList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range groupList {
		if v.OrgID != orgID {
			err = errors.New("group not org")
			return
		}
	}
	return
}

// 缓冲
func getGroupCacheMark(id int64) string {
	return fmt.Sprint("user:core:group:id:", id)
}

func deleteGroupCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getGroupCacheMark(id))
}
