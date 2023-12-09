package MallCoreMod

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

// SendProductVirtual 处理虚拟商品，交付给用户
func SendProductVirtual(productID int64, count int64, userID, orgID int64, orderID int64) {
	CoreNats.PushDataNoErr("/mall/core/product_virtual", "send", productID, "", map[string]interface{}{
		"count":   count,
		"userID":  userID,
		"orgID":   orgID,
		"orderID": orderID,
	})
}
