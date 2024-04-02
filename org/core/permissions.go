package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// GetAllPermission 获取所有权限
func GetAllPermission() (dataList []FieldsPermission) {
	cacheMark := getPermissionAllCacheMark()
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	var rawList []FieldsPermission
	_ = Router2SystemConfig.MainDB.Select(&rawList, "SELECT mark FROM org_core_permission")
	for _, v := range rawList {
		vData := getPermissionByMark(v.Mark)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, CoreCache.CacheTime1Week)
	return
}

// GetPermissionsByOrg 获取指定组织具备的权限数据
func GetPermissionsByOrg(orgID int64) (permissions []FieldsPermission) {
	//获取缓冲
	cacheMark := getPermissionOrgAllCacheMark(orgID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &permissions); err == nil && len(permissions) > 0 {
		return
	}
	//获取组织数据
	var orgData FieldsOrg
	orgData = getOrgByID(orgID)
	if orgData.ID < 1 {
		return
	}
	//获取平台所有权限
	allPermissionList := GetAllPermission()
	if len(allPermissionList) < 1 {
		return
	}
	//遍历功能数据集合
	for _, v := range orgData.OpenFunc {
		switch v {
		case "all":
			//如果不是全部模块启动，则继续
			permissions = allPermissionList
			return
		default:
			//识别该值
			for _, v2 := range allPermissionList {
				if v != v2.FuncMark {
					continue
				}
				permissions = append(permissions, v2)
			}
		}
	}
	//保存缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, permissions, CoreCache.CacheTime3Day)
	//全部完成
	return
}

// ArgsCheckPermissionsByBind 检查指定的分组或人是否具备指定的权限参数
type ArgsCheckPermissionsByBind struct {
	//绑定ID
	BindID int64
	//要检查的权限组
	// 必须同时满足
	Permissions []string
}

// ReplyCheckPermissionsByBind 检查指定的分组或人是否具备指定的权限反馈集合
type ReplyCheckPermissionsByBind struct {
	//分组列
	GroupList []FieldsGroup
	//绑定数据
	BindData FieldsBind
	//权限列
	Permissions []string
	//是否允许访问
	Allow bool
}

// CheckPermissionsByBind 检查指定的分组或人是否具备指定的权限
// Deprecated
func CheckPermissionsByBind(args *ArgsCheckPermissionsByBind) (reply ReplyCheckPermissionsByBind, err error) {
	//获取绑定数据
	reply.BindData, err = GetBind(&ArgsGetBind{
		ID:     args.BindID,
		OrgID:  0,
		UserID: 0,
	})
	if err != nil {
		//绑定关系不存在，禁止访问
		return
	}
	//获取成员分组
	reply.GroupList, _ = GetGroupMore(&ArgsGetGroupMore{
		IDs:        reply.BindData.GroupIDs,
		HaveRemove: false,
	})
	//获取权限
	reply.Permissions = GetPermissionByBindID(reply.BindData.ID)
	//检查权限
	for _, v := range args.Permissions {
		for _, v2 := range reply.Permissions {
			if v2 == "all" {
				reply.Allow = true
				return
			}
			if v == v2 {
				reply.Allow = true
				break
			}
		}
		if !reply.Allow {
			err = errors.New("no permission mark: " + v)
			return
		}
	}
	return
}

// GetPermissionByBindDataAndGroupList 通过组织成员和分组，获取符合条件的权限列
// Deprecated
func GetPermissionByBindDataAndGroupList(orgUserID int64, bindData FieldsBind, groupList []FieldsGroup) (permissionList []string) {
	//获取组织功能列
	orgPermissions := GetPermissionsByOrg(bindData.OrgID)
	if len(orgPermissions) < 1 {
		//不具备任何功能？禁止访问
		return
	}
	//如果是组织的管理人员
	if orgUserID == bindData.UserID {
		for _, v := range orgPermissions {
			permissionList = append(permissionList, v.Mark)
		}
		return
	}
	//将绑定的权限，写入反馈权限列
	for _, v := range bindData.Manager {
		if v == "all" {
			for _, v2 := range orgPermissions {
				permissionList = append(permissionList, v2.Mark)
			}
			break
		} else {
			isFind := false
			for _, v2 := range orgPermissions {
				if v == v2.Mark {
					isFind = true
					break
				}
			}
			if !isFind {
				continue
			}
			permissionList = append(permissionList, v)
		}
	}
	//获取绑定的所有分组
	for _, vGroup := range groupList {
		for _, vGroupManager := range vGroup.Manager {
			if vGroupManager == "all" {
				for _, v3 := range orgPermissions {
					isFind := false
					for _, v4 := range permissionList {
						if v3.Mark == v4 {
							isFind = true
							break
						}
					}
					if !isFind {
						permissionList = append(permissionList, v3.Mark)
					}
				}
				break
			} else {
				isFind := false
				for _, v3 := range orgPermissions {
					if vGroupManager == v3.Mark {
						isFind = true
						break
					}
				}
				if !isFind {
					continue
				}
				isFind = false
				for _, v3 := range permissionList {
					if v3 == vGroupManager {
						isFind = true
						break
					}
				}
				if !isFind {
					permissionList = append(permissionList, vGroupManager)
				}
			}
		}
	}
	return
}

// GetPermissionByBindID 通过组织成员，获取权限列
func GetPermissionByBindID(bindID int64) (permissionList []string) {
	//获取缓冲
	cacheMark := getPermissionBindCacheMark(bindID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &permissionList); err == nil && len(permissionList) > 0 {
		return
	}
	//获取成员信息
	bindData := getBindByID(bindID)
	if bindData.ID < 1 {
		return
	}
	//获取组织的权限列
	orgPermissions := GetPermissionsByOrg(bindData.OrgID)
	if len(orgPermissions) < 1 {
		return
	}
	//组装已经授权的所有权限
	var havePermissions []string
	//检查是否具备全部权限？
	haveAll := false
	//遍历组织的权限列，依次给与最终权限组
	for _, v := range bindData.Manager {
		if v == "all" {
			haveAll = true
			break
		}
		havePermissions = append(havePermissions, v)
		if haveAll {
			break
		}
	}
	//遍历组织分组的权限列，依次给与权限组
	for _, vGroupID := range bindData.GroupIDs {
		//获取分组
		vGroup := getGroupByID(vGroupID)
		if vGroup.ID < 1 {
			continue
		}
		//检查分组授予的权限列
		for _, vManager := range vGroup.Manager {
			if vManager == "all" {
				haveAll = true
				break
			}
			havePermissions = append(havePermissions, vManager)
		}
		if haveAll {
			break
		}
	}
	//如果需要授予全部权限，则直接根据组织具备权限给与
	if haveAll {
		permissionList = []string{}
		for _, v := range orgPermissions {
			permissionList = append(permissionList, v.Mark)
		}
	} else {
		//剔除组织不具备的权限然后写入权限列
		var newPermissionList []string
		for _, vPermission := range havePermissions {
			isFind := false
			for _, v2 := range orgPermissions {
				if vPermission == v2.Mark {
					isFind = true
					break
				}
			}
			if !isFind {
				continue
			}
			newPermissionList = append(newPermissionList, vPermission)
		}
		//替代最终的权限列，避免组织给与的权限超出范围
		permissionList = newPermissionList
	}
	//保存缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, permissionList, CoreCache.CacheTime1Hour)
	//反馈
	return
}

// CheckPermissionByBindID 检查成员是否具备权限
func CheckPermissionByBindID(bindID int64, permissions []string) bool {
	permissionList := GetPermissionByBindID(bindID)
	for _, v := range permissions {
		isFind := false
		for _, v2 := range permissionList {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			return false
		}
	}
	return true
}

// ArgsSetPermission 设置权限参数
type ArgsSetPermission struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//分组标识码
	FuncMark string `db:"func_mark" json:"funcMark"`
	//名称
	Name string `db:"name" json:"name"`
}

// SetPermission 设置权限
func SetPermission(args *ArgsSetPermission) (err error) {
	var dataFunc FieldsPermissionFunc
	err = Router2SystemConfig.MainDB.Get(&dataFunc, "SELECT mark FROM org_core_permission_func WHERE mark = $1", args.FuncMark)
	if err != nil {
		return
	}
	if dataFunc.Mark == "" {
		err = errors.New("func mark not exist")
		return
	}
	var data FieldsPermission
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_permission WHERE mark = $1;", args.Mark)
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_permission SET mark = :mark, func_mark = :func_mark, name = :name WHERE id = :id", map[string]interface{}{
			"id":        data.ID,
			"mark":      args.Mark,
			"func_mark": args.FuncMark,
			"name":      args.Name,
		})
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_core_permission (mark, func_mark, name) VALUES (:mark, :func_mark, :name)", args)
	}
	if err != nil {
		return
	}
	deletePermissionCache(args.Mark)
	return
}

// ArgsDeletePermission 删除指定的权限参数
type ArgsDeletePermission struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
}

// DeletePermission 删除指定的权限
func DeletePermission(args *ArgsDeletePermission) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "org_core_permission", "mark", args)
	if err != nil {
		return
	}
	deletePermissionCache(args.Mark)
	return
}

// CheckPermissionsByBindOrGroupOnlyBool 只检查权限不反馈数据
func CheckPermissionsByBindOrGroupOnlyBool(args *ArgsCheckPermissionsByBind) bool {
	permissionList := GetPermissionByBindID(args.BindID)
	for _, v := range args.Permissions {
		isFind := false
		for _, v2 := range permissionList {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			return false
		}
	}
	return true
}

// 获取指定权限
func getPermissionByMark(mark string) (data FieldsPermission) {
	cacheMark := getPermissionCacheMark(mark)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, mark, func_mark, name FROM org_core_permission WHERE mark = $1", mark)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime1Week)
	return
}

// 缓冲
func getPermissionCacheMark(mark string) string {
	return fmt.Sprint("org:core:permission:mark:", mark)
}

func getPermissionAllCacheMark() string {
	return fmt.Sprint("org:core:permission:all")
}

func getPermissionOrgAllCacheMark(orgID int64) string {
	if orgID < 1 {
		return fmt.Sprint("org:core:permission:org:all:")
	} else {
		return fmt.Sprint("org:core:permission:org:all:", orgID)
	}
}

func getPermissionBindCacheMark(bindID int64) string {
	if bindID < 1 {
		return fmt.Sprint("org:core:permission:bind:id:")
	} else {
		return fmt.Sprint("org:core:permission:bind:id:", bindID)
	}
}

func deletePermissionCache(mark string) {
	Router2SystemConfig.MainCache.DeleteMark(getPermissionCacheMark(mark))
	Router2SystemConfig.MainCache.DeleteMark(getPermissionAllCacheMark())
	Router2SystemConfig.MainCache.DeleteSearchMark(getPermissionOrgAllCacheMark(-1))
	Router2SystemConfig.MainCache.DeleteSearchMark(getPermissionBindCacheMark(-1))
}

func deletePermissionByOrgCache(orgID int64) {
	Router2SystemConfig.MainCache.DeleteSearchMark(getPermissionOrgAllCacheMark(orgID))
}

func deletePermissionByBindCache(bindID int64) {
	Router2SystemConfig.MainCache.DeleteSearchMark(getPermissionBindCacheMark(bindID))
}
