package ServiceOrderMod

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetByID 获取订单ID参数
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
}

// GetByID 获取订单ID
func GetByID(args *ArgsGetByID) (data FieldsOrder, err error) {
	data = getByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.UserID, data.UserID) {
		data = FieldsOrder{}
		err = errors.New("no data")
		return
	}
	return
}

// GetByIDNoErr 无错误获取订单信息
func GetByIDNoErr(orderID int64) (data FieldsOrder) {
	data = getByID(orderID)
	return
}

// 获取订单
func getByID(orderID int64) (data FieldsOrder) {
	cacheMark := getOrderCacheMark(orderID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, system_mark, org_id, user_id, company_id, create_from, serial_number, serial_number_day, status, refund_status, refund_way, refund_des, refund_file_ids, refund_have_good, refund_transport_system, refund_transport_sn, refund_transport_info, refund_pay_id, refund_price, refund_pay_finish, refund_expire_at, refund_tip_at, address_from, address_to, goods, exemptions, allow_auto_audit, transport_id, transport_allow_auto, transport_task_at, transport_pay_after, transport_ids, transport_system, transport_sn, transport_info, transport_status, price_list, price_pay, currency, price, price_total, pay_status, pay_id, pay_list, pay_from, des, logs, params FROM service_order WHERE id = $1", orderID)
	if err != nil || data.ID < 1 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, system_mark, org_id, user_id, company_id, create_from, serial_number, serial_number_day, status, refund_status, refund_way, refund_des, refund_file_ids, refund_have_good, refund_transport_system, refund_transport_sn, refund_transport_info, refund_pay_id, refund_price, refund_pay_finish, refund_expire_at, refund_tip_at, address_from, address_to, goods, exemptions, allow_auto_audit, transport_id, transport_allow_auto, transport_task_at, transport_pay_after, transport_ids, transport_system, transport_sn, transport_info, transport_status, price_list, price_pay, currency, price, price_total, pay_status, pay_id, pay_list, pay_from, des, logs, params FROM service_order_history WHERE id = $1", orderID)
		if err != nil || data.ID < 1 {
			err = errors.New("no data")
			return
		}
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}
