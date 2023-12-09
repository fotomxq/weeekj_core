package UserFocus

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsFocus 用户关注
type FieldsFocus struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//绑定组织
	// 该组织根据资源来源设定
	// 如果是平台资源，则为0
	OrgID int64 `db:"org_id" json:"orgID"`
	//关注类型
	Mark string `db:"mark" json:"mark"`
	//关注内容来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
}
