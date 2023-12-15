package MarketCore

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/lib/pq"
	"time"
)

// FieldsLog 营销记录
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//对表成员的用户ID
	// 和成员对等，可用于一次性推荐的记录处理
	UserID int64 `db:"user_id" json:"userID"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//绑定的用户
	BindUserID int64 `db:"bind_user_id" json:"bindUserID"`
	//绑定的档案
	BindInfoID int64 `db:"bind_info_id" json:"bindInfoID"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//领取积分
	UserIntegral int64 `db:"user_integral" json:"userIntegral"`
	//领取用户订阅
	UserSubs FieldsConfigUserSubs `db:"user_subs" json:"userSubs"`
	//领取票据
	UserTickets FieldsConfigUserTickets `db:"user_tickets" json:"userTickets"`
	//奖励的来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//交易的金额
	// 用户发生交易的总金额
	PriceTotal int64 `db:"price_total" json:"priceTotal"`
	//奖励金储蓄标识码
	DepositConfigMark string `db:"deposit_config_mark" json:"depositConfigMark"`
	//奖励金额
	Price int64 `db:"price" json:"price"`
	//奖励的次数
	// 部分系统中不一定非要奖励金额，可能只是次数
	Count int64 `db:"count" json:"count"`
	//奖励依据配置
	ConfigID int64 `db:"config_id" json:"configID"`
	//奖励原因描述
	// eg: 推荐用户注册新用户 / 推荐用户购买商品
	Des string `db:"des" json:"des"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
