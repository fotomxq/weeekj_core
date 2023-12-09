package OrgCoreCore

import (
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// GetAllPermissionFunc 查询所有权限业务
func GetAllPermissionFunc() (dataList []FieldsPermissionFunc, err error) {
	cacheMark := getPermissionFuncCacheMark()
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT mark, name, des, parent_marks FROM org_core_permission_func")
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, permissionFuncCacheTime)
	return
}

// CheckOrgPermissionFunc 检查组织是否具备指定的服务
func CheckOrgPermissionFunc(orgID int64, mark string) bool {
	orgData := getOrgByID(orgID)
	if orgData.ID < 1 {
		return false
	}
	for _, v := range orgData.OpenFunc {
		if v == "all" {
			return true
		}
		if v == mark {
			return true
		}
	}
	return false
}

// ArgsSetPermissionFunc 修改权限业务组参数
type ArgsSetPermissionFunc struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//所需业务
	ParentMarks pq.StringArray `db:"parent_marks" json:"parentMarks"`
}

// SetPermissionFunc 修改权限业务组
func SetPermissionFunc(args *ArgsSetPermissionFunc) (err error) {
	var data FieldsPermissionFunc
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT mark FROM org_core_permission_func WHERE mark = $1;", args.Mark)
	if err == nil {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core_permission_func SET name = :name, des = :des, parent_marks = :parent_marks WHERE mark = :mark", args)
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_core_permission_func (mark, name, des, parent_marks) VALUES (:mark,:name,:des,:parent_marks)", args)
	}
	if err != nil {
		return
	}
	deletePermissionFuncCache()
	return
}

// ArgsDeletePermissionFunc 删除业务组参数
type ArgsDeletePermissionFunc struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
}

// DeletePermissionFunc 删除业务组
func DeletePermissionFunc(args *ArgsDeletePermissionFunc) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "org_core_permission_func", "mark", args)
	if err != nil {
		return
	}
	deletePermissionFuncCache()
	return
}

// 缓冲设计
func getPermissionFuncCacheMark() string {
	return fmt.Sprint("org:core:permission:func:all")
}

func deletePermissionFuncCache() {
	Router2SystemConfig.MainCache.DeleteMark(getPermissionFuncCacheMark())
}
