package ServiceUserInfoCost

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/golang-module/carbon"
	"time"
)

// ArgsGetCostList 获取耗能记录参数
type ArgsGetCostList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark" check:"mark" empty:"true"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark" check:"mark" empty:"true"`
}

// GetCostList 获取耗能记录
func GetCostList(args *ArgsGetCostList) (dataList []FieldsCost, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.RoomID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "room_id = :room_id"
		maps["room_id"] = args.RoomID
	}
	if args.InfoID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "info_id = :info_id"
		maps["info_id"] = args.InfoID
	}
	if args.RoomBindMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "room_bind_mark = :room_bind_mark"
		maps["room_bind_mark"] = args.RoomBindMark
	}
	if args.SensorMark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sensor_mark = :sensor_mark"
		maps["sensor_mark"] = args.SensorMark
	}
	tableName := "service_user_info_cost"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, room_id, info_id, config_id, room_bind_mark, sensor_mark, unit, currency, price FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsSetCost 记录耗能参数
type ArgsSetCost struct {
	//创建时间
	// 每次计算上一个小时形成的数据
	CreateAt string `db:"create_at" json:"createAt"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id"`
	//老人ID
	// 可能不存在，则忽略
	InfoID int64 `db:"info_id" json:"infoID" check:"id"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark" check:"mark"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark" check:"mark"`
	//阶段累计总量
	Unit float64 `db:"unit" json:"unit"`
	//阶段累计金额
	Currency int   `db:"currency" json:"currency" check:"currency"`
	Price    int64 `db:"price" json:"price" check:"price"`
}

// SetCost 记录耗能
// 该设计会强制修改数据，如果是新创建的数据，则不会
func SetCost(args *ArgsSetCost) (err error) {
	//获取最近1小时的数据
	var createAt time.Time
	createAt, err = CoreFilter.GetTimeByISO(args.CreateAt)
	if err != nil {
		return
	}
	startAt := carbon.CreateFromGoTime(createAt).StartOfDay()
	endAt := carbon.CreateFromGoTime(createAt).EndOfDay()
	var data FieldsCost
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_user_info_cost WHERE create_at >= $1 AND create_at <= $2 AND org_id = $3 AND room_id = $4 AND info_id = $5 AND room_bind_mark = $6 AND sensor_mark = $7", startAt, endAt, args.OrgID, args.RoomID, args.InfoID, args.RoomBindMark, args.SensorMark)
	//创建或直接修改数据
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info_cost SET unit = :unit, currency = :currency, price = :price WHERE id = :id", map[string]interface{}{
			"id":       data.ID,
			"unit":     args.Unit,
			"currency": args.Currency,
			"price":    args.Price,
		})
	} else {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_user_info_cost (create_at, org_id, room_id, info_id, config_id, room_bind_mark, sensor_mark, unit, currency, price) VALUES (:create_at,:org_id,:room_id,:info_id,0,:room_bind_mark,:sensor_mark,:unit,:currency,:price)", map[string]interface{}{
			"create_at":      args.CreateAt,
			"org_id":         args.OrgID,
			"room_id":        args.RoomID,
			"info_id":        args.InfoID,
			"room_bind_mark": args.RoomBindMark,
			"sensor_mark":    args.SensorMark,
			"unit":           args.Unit,
			"currency":       args.Currency,
			"price":          args.Price,
		})
	}
	return
}

// ArgsGetCostLast 获取最新的耗能数据参数
type ArgsGetCostLast struct {
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//信息ID
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark" check:"mark" empty:"true"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark" check:"mark" empty:"true"`
}

// GetCostLast 获取最新的耗能数据
func GetCostLast(args *ArgsGetCostLast) (data FieldsCost, err error) {
	if args.RoomID < 1 && args.InfoID < 1 {
		err = errors.New("room id and info id is empty")
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, unit, currency, price FROM service_user_info_cost WHERE ($1 < 1 OR org_id = $1) AND (($2 > 0 AND room_id = $2) OR ($3 > 0 AND info_id = $3)) AND ($4 = '' OR room_bind_mark = $4) AND ($5 = '' OR sensor_mark = $5) ORDER BY id DESC LIMIT 1", args.OrgID, args.RoomID, args.InfoID, args.RoomBindMark, args.SensorMark)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = errors.New("data is empty")
		return
	}
	return
}
