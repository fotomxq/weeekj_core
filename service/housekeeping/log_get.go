package ServiceHousekeeping

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//成员ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//是否已经上门
	NeedIsNeed bool `json:"needIsNeed"`
	IsNeed     bool `json:"isNeed"`
	//是否完成
	NeedIsFinish bool `json:"needIsFinish"`
	IsFinish     bool `json:"isFinish"`
	//是否支付
	NeedIsPay bool `json:"needIsPay"`
	IsPay     bool `json:"isPay"`
	//编号
	// 商户下唯一
	SN int64 `db:"sn" json:"sn" check:"int64Than0" empty:"true"`
	//今日编号
	SNDay int64 `db:"sn_day" json:"snDay" check:"int64Than0" empty:"true"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		where = where + " AND (bind_id = :bind_id OR :bind_id = ANY(other_binds))"
		maps["bind_id"] = args.BindID
	}
	if args.OrderID > -1 {
		where = where + " AND order_id = :order_id"
		maps["order_id"] = args.OrderID
	}
	if args.NeedIsNeed {
		if args.IsNeed {
			where = where + " AND need_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND need_at < to_timestamp(1000000)"
		}
	}
	if args.NeedIsFinish {
		if args.IsFinish {
			where = where + " AND finish_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND finish_at < to_timestamp(1000000)"
		}
	}
	if args.NeedIsPay {
		if args.IsPay {
			where = where + " AND pay_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND pay_at < to_timestamp(1000000)"
		}
	}
	if args.SN > 0 {
		where = where + " AND sn = :sn"
		maps["sn"] = args.SN
	}
	if args.SNDay > 0 {
		where = where + " AND sn_day = :sn_day"
		maps["sn_day"] = args.SNDay
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		where = where + " AND create_at >= :start_at AND create_at <= :end_at"
		maps["start_at"] = timeBetween.MinTime
		maps["end_at"] = timeBetween.MaxTime
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%' OR address::text ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_housekeeping_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, sn, sn_day, need_at, finish_at, user_id, org_id, bind_id, other_binds, mall_product_id, order_id, currency, price, pay_id, pay_at, des, address, config_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "need_at", "finish_at", "pay_at"},
	)
	return
}

// ArgsGetLogID 获取数据参数
type ArgsGetLogID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// GetLogID 获取数据
func GetLogID(args *ArgsGetLogID) (data FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, sn, sn_day, need_at, finish_at, user_id, org_id, bind_id, other_binds, mall_product_id, order_id, currency, price, pay_id, pay_at, des, address, config_id, params FROM service_housekeeping_log WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR (bind_id = $3 OR $3 = ANY(other_binds)))", args.ID, args.OrgID, args.BindID)
	return
}

// 获取服务单
func getLogID(id int64) (data FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, sn, sn_day, need_at, finish_at, user_id, org_id, bind_id, other_binds, mall_product_id, order_id, currency, price, pay_id, pay_at, des, address, config_id, params FROM service_housekeeping_log WHERE id = $1", id)
	return
}

// ArgsGetLogByUserID 获取用户的ID参数
type ArgsGetLogByUserID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetLogByUserID 获取用户的ID
func GetLogByUserID(args *ArgsGetLogByUserID) (data FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, sn, sn_day, need_at, finish_at, user_id, org_id, bind_id, other_binds, mall_product_id, order_id, currency, price, pay_id, pay_at, des, address, config_id, params FROM service_housekeeping_log WHERE id = $1 AND user_id = $2", args.ID, args.UserID)
	return
}

// ArgsGetLogByOrderIDs 通过订单ID列查询服务单参数
type ArgsGetLogByOrderIDs struct {
	//订单ID列
	OrderIDs pq.Int64Array `db:"order_ids" json:"orderIDs" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//服务负责人
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//服务用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetLogByOrderIDs 通过订单ID列查询服务单
func GetLogByOrderIDs(args *ArgsGetLogByOrderIDs) (dataList []FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, sn, sn_day, need_at, finish_at, user_id, org_id, bind_id, other_binds, mall_product_id, order_id, currency, price, pay_id, pay_at, des, address, config_id, params FROM service_housekeeping_log WHERE ($1 < 1 OR user_id = $1) AND delete_at < to_timestamp(1000000) AND order_id = ANY($2) AND ($3 < 1 OR org_id = $3) AND ($4 < 1 OR bind_id = $4)", args.UserID, args.OrderIDs, args.OrgID, args.BindID)
	return
}

// 获取订单关联的服务单
func getLogByOrderID(orderID int64) (dataList []FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, sn, sn_day, need_at, finish_at, user_id, org_id, bind_id, other_binds, mall_product_id, order_id, currency, price, pay_id, pay_at, des, address, config_id, params FROM service_housekeeping_log WHERE order_id = $1 AND delete_at < to_timestamp(1000000)", orderID)
	return
}

// 获取支付关联的服务单
func getLogByPayID(payID int64) (dataList []FieldsLog, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, sn, sn_day, need_at, finish_at, user_id, org_id, bind_id, other_binds, mall_product_id, order_id, currency, price, pay_id, pay_at, des, address, config_id, params FROM service_housekeeping_log WHERE pay_id = $1 AND delete_at < to_timestamp(1000000)", payID)
	return
}

// 获取订单关联的服务单数量
func getCountByOrderID(orderID int64) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_housekeeping_log WHERE order_id = $1 AND delete_at < to_timestamp(1000000)", orderID)
	return
}

// 批量获取服务单
func getListByIDs(ids pq.Int64Array) (dataList []FieldsLog, err error) {
	if err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, org_id, bind_id FROM service_housekeeping_log WHERE id = ANY($1) AND delete_at < to_timestamp(1000000)", ids); err != nil {
		return
	}
	return
}
