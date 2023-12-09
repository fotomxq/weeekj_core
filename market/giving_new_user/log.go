package MarketGivingNewUser

import (
	"errors"
	"fmt"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	MarketGivingCore "gitee.com/weeekj/weeekj_core/v5/market/giving_core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "gitee.com/weeekj/weeekj_core/v5/service/order/mod"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
)

// argsCreateLog 创建新的请求参数
type argsCreateLog struct {
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
	//是否为订单发生时
	IsOrder bool `json:"isOrder"`
}

// createLog 创建新的请求
func createLog(args *argsCreateLog) (errCode string, err error) {
	//检查符合的条件
	var conditionsList []FieldsConditions
	err = Router2SystemConfig.MainDB.Select(&conditionsList, "SELECT id, config_id, have_phone, after_sign, before_sign, have_order, params FROM market_giving_new_user WHERE delete_at < to_timestamp(1000000) AND org_id = $1", args.OrgID)
	if err != nil || len(conditionsList) < 1 {
		//没有符合条件则退出
		err = nil
		return
	}
	//遍历条件并触发
	for _, v := range conditionsList {
		//是否为订单发生时，必须是订单后才发生该奖励
		if args.IsOrder && !v.HaveOrder {
			continue
		}
		//检查用户的是否符合奖励条件
		var vUserInfo UserCore.FieldsUserType
		vUserInfo, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
			ID:    args.UserID,
			OrgID: -1,
		})
		if err != nil {
			err = nil
			continue
		}
		if v.HavePhone {
			if vUserInfo.Phone == "" {
				continue
			}
		}
		if v.AfterSign.Unix() > 1000000 {
			if vUserInfo.CreateAt.Unix() < v.AfterSign.Unix() {
				continue
			}
		}
		if v.BeforeSign.Unix() > 1000000 {
			if vUserInfo.CreateAt.Unix() > v.BeforeSign.Unix() {
				continue
			}
		}
		if v.HaveOrder {
			var orderCount int64
			orderCount, _ = ServiceOrderMod.GetUserOrderCount(&ServiceOrderMod.ArgsGetUserOrderCount{
				OrgID:    args.OrgID,
				UserID:   args.UserID,
				IsFinish: true,
			})
			if orderCount < 1 {
				err = UserCore.UpdateUserInfosByID(&UserCore.ArgsUpdateUserInfosByID{
					ID:       vUserInfo.ID,
					OrgID:    vUserInfo.OrgID,
					Mark:     "referrerUserID",
					Val:      fmt.Sprint(args.ReferrerUserID),
					IsRemove: false,
				})
				if err != nil {
					errCode = "update_user_info"
					err = errors.New(fmt.Sprint("update user info, ", err))
					return
				}
				err = UserCore.UpdateUserInfosByID(&UserCore.ArgsUpdateUserInfosByID{
					ID:       vUserInfo.ID,
					OrgID:    vUserInfo.OrgID,
					Mark:     "referrerBindID",
					Val:      fmt.Sprint(args.ReferrerBindID),
					IsRemove: false,
				})
				if err != nil {
					errCode = "update_user_info"
					err = errors.New(fmt.Sprint("update user info, ", err))
					return
				}
				continue
			}
		}
		//触发赠礼
		_, errCode, err = MarketGivingCore.CreateLog(&MarketGivingCore.ArgsCreateLog{
			OrgID: args.OrgID,
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "new_user",
				ID:     0,
				Mark:   "",
				Name:   "",
			},
			UserID:         args.UserID,
			ReferrerUserID: args.ReferrerUserID,
			ReferrerBindID: args.ReferrerBindID,
			ConfigID:       v.ConfigID,
			PriceTotal:     args.PriceTotal,
			Des:            fmt.Sprint("新注册用户奖励"),
		})
		if err != nil {
			if errCode == "config_limit" {
				err = nil
				continue
			}
			return
		}
	}
	//反馈
	return
}
