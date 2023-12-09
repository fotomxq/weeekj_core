package ServiceOrderWait

import (
	"errors"
	"fmt"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceOrderWaitFields "gitee.com/weeekj/weeekj_core/v5/service/order/wait_fields"
)

// 获取等待订单
func getCreateWait(id int64) (data ServiceOrderWaitFields.FieldsWait, err error) {
	cacheMark := fmt.Sprint("service:order:wait:id:", id)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, order_id, system_mark, org_id, user_id, create_from, hash, address_from, address_to, goods, exemptions, allow_auto_audit, transport_allow_auto, transport_task_at, transport_pay_after, price_list, price_pay, currency, price, price_total, des, logs, params, err_code, err_msg, transport_system FROM service_order_wait WHERE id = $1", id)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 3600)
	return
}

// 删除等待订单ID
func deleteWaitCache(id int64) {
	cacheMark := fmt.Sprint("service:order:wait:id:", id)
	Router2SystemConfig.MainCache.DeleteMark(cacheMark)
}
