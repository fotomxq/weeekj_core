package Market2LogMod

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
)

// ArgsAppendLog 添加新的日志参数
type ArgsAppendLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	// 允许给0，系统将自动根据orgID和用户ID找到匹配的组织成员
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//用户ID
	// 得到奖励的用户
	UserID int64 `db:"user_id" json:"userID"`
	//触发奖励的设置ID
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
	//同一个行为禁止重复
	// 同一个被奖励来源和触发奖励来源，将被禁止触发奖励
	// 如果是存在时间限制为主，请在外围具体实施模块加以判断处理后提交给本方法，同时禁止使用此参数，避免永久性无法触发
	NoReplaceByFrom bool `json:"noReplaceByFrom"`
	//资源导向的来源组织
	// 相关奖励资源的让渡方
	SourceOrgID int64 `json:"sourceOrgID"`
}

// AppendLog 添加新的日志
func AppendLog(args ArgsAppendLog) {
	CoreNats.PushDataNoErr("/market2_log/create", "", 0, "", args)
}
