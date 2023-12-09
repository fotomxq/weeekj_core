package OrgCoreCore

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsGroup 组织内部分组设计
type FieldsGroup struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分组名称
	Name string `db:"name" json:"name"`
	//权利主张
	Manager pq.StringArray `db:"manager" json:"manager"`
	//部门领导
	ManagerOrgBindID int64 `db:"manager_org_bind_id" json:"managerOrgBindID"`
	//上级部门
	ParentID int64 `db:"parent_id" json:"parentID"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
