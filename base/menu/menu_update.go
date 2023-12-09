package BaseMenu

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsUpdateMenu 更新目录参数
type ArgsUpdateMenu struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//排序
	Sort int `db:"sort" json:"sort" check:"intThan0"`
	//目录名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//图标
	Icon string `db:"icon" json:"icon"`
	//上级
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//所需权限
	OrgPermissions pq.StringArray `db:"org_permissions" json:"orgPermissions" check:"marks" empty:"true"`
	//指定组织分组
	OrgGroupIDs pq.Int64Array `db:"org_group_ids" json:"orgGroupIDs" check:"ids" empty:"true"`
	//指定组织角色
	OrgRoleIDs pq.Int64Array `db:"org_role_ids" json:"orgRoleIDs" check:"ids" empty:"true"`
	//指定组织成员
	OrgBindIDs pq.Int64Array `db:"org_bind_ids" json:"orgBindIDs" check:"ids" empty:"true"`
	//外挂模块
	WidgetSystem string `db:"widget_system" json:"widgetSystem" check:"mark"`
	// 指定对应模块配置ID
	WidgetID int64 `db:"widget_id" json:"widgetID" check:"id"`
	//访问级别
	VisitPermission string `db:"visit_permission" json:"visitPermission" check:"mark"`
}

// UpdateMenu 更新目录
func UpdateMenu(args *ArgsUpdateMenu) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_menu SET update_at = NOW(), sort = :sort, name = :name, icon = :icon, parent_id = :parent_id, org_permissions = :org_permissions, org_group_ids = :org_group_ids, org_role_ids = :org_role_ids, org_bind_ids = :org_bind_ids, widget_system = :widget_system, widget_id = :widget_id, visit_permission = :visit_permission WHERE id = :id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	deleteMenuCache(args.ID)
	return
}

// ArgsUpdateMenuSort 修改目录排序参数
type ArgsUpdateMenuSort struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//排序
	Sort int `db:"sort" json:"sort" check:"intThan0"`
}

// UpdateMenuSort 修改目录排序
func UpdateMenuSort(args *ArgsUpdateMenuSort) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_menu SET sort = :sort WHERE id = :id AND org_id = :org_id", args)
	if err != nil {
		return
	}
	deleteMenuCache(args.ID)
	return
}
