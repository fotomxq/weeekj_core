package OrgCoreCore

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetRoleConfigList 获取角色配置列表参数
type ArgsGetRoleConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRoleConfigList 获取角色配置列表
func GetRoleConfigList(args *ArgsGetRoleConfigList) (dataList []FieldsRoleConfig, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsGroup
	tableName := "org_core_role_config"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getRoleConfig(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetRoleConfig 获取指定角色配置
func GetRoleConfig(id int64, orgID int64) (data FieldsRoleConfig) {
	if id < 1 {
		return
	}
	data = getRoleConfig(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsRoleConfig{}
		return
	}
	return
}

// GetRoleConfigName 获取角色配置名称
func GetRoleConfigName(id int64) (name string) {
	if id < 1 {
		return
	}
	data := GetRoleConfig(id, -1)
	name = data.Name
	return
}

// GetRoleConfigNames 获取一组角色配置名称
func GetRoleConfigNames(ids []int64) (nameList []string) {
	nameList = []string{}
	if len(ids) < 1 {
		return
	}
	for _, v := range ids {
		vName := GetRoleConfigName(v)
		nameList = append(nameList, vName)
	}
	return
}

// GetRoleConfigMoreNames 获取一组角色配置名称
func GetRoleConfigMoreNames(ids []int64) (nameList map[int64]string) {
	nameList = map[int64]string{}
	if len(ids) < 1 {
		return
	}
	for _, v := range ids {
		nameList[v] = GetRoleConfigName(v)
	}
	return
}

// ArgsCreateRoleConfig 创建新的角色配置参数
type ArgsCreateRoleConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateRoleConfig 创建新的角色配置
func CreateRoleConfig(args *ArgsCreateRoleConfig) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_core_role_config (org_id, name, params) VALUES (:org_id,:name,:params)", args)
	if err != nil {
		return
	}
	return
}

// ArgsUpdaterRoleConfig 修改角色配置参数
type ArgsUpdaterRoleConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 无法修改，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdaterRoleConfig 修改角色配置
func UpdaterRoleConfig(args *ArgsUpdaterRoleConfig) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_role_config SET update_at = NOW(), name = :name, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteRoleConfigCache(args.ID)
	return
}

// ArgsDeleteRoleConfig 删除角色配置参数
type ArgsDeleteRoleConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 无法修改，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteRoleConfig 删除角色配置
func DeleteRoleConfig(args *ArgsDeleteRoleConfig) (err error) {
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_role_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteRoleConfigCache(args.ID)
	//获取所有符合条件的成员，并移除该角色
	deleteAllBindRoleConfigID(args.ID)
	//反馈
	return
}

func getRoleConfig(id int64) (data FieldsRoleConfig) {
	cacheMark := getRoleConfigCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, params FROM org_core_role_config WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, roleConfigCacheTime)
	return
}

// 缓冲
func getRoleConfigCacheMark(id int64) string {
	return fmt.Sprint("org:core:role:config:id:", id)
}

func deleteRoleConfigCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getRoleConfigCacheMark(id))
}
