package MarketGivingNewUser

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//新注册用户
	CoreNats.SubDataByteNoErr("/user/login2/new", subNatsUserLogin2New)
	//新用户注册后购买行为
	CoreNats.SubDataByteNoErr("/market_giving/new_user/buy", subNatsNewUserBuy)
}

// 新注册用户
func subNatsUserLogin2New(_ *nats.Msg, action string, userID int64, _ string, data []byte) {
	if action != "new" {
		return
	}
	//获取参数
	orgID := gjson.GetBytes(data, "orgID").Int()
	referrerNationCode := gjson.GetBytes(data, "referrerNationCode").String()
	referrerPhone := gjson.GetBytes(data, "referrerPhone").String()
	//如果referrerNationCode为空，则默认采用86
	if referrerNationCode == "" {
		referrerNationCode = "86"
	}
	//查询推荐人
	var referrerUserID, referrerBindID int64
	userData, err := UserCore.GetUserByPhone(&UserCore.ArgsGetUserByPhone{
		OrgID:      orgID,
		NationCode: referrerNationCode,
		Phone:      referrerPhone,
	})
	if err == nil {
		referrerUserID = userData.ID
		bindData, err := OrgCore.GetBindByUserAndOrg(&OrgCore.ArgsGetBindByUserAndOrg{
			UserID: userData.ID,
			OrgID:  orgID,
		})
		if err == nil {
			referrerBindID = bindData.ID
		}
	}
	//开始赠礼
	errCode, err := createLog(&argsCreateLog{
		OrgID:          orgID,
		UserID:         userID,
		ReferrerUserID: referrerUserID,
		ReferrerBindID: referrerBindID,
		PriceTotal:     0,
		IsOrder:        false,
	})
	if err != nil {
		CoreLog.Warn("market giving new user sub nats, user login2 new, create log failed: ", errCode, ", err: ", err)
	}
}

// 新用户注册后购买行为
func subNatsNewUserBuy(_ *nats.Msg, _ string, userID int64, _ string, data []byte) {
	appendLog := "market giving new user sub nats, new user buy, "
	//处理新注册用户赠礼
	userData, err := UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    userID,
		OrgID: -1,
	})
	if err != nil {
		CoreLog.Warn(appendLog, "get user by id: ", userID, ", err: ", err)
		return
	}
	regReferrerUserID, _ := userData.Infos.GetValInt64("referrerUserID")
	regReferrerBindID, _ := userData.Infos.GetValInt64("referrerBindID")
	//获取参数
	orderPrice := gjson.GetBytes(data, "orderPrice").Int()
	isOrder := gjson.GetBytes(data, "isOrder").Bool()
	//触发奖励
	errCode, err := createLog(&argsCreateLog{
		OrgID:          userData.OrgID,
		UserID:         userData.ID,
		ReferrerUserID: regReferrerUserID,
		ReferrerBindID: regReferrerBindID,
		PriceTotal:     orderPrice,
		IsOrder:        isOrder,
	})
	if err != nil {
		CoreLog.Warn(appendLog, "create log failed: ", errCode, ", err: ", err)
		return
	}
}
