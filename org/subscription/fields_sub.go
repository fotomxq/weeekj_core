package OrgSubscription

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsSub 订阅存储主体
// 同一个组织的不同配置会有多个记录
type FieldsSub struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//开通配置
	ConfigID int64 `db:"config_id" json:"configID"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
