package OrgCoreCore

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/lib/pq"
	"time"
)

// FieldsOperate 资源控制关系
type FieldsOperate struct {
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
	//绑定来源
	// system: bind | id: 成员ID
	// system: group | id: 分组ID
	BindInfo CoreSQLFrom.FieldsFrom `db:"bind_info" json:"bindInfo"`
	//控制资源的来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//控制权限列
	Manager pq.StringArray `db:"manager" json:"manager"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
