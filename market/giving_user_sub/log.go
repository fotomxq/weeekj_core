package MarketGivingUserSub

import (
	"fmt"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	MarketGivingCore "gitee.com/weeekj/weeekj_core/v5/market/giving_core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsCreateLog 创建新的请求参数
type ArgsCreateLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//奖励的用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//推荐人用户ID
	ReferrerUserID int64 `db:"referrer_user_id" json:"referrerUserID" check:"id" empty:"true"`
	//推荐成员ID
	ReferrerBindID int64 `db:"referrer_bind_id" json:"referrerBindID" check:"id" empty:"true"`
	//交易的金额
	// 用户发生交易的总金额
	PriceTotal int64 `db:"price_total" json:"priceTotal" check:"price" empty:"true"`
	//订阅ID
	SubConfigID int64 `db:"sub_config_id" json:"subConfigID" check:"id"`
	//订阅单位
	SubBuyCount int64 `db:"sub_buy_count" json:"subBuyCount" check:"int64Than0"`
	//锁定赠礼ID
	LockGivingUserSubID int64 `db:"lock_giving_user_sub_id" json:"lockGivingUserSubID"`
}

// CreateLog 创建新的请求
// 注意，如果多个相同订阅，则以最大符合单位个数的为基准赠送
func CreateLog(args *ArgsCreateLog) (errCode string, err error) {
	//检查符合的条件
	var conditionsData FieldsConditions
	if args.LockGivingUserSubID > 0 {
		err = Router2SystemConfig.MainDB.Get(&conditionsData, "SELECT id, config_id, params FROM market_giving_user_sub WHERE delete_at < to_timestamp(1000000) AND id = $1 ORDER BY sub_buy_count DESC LIMIT 1", args.LockGivingUserSubID)
	} else {
		err = Router2SystemConfig.MainDB.Get(&conditionsData, "SELECT id, config_id, params FROM market_giving_user_sub WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND sub_config_id = $2 AND sub_buy_count = $3 ORDER BY sub_buy_count DESC LIMIT 1", args.OrgID, args.SubConfigID, args.SubBuyCount)
	}
	if err != nil || conditionsData.ID < 1 {
		//没有符合条件则退出
		err = nil
		return
	}
	//触发赠礼
	_, errCode, err = MarketGivingCore.CreateLog(&MarketGivingCore.ArgsCreateLog{
		OrgID: args.OrgID,
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user_sub",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		UserID:         args.UserID,
		ReferrerUserID: args.ReferrerUserID,
		ReferrerBindID: args.ReferrerBindID,
		ConfigID:       conditionsData.ConfigID,
		PriceTotal:     args.PriceTotal,
		Des:            fmt.Sprint("购买会员[", args.SubConfigID, "]超过[", args.SubBuyCount, "]个单位奖励"),
	})
	if err != nil {
		if errCode == "config_limit" {
			err = nil
			return
		}
		return
	}
	//反馈
	return
}
