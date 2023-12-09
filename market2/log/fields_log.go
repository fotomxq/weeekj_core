package Market2Log

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsLog 奖励记录
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//用户ID
	// 得到奖励的用户/推荐的人
	UserID int64 `db:"user_id" json:"userID"`
	//触发奖励的设置ID
	// 部分奖励模式下，将给与0，因为这些奖励没有具体的设置
	// 只有具体设置才会给与值
	BindID int64 `db:"bind_id" json:"bindID"`
	//触发的用户
	BindUserID int64 `db:"bind_user_id" json:"bindUserID"`
	//奖励积分
	GivingUserIntegral int64 `db:"giving_user_integral" json:"givingUserIntegral"`
	//奖励资金
	// savings 储蓄; deposit 押金; free 免费资金
	GivingDepositType  string `db:"giving_deposit_type" json:"givingDepositType"`
	GivingDepositPrice int64  `db:"giving_deposit_price" json:"givingDepositPrice"`
	//奖励票据
	GivingTicketConfigID int64 `db:"giving_ticket_config_id" json:"givingTicketConfigID"`
	GivingTicketCount    int64 `db:"giving_ticket_count" json:"givingTicketCount"`
	//奖励会员
	GivingUserSubAddHour int64 `db:"giving_user_sub_add_hour" json:"givingUserSubAddHour"`
	//行为范畴
	// 聚合统计中，按照本行为的列队行为具体定义
	// 0 new_user 新用户奖励; 1 referrer_new_user 邀请新用户奖励; 2 qrcode 扫码奖励; 3 user_sub 用户会员奖励; 4 referrer_user_sub 推荐用户会员奖励
	Action string `db:"action" json:"action"`
	//奖励原因描述
	Des string `db:"des" json:"des"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
