package UserChat

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsGroup 聊天室
type FieldsGroup struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//最后一次聊天时间
	LastAt time.Time `db:"last_at" json:"lastAt"`
	//绑定组织
	// 商户可以查看构建的相关聊天室
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//聊天室创建人
	UserID int64 `db:"user_id" json:"userID"`
	//只有创建人能邀请其他人？
	OnlyCreateInvite bool `db:"only_create_invite" json:"onlyCreateInvite"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
