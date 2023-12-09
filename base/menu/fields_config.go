package BaseMenu

import (
	"github.com/lib/pq"
	"time"
)

type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//排序
	Sort int `db:"sort" json:"sort"`
	//目录名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//图标
	Icon string `db:"icon" json:"icon"`
	//上级
	ParentID int64 `db:"parent_id" json:"parentID"`
	//所需权限
	OrgPermissions pq.StringArray `db:"org_permissions" json:"orgPermissions"`
	//指定组织分组
	OrgGroupIDs pq.Int64Array `db:"org_group_ids" json:"orgGroupIDs"`
	//指定组织角色
	OrgRoleIDs pq.Int64Array `db:"org_role_ids" json:"orgRoleIDs"`
	//指定组织成员
	OrgBindIDs pq.Int64Array `db:"org_bind_ids" json:"orgBindIDs"`
	//外挂模块
	// 支持: menu 目录模块; menu_more 多级目录模块; erp_audit 审批流程; erp_doc 文档数据集
	WidgetSystem string `db:"widget_system" json:"widgetSystem"`
	// 指定对应模块配置ID
	WidgetID int64 `db:"widget_id" json:"widgetID"`
	//访问级别
	// all 全部权限; edit 仅编辑和自己相关的数据; create 仅创建和查看; view 仅查看
	VisitPermission string `db:"visit_permission" json:"visitPermission"`
}

// FieldsConfigList 排序方法实现
type FieldsConfigList []FieldsConfig

func (t FieldsConfigList) Len() int {
	return len(t)
}
func (t FieldsConfigList) Less(i, j int) bool {
	return t[i].Sort < t[j].Sort
}

func (t FieldsConfigList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
