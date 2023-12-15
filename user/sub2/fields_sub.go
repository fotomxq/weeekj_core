package UserSub2

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsSub 订阅记录
// 配置文件记录会员的级别、对应时间长度、对应价格
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
	// 留空则为平台级别会员，否则为组织的会员
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//上次支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//上次支付状态
	PayAt time.Time `db:"pay_at" json:"payAt"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
