package Market2Log

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	FinanceDeposit2 "gitee.com/weeekj/weeekj_core/v5/finance/deposit2"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserIntegral "gitee.com/weeekj/weeekj_core/v5/user/integral"
	UserSubscription "gitee.com/weeekj/weeekj_core/v5/user/subscription"
	UserTicket "gitee.com/weeekj/weeekj_core/v5/user/ticket"
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
func AppendLog(args *ArgsAppendLog) (errCode string, err error) {
	//修正参数
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//如果不存在被奖励人，将退出
	if args.UserID < 1 {
		return
	}
	//如果禁止重复奖励，则绕过处理
	if args.NoReplaceByFrom {
		logData := GetLastLogByFrom(args.Action, args.BindID, args.OrgID, args.OrgBindID, args.UserID, args.BindUserID)
		if logData.ID > 0 {
			errCode = "err_replace"
			err = errors.New("have replace data")
			return
		}
	}
	//找到用户关联的组织成员
	if args.UserID > 0 && args.OrgBindID < 1 && args.OrgID > 0 {
		bindData := OrgCore.GetBindByUserIDOnly(args.OrgID, args.UserID)
		if bindData.ID > 0 {
			args.OrgBindID = bindData.ID
		}
	}
	if args.OrgBindID < 1 {
		args.OrgBindID = 0
	}
	//是否存在奖励
	haveGivingData := false
	//奖励对应资源
	if args.GivingUserIntegral > 0 {
		_ = UserIntegral.AddCount(&UserIntegral.ArgsAddCount{
			OrgID:    args.SourceOrgID,
			UserID:   args.UserID,
			AddCount: args.GivingUserIntegral,
			Des:      args.Des,
		})
		haveGivingData = true
	}
	if args.GivingDepositPrice > 0 {
		// savings 储蓄; deposit 押金; free 免费资金
		switch args.GivingDepositType {
		case "savings":
			_, _ = FinanceDeposit2.SetUserSaving("", args.SourceOrgID, args.UserID, args.GivingDepositPrice)
		case "deposit":
			_, _ = FinanceDeposit2.SetUserDeposit("", args.SourceOrgID, args.UserID, args.GivingDepositPrice)
		case "free":
			_, _ = FinanceDeposit2.SetUserFree("", args.SourceOrgID, args.UserID, args.GivingDepositPrice)
		}
		haveGivingData = true
	}
	if args.GivingTicketConfigID > 0 && args.GivingTicketCount > 0 {
		_ = UserTicket.AddTicket(&UserTicket.ArgsAddTicket{
			OrgID:       args.SourceOrgID,
			ConfigID:    args.GivingTicketConfigID,
			UserID:      args.UserID,
			Count:       args.GivingTicketCount,
			UseFromName: "奖励",
		})
		haveGivingData = true
	}
	if args.GivingUserSubAddHour > 0 {
		//计算为秒单位
		addSec := args.GivingUserSubAddHour * 60 * 60
		_ = UserSubscription.AddUserSubAny(args.SourceOrgID, args.UserID, addSec)
		haveGivingData = true
	}
	//如果没有奖励，则退出
	if !haveGivingData {
		return
	}
	//写入记录
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO market2_log (org_id, org_bind_id, user_id, bind_id, bind_user_id, giving_user_integral, giving_deposit_type, giving_deposit_price, giving_ticket_config_id, giving_ticket_count, giving_user_sub_add_hour, action, des, params) VALUES (:org_id,:org_bind_id,:user_id,:bind_id,:bind_user_id,:giving_user_integral,:giving_deposit_type,:giving_deposit_price,:giving_ticket_config_id,:giving_ticket_count,:giving_user_sub_add_hour,:action,:des,:params)", map[string]interface{}{
		"org_id":                   args.OrgID,
		"org_bind_id":              args.OrgBindID,
		"user_id":                  args.UserID,
		"bind_id":                  args.BindID,
		"bind_user_id":             args.BindUserID,
		"giving_user_integral":     args.GivingUserIntegral,
		"giving_deposit_type":      args.GivingDepositType,
		"giving_deposit_price":     args.GivingDepositPrice,
		"giving_ticket_config_id":  args.GivingTicketConfigID,
		"giving_ticket_count":      args.GivingTicketCount,
		"giving_user_sub_add_hour": args.GivingUserSubAddHour,
		"action":                   args.Action,
		"des":                      args.Des,
		"params":                   args.Params,
	})
	if err != nil {
		return
	}
	//反馈
	return
}
