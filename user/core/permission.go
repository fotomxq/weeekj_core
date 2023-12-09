package UserCore

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

//权限

// ArgsGetAllPermission 获取所有权限参数
type ArgsGetAllPermission struct {
	//组织是否可以授权
	// 如果为false，则表明为平台方
	AllowOrg bool `db:"allowOrg" json:"allowOrg"`
}

// GetAllPermission 获取所有权限
func GetAllPermission(args *ArgsGetAllPermission) (dataList []FieldsPermissionType, err error) {
	var rawList []FieldsPermissionType
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT mark FROM user_permissions WHERE allow_org = $1;", args.AllowOrg)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getPermissionByMark(v.Mark)
		if vData.Mark == "" {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetPermissionByMark 获取某个权限参数
type ArgsGetPermissionByMark struct {
	//标识码
	Mark string
}

// GetPermissionByMark 获取某个权限
func GetPermissionByMark(args *ArgsGetPermissionByMark) (data FieldsPermissionType, err error) {
	data = getPermissionByMark(args.Mark)
	if data.Mark == "" {
		err = errors.New("no data")
		return
	}
	return
}

func getPermissionByMark(mark string) (data FieldsPermissionType) {
	cacheMark := getPermissionCacheMark(mark)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.Mark != "" {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT mark, name, des, allow_org FROM user_permissions WHERE mark = $1", mark)
	if data.Mark == "" {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cachePermissionTime)
	return
}

// ArgsCreatePermission 写入新的权限参数
type ArgsCreatePermission struct {
	//标识码
	Mark string `db:"mark"`
	//名称
	Name string `db:"name"`
	//备注
	Des string `db:"des"`
	//组织是否可以授权
	AllowOrg bool `db:"allowOrg" json:"allowOrg"`
}

// CreatePermission 写入新的权限
func CreatePermission(args *ArgsCreatePermission) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_permissions (mark, name, des, allow_org) VALUES (:mark,:name,:des,:allowOrg)", args)
	if err != nil {
		return
	}
	deletePermissionCache(args.Mark)
	return
}

// ArgsUpdatePermission 修改权限描述信息参数
type ArgsUpdatePermission struct {
	//标识码
	Mark string `db:"mark"`
	//名称
	Name string `db:"name"`
	//备注
	Des string `db:"des"`
	//组织是否可以授权
	AllowOrg bool `db:"allowOrg" json:"allowOrg"`
}

// UpdatePermission 修改权限描述信息
func UpdatePermission(args *ArgsUpdatePermission) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_permissions SET  name = :name, des = :des, allow_org = :allowOrg WHERE mark = :mark", args)
	if err != nil {
		return
	}
	deletePermissionCache(args.Mark)
	return
}

// ArgsDeletePermission 删除权限参数
type ArgsDeletePermission struct {
	//标识码
	Mark string `db:"mark"`
	//是否跳过验证？
	SkipCheckGroup bool
}

// DeletePermission 删除权限
func DeletePermission(args *ArgsDeletePermission) (err error) {
	if !args.SkipCheckGroup && getGroupCountByPermission(args.Mark) > 0 {
		return errors.New("permission use in user group")
	}
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "user_permissions", "mark", args)
	if err != nil {
		return
	}
	deletePermissionCache(args.Mark)
	return
}

// DeleteAllPermission 清理权限
func DeleteAllPermission() (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_permissions", "", nil)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.DeleteSearchMark("user:core:permission:mark:")
	return
}

// 缓冲
func getPermissionCacheMark(mark string) string {
	return fmt.Sprint("user:core:permission:mark:", mark)
}

func deletePermissionCache(mark string) {
	Router2SystemConfig.MainCache.DeleteMark(getPermissionCacheMark(mark))
}
