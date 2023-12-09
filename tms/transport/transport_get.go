package TMSTransport

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetTransportList 获取配送列表参数
type ArgsGetTransportList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//当前配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//编号
	// 商户下唯一
	SN int64 `db:"sn" json:"sn" check:"int64Than0" empty:"true"`
	//今日编号
	SNDay int64 `db:"sn_day" json:"snDay" check:"int64Than0" empty:"true"`
	//符合条件的一组配送状态
	// 0 等待分配人员; 1 取货中; 2 送货中; 3 完成配送
	Status pq.Int32Array `db:"status" json:"status"`
	//是否需要支付参数
	NeedIsPay bool `db:"need_is_pay" json:"needIsPay" check:"bool"`
	//是否已经支付配送
	IsPay bool `db:"is_pay" json:"isPay" check:"bool"`
	//缴费交易ID
	PayID int64 `db:"pay_id" json:"payID" check:"id" empty:"true"`
	//是否为完成时间的时间范围
	IsFinishAt bool `json:"isFinishAt" check:"bool"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//是否为历史
	IsHistory bool `db:"is_history" json:"isHistory" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTransportList 获取配送列表
func GetTransportList(args *ArgsGetTransportList) (dataList []FieldsTransport, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.InfoID > -1 {
		where = where + " AND info_id = :info_id"
		maps["info_id"] = args.InfoID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.OrderID > -1 {
		where = where + " AND order_id = :order_id"
		maps["order_id"] = args.OrderID
	}
	if args.SN > 0 {
		where = where + " AND sn = :sn"
		maps["sn"] = args.SN
	}
	if args.SNDay > 0 {
		where = where + " AND sn_day = :sn_day"
		maps["sn_day"] = args.SNDay
	}
	if len(args.Status) > 0 {
		where = where + " AND status = ANY(:status)"
		maps["status"] = args.Status
	}
	if args.NeedIsPay {
		if args.IsPay {
			where = where + " AND pay_finish_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND pay_finish_at < to_timestamp(1000000)"
		}
	}
	if args.PayID > -1 {
		where = where + " AND pay_id = :pay_id"
		maps["pay_id"] = args.PayID
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		if args.IsFinishAt {
			where = where + " AND finish_at >= :start_at AND finish_at <= :end_at"
		} else {
			where = where + " AND create_at >= :start_at AND create_at <= :end_at"
		}
		maps["start_at"] = timeBetween.MinTime
		maps["end_at"] = timeBetween.MaxTime
	}
	if args.Search != "" {
		where = where + " AND (from_address ->> 'address' ILIKE '%' || :search || '%' OR from_address ->> 'name' ILIKE '%' || :search || '%' OR from_address ->> 'phone' ILIKE '%' || :search || '%' OR to_address ->> 'address' ILIKE '%' || :search || '%' OR to_address ->> 'name' ILIKE '%' || :search || '%' OR to_address ->> 'phone' ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "tms_transport"
	if args.IsHistory {
		tableName = "tms_transport_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, finish_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, pay_ids, task_at, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "finish_at"},
	)
	return
}

// ArgsGetTransport 获取配送信息参数
type ArgsGetTransport struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//客户档案ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//客户用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetTransport 获取配送信息
func GetTransport(args *ArgsGetTransport) (data FieldsTransport, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, finish_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, pay_ids, task_at, params FROM tms_transport WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR info_id = $3) AND ($4 < 1 OR user_id = $4)", args.ID, args.OrgID, args.InfoID, args.UserID)
	if err != nil {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, finish_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, pay_ids, task_at, params FROM tms_transport_history WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR info_id = $3) AND ($4 < 1 OR user_id = $4)", args.ID, args.OrgID, args.InfoID, args.UserID)
		return
	}
	return
}

func getTransportByID(id int64) (data FieldsTransport, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, finish_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, pay_ids, task_at, params FROM tms_transport WHERE id = $1", id)
	return
}

// ArgsGetTransports 查询多个配送单参数
type ArgsGetTransports struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetTransports 查询多个配送单
func GetTransports(args *ArgsGetTransports) (dataList []FieldsTransport, err error) {
	err = CoreSQLIDs.GetIDsOrgAndDelete(&dataList, "tms_transport", "id, create_at, update_at, delete_at, finish_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, pay_ids, task_at, params", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// ArgsCheckTransport 检查ID列最近是否更新过参数
type ArgsCheckTransport struct {
	//检查的ID列
	Data []ArgsCheckTransportData `json:"data"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//当前配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

type ArgsCheckTransportData struct {
	//ID
	ID int64 `json:"id"`
	//更新时间
	UpdateAt string `json:"updateAt"`
}

// CheckTransport 检查ID列最近是否更新过
// 如果存在更新，才会反馈数据
func CheckTransport(args *ArgsCheckTransport) (updateIDs []int64) {
	ids := pq.Int64Array{}
	for _, v := range args.Data {
		ids = append(ids, v.ID)
	}
	if len(ids) < 1 {
		return
	}
	var dataList []FieldsTransport
	err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, update_at FROM tms_transport WHERE id = ANY($1) AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000) AND ($3 < 1 OR bind_id = $3)", ids, args.OrgID, args.BindID)
	if err != nil {
		return
	}
	for _, v := range dataList {
		for _, v2 := range args.Data {
			if v.ID != v2.ID {
				continue
			}
			updateAt, err := CoreFilter.GetTimeByISO(v2.UpdateAt)
			if err != nil {
				continue
			}
			if v.UpdateAt.Unix() > updateAt.Unix() {
				updateIDs = append(updateIDs, v.ID)
			}
		}
	}
	return
}

// 获取个别核心信息
func getTransportOrgAndBindByID(id int64) (data FieldsTransport, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id, bind_id FROM tms_transport WHERE id = $1", id)
	return
}

// 获取支付ID的所有配送
func getTransportOrgAndBindByPayID(payID int64) (dataList []FieldsTransport, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, finish_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, pay_ids, task_at, params FROM tms_transport WHERE pay_id = $1 OR pay_id = ANY(pay_ids)", payID)
	return
}

// 批量获取配送单ID列
func getTransportIDs(ids pq.Int64Array) (dataList []FieldsTransport, err error) {
	if err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, org_id, bind_id FROM tms_transport WHERE id = ANY($1) AND delete_at < to_timestamp(1000000)", ids); err != nil {
		return
	}
	return
}

// 获取订单关联的配送单数量
func getTransportCountByOrderID(orderID int64) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM tms_transport WHERE order_id = $1 AND delete_at < to_timestamp(1000000)", orderID)
	return
}

// 获取订单关联的配送单
func getTransportByOrderID(orderID int64) (dataList []FieldsTransport, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, finish_at, org_id, bind_id, info_id, user_id, sn, sn_day, status, from_address, to_address, order_id, goods, weight, length, width, currency, price, pay_finish_at, pay_id, pay_ids, task_at, params FROM tms_transport WHERE order_id = $1 AND delete_at < to_timestamp(1000000)", orderID)
	return
}
