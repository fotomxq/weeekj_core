package ServiceOrderMod

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
)

// UpdateOrderParams 请求更新订单的扩展参数
func UpdateOrderParams(orderID int64, params []CoreSQLConfig.FieldsConfigType) {
	CoreNats.PushDataNoErr("/service/order/params", "", orderID, "", map[string]interface{}{
		"params": params,
	})
}
