package MarketGivingCore

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//奖励的来源
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//发生用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//推荐人用户ID
	ReferrerUserID int64 `db:"referrer_user_id" json:"referrerUserID"`
	//推荐成员ID
	ReferrerBindID int64 `db:"referrer_bind_id" json:"referrerBindID"`
	//赠送配置
	ConfigID int64 `db:"config_id" json:"configID"`
	//领取积分
	UserIntegral int64 `db:"user_integral" json:"userIntegral"`
	//领取用户订阅
	UserSubs FieldsConfigUserSubs `db:"user_subs" json:"userSubs"`
	//领取票据
	UserTickets FieldsConfigUserTickets `db:"user_tickets" json:"userTickets"`
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
	//描述
	Des string `db:"des" json:"des"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
