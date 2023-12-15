package ServiceOrder

import (
	"encoding/json"
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetList 获取订单列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//来源系统
	// 该订单创建来源的系统
	// eg: user_sub / org_sub / mall / core_api
	SystemMark string `db:"system_mark" json:"systemMark" check:"mark" empty:"true"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// -1 跳过
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//创建来源
	// 0 商户平台下单、1 用户APP下单、2 微信小程序下单、3 支付宝小程序下单、4 电话系统下单、5 线下柜台下单、6 辅助设备下单、7 其他渠道
	CreateFrom int `db:"create_from" json:"createFrom"`
	//状态
	// 0 草稿等待提交
	// 1 提交等待审核中
	// 2 送货中，内部状态根据配送状态确认
	// 3 送货完成，可能包含货到付款
	// 4 送货完成且付款完成
	// 5 订单失败，发货失败等因素
	// 6 取消，包括超时、人为因素
	Status []int `db:"status" json:"status"`
	//退货状态
	// 0 没有退货申请
	// 1 提交退货申请
	// 2 退货中
	// 3 退货完成，退款需配合pay_status进行
	RefundStatus []int `db:"refund_status" json:"refundStatus"`
	//配送ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id" empty:"true"`
	//允许自动配送
	NeedTransportAllowAuto bool `json:"needTransportAllowAuto" check:"bool"`
	TransportAllowAuto     bool `db:"transport_allow_auto" json:"transportAllowAuto" check:"bool"`
	//付费状态
	// 0 尚未付款
	// 1 已经付款
	// 2 发起退款
	// 3 完成退款
	PayStatus []int `db:"pay_status" json:"payStatus"`
	//当前匹配的支付ID
	PayID int64 `db:"pay_id" json:"payID" check:"id" empty:"true"`
	//支付渠道
	PayFrom string `db:"pay_from" json:"payFrom"`
	//货物清单
	GoodFrom CoreSQLFrom.FieldsFrom `db:"good_from" json:"goodFrom"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime2.DataCoreTime `json:"timeBetween"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//是否为归档订单
	IsHistory bool `db:"is_history" json:"isHistory" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取订单列表
func GetList(args *ArgsGetList) (dataList []FieldsOrder, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.SystemMark != "" {
		where = where + " AND system_mark = :system_mark"
		maps["system_mark"] = args.SystemMark
	}
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.CompanyID > -1 {
		where = where + " AND company_id = :company_id"
		maps["company_id"] = args.CompanyID
	}
	if args.CreateFrom > -1 {
		where = where + " AND create_from = :create_from"
		maps["create_from"] = args.CreateFrom
	}
	if len(args.Status) > 0 {
		where = where + " AND status = ANY(:status)"
		maps["status"] = pq.Array(args.Status)
	}
	if len(args.RefundStatus) > 0 {
		where = where + " AND refund_status = ANY(:refund_status)"
		maps["refund_status"] = pq.Array(args.RefundStatus)
	}
	if args.TransportID > -1 {
		where = where + " AND transport_id = :transport_id"
		maps["transport_id"] = args.TransportID
	}
	if args.NeedTransportAllowAuto {
		where = where + " AND transport_allow_auto = :transport_allow_auto"
		maps["transport_allow_auto"] = args.TransportAllowAuto
	}
	if len(args.PayStatus) > 0 {
		where = where + " AND pay_status = ANY(:pay_status)"
		maps["pay_status"] = pq.Array(args.PayStatus)
	}
	if args.PayID > -1 {
		where = where + " AND pay_id = :pay_id"
		maps["pay_id"] = args.PayID
	}
	if args.PayFrom != "" {
		where = where + " AND pay_from = :pay_from"
		maps["pay_from"] = args.PayFrom
	}
	if args.GoodFrom.System != "" {
		where, maps, err = args.GoodFrom.GetListAnd("goods", "goods", where, maps)
		if err != nil {
			return
		}
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		where = where + " AND create_at >= :start_at AND create_at <= :end_at"
		maps["start_at"] = args.TimeBetween.MinTime
		maps["end_at"] = args.TimeBetween.MaxTime
	}
	if args.Search != "" {
		where = where + " AND (address_from ->> 'address' ILIKE '%' || :search || '%' OR address_from ->> 'name' ILIKE '%' || :search || '%' OR address_from ->> 'phone' ILIKE '%' || :search || '%' OR address_to ->> 'address' ILIKE '%' || :search || '%' OR address_to ->> 'name' ILIKE '%' || :search || '%' OR address_to ->> 'phone' ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%' OR logs ->> 'des' ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_order"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	var rawList []FieldsOrder
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil || len(rawList) < 1 {
		return
	}
	for _, v := range rawList {
		vData := getByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

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
	if data.ID < 1 {
		return
	}
	err = checkOrgAndUser(&data, args.OrgID, args.UserID)
	if err != nil {
		return
	}
	return
}

func GetByIDNoErr(id int64, orgID int64, userID int64) (data FieldsOrder) {
	data = getByID(id)
	if data.ID < 1 {
		return
	}
	err := checkOrgAndUser(&data, orgID, userID)
	if err != nil {
		return
	}
	return
}

// GetNearDayCountByOrgID 获取近期增加的下单客户人数
func GetNearDayCountByOrgID(orgID int64, subDay int) (count int64) {
	findAt := CoreFilter.GetNowTimeCarbon().SubDays(subDay)
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_order WHERE org_id = $1 AND create_at >= $2", orgID, findAt.Time)
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

// 获取一组订单
func getByIDs(ids pq.Int64Array) (dataList []FieldsOrder, err error) {
	for _, v := range ids {
		var data FieldsOrder
		data = getByID(v)
		if data.ID < 1 {
			continue
		}
		dataList = append(dataList, data)
	}
	return
}

// 检查订单的所属权
func checkOrgAndUser(orderData *FieldsOrder, orgID, userID int64) (err error) {
	if CoreFilter.EqID2(orgID, orderData.OrgID) && CoreFilter.EqID2(userID, orderData.UserID) {
		return
	}
	err = errors.New("no data")
	return
}

// 根据支付ID获取订单列表
func getListByPayID(payID int64) (dataList []FieldsOrder, err error) {
	var rawList []FieldsOrder
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_order WHERE pay_id = $1 OR ($1 = ANY(pay_list)) AND delete_at < to_timestamp(1000000)", payID)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// 获取符合配送单ID的订单列表
func getListByTransportID(transportSystem string, transportID int64) (dataList []FieldsOrder, err error) {
	var rawList []FieldsOrder
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM service_order WHERE transport_system = $1 AND transport_id = $2 AND delete_at < to_timestamp(1000000)", transportSystem, transportID)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// 生成追加日志的部分
func getLogData(userID int64, orgBindID int64, mark, des string) (logData string, err error) {
	newLog := []FieldsLog{
		{
			CreateAt:  CoreFilter.GetNowTime(),
			UserID:    userID,
			OrgBindID: orgBindID,
			Mark:      mark,
			Des:       des,
		},
	}
	var newLogByte []byte
	newLogByte, err = json.Marshal(newLog)
	if err != nil {
		return
	}
	logData = string(newLogByte)
	return
}
