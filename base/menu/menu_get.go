package BaseMenu

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"sort"
)

// ArgsGetMenuList 获取目录配置列表参数
type ArgsGetMenuList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//上级
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetMenuList 获取目录配置列表
func GetMenuList(args *ArgsGetMenuList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "core_menu"
	var rawList []FieldsConfig
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "sort"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getMenuByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetMenuByID 获取指定到目录
func GetMenuByID(id int64, orgID int64) (data FieldsConfig) {
	data = getMenuByID(id)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(orgID, data.OrgID) {
		data = FieldsConfig{}
		return
	}
	return
}

// GetMenuByOrgBindID 获取指定成员的目录集合
func GetMenuByOrgBindID(orgBindID int64) (dataList FieldsConfigList) {
	//获取成员信息
	bindData := OrgCoreCore.GetBindNoErr(orgBindID, -1, -1)
	if bindData.ID < 1 {
		return
	}
	//获取成员的所有权限
	bindPermissions := OrgCoreCore.GetPermissionByBindID(bindData.ID)
	//获取符合条件的数据
	var rawList []FieldsConfig
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM core_menu WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND ((org_group_ids && $2) OR (org_role_ids && $3) OR $4 = ANY(org_bind_ids)) AND parent_id = 0", bindData.OrgID, bindData.GroupIDs, bindData.RoleConfigIDs, bindData.ID)
	if err != nil || len(rawList) < 1 {
		return
	}
	//写入数据
	for _, v := range rawList {
		vData := getMenuByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		if !CoreFilter.CheckArrayStringLeftMustHaveRight(bindPermissions, vData.OrgPermissions) {
			continue
		}
		dataList = append(dataList, vData)
	}
	//重新排序
	sort.Sort(dataList)
	//反馈
	return
}

// GetMenuCountByParentID 获取存在多少个下级
func GetMenuCountByParentID(parentID int64) (count int64) {
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM core_menu WHERE parent_id = $1 AND delete_at < to_timestamp(1000000)", parentID)
	if err != nil {
		return
	}
	return
}

// getMenuByID 获取ID
func getMenuByID(id int64) (data FieldsConfig) {
	cacheMark := getMenuCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, sort, name, icon, parent_id, org_permissions, org_group_ids, org_role_ids, org_bind_ids, widget_system, widget_id, visit_permission FROM core_menu WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 259200)
	return
}
