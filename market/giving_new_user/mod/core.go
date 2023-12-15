package MarketGivingNewUserMod

import CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"

func PushNewUserBuy(userID int64, orderPrice int64, isOrder bool) {
	CoreNats.PushDataNoErr("/market_giving/new_user/buy", "", userID, "", map[string]interface{}{
		"orderPrice": orderPrice,
		"isOrder":    isOrder,
	})
}
