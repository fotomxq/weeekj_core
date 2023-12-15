package MarketGivingQrcode

import (
	"errors"
	"fmt"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	MarketGivingCore "github.com/fotomxq/weeekj_core/v5/market/giving_core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceOrderMod "github.com/fotomxq/weeekj_core/v5/service/order/mod"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
)

// ArgsCreateLog 创建新的请求参数
type ArgsCreateLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//触发的条件ID
	// 二维码比较特殊，需指定ID才能触发
	QrcodeID int64 `db:"qrcode_id" json:"qrcodeID" check:"id"`
	//奖励的用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//推荐人用户ID
	ReferrerUserID int64 `db:"referrer_user_id" json:"referrerUserID" check:"id" empty:"true"`
	//推荐成员ID
	ReferrerBindID int64 `db:"referrer_bind_id" json:"referrerBindID" check:"id" empty:"true"`
	//交易的金额
	// 用户发生交易的总金额
	PriceTotal int64 `db:"price_total" json:"priceTotal" check:"price" empty:"true"`
}

// CreateLog 创建新的请求
func CreateLog(args *ArgsCreateLog) (errCode string, err error) {
	//检查符合的条件
	var conditionsData FieldsConditions
	err = Router2SystemConfig.MainDB.Get(&conditionsData, "SELECT id, config_id, have_phone, have_order, params FROM market_giving_qrcode WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND id = $2", args.OrgID, args.QrcodeID)
	if err != nil || conditionsData.ID < 1 {
		//没有符合条件则退出
		err = nil
		return
	}
	//检查用户的是否符合奖励条件
	var vUserInfo UserCore.FieldsUserType
	vUserInfo, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    args.UserID,
		OrgID: -1,
	})
	if err != nil {
		errCode = "no_user"
		err = errors.New("no user")
		return
	}
	if conditionsData.HavePhone {
		if vUserInfo.Phone == "" {
			errCode = "no_phone"
			err = errors.New("no phone")
			return
		}
	}
	if conditionsData.HaveOrder {
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
			errCode = "no_order"
			err = errors.New("no order")
			return
		}
	}
	//触发赠礼
	_, errCode, err = MarketGivingCore.CreateLog(&MarketGivingCore.ArgsCreateLog{
		OrgID: args.OrgID,
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "qrcode",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		UserID:         args.UserID,
		ReferrerUserID: args.ReferrerUserID,
		ReferrerBindID: args.ReferrerBindID,
		ConfigID:       conditionsData.ConfigID,
		PriceTotal:     args.PriceTotal,
		Des:            fmt.Sprint("用户扫码奖励"),
	})
	if err != nil {
		if errCode == "config_limit" {
			return
		}
		return
	}
	//反馈
	return
}
