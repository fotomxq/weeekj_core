package OrgCoreCoreMod

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsBind 分组和其他模块来源的关系
// 来源主要以用户为主
type FieldsBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//最后1次登陆时间
	LastAt time.Time `db:"last_at" json:"lastAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//用户头像
	Avatar int64 `db:"avatar" json:"avatar" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//组织分组ID列
	GroupIDs pq.Int64Array `db:"group_ids" json:"groupIDs"`
	//角色配置列
	RoleConfigIDs pq.Int64Array `db:"role_config_ids" json:"roleConfigIDs"`
	//权利主张
	Manager pq.StringArray `db:"manager" json:"manager"`
	//联系电话
	NationCode string `db:"nation_code" json:"nationCode"`
	Phone      string `db:"phone" json:"phone"`
	//邮件地址
	Email string `db:"email" json:"email"`
	//同步专用设计
	// 可用于同步其他系统来源
	SyncSystem string `db:"sync_system" json:"syncSystem"`
	SyncID     int64  `db:"sync_id" json:"syncID"`
	SyncHash   string `db:"sync_hash" json:"syncHash"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
