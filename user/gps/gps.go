package UserGPS

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"math"
	"time"
)

// ArgsGetList 获取定位追踪列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//所属用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country" empty:"true"`
	//所属城市
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
}

// GetList 获取定位追踪列表
func GetList(args *ArgsGetList) (dataList []FieldsGPS, dataCount int64, err error) {
	if err = checkMapType(args.MapType); err != nil {
		return
	}
	var where string
	maps := map[string]interface{}{}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.Country > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "country = :country"
		maps["country"] = args.Country
	}
	if args.City > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "city = :city"
		maps["city"] = args.City
	}
	if args.MapType > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "map_type = :map_type"
		maps["map_type"] = args.MapType
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_gps",
		"id",
		"SELECT id, create_at, user_id, country, city, map_type, longitude, latitude FROM user_gps WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetLast 获取最新的定位参数
type ArgsGetLast struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetLast 获取最新的定位
func GetLast(args *ArgsGetLast) (data FieldsGPS, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, user_id, country, city, map_type, longitude, latitude FROM user_gps WHERE user_id = $1 ORDER BY id DESC LIMIT 1", args.UserID)
	return
}

// ArgsGetLastByTime 获取指定时间之前最新的数据参数
type ArgsGetLastByTime struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//时间
	TimeAt time.Time `db:"time_at" json:"timeAt" check:"isoTime"`
}

// GetLastByTime 获取指定时间之前最新的数据
func GetLastByTime(args *ArgsGetLastByTime) (data FieldsGPS, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, user_id, country, city, map_type, longitude, latitude FROM user_gps WHERE user_id = $1 AND create_at <= $2 ORDER BY id DESC LIMIT 1", args.UserID, args.TimeAt)
	return
}

// ArgsGetMore 获取多人的最新订单参数
type ArgsGetMore struct {
	//用户IDs
	UserIDs pq.Int64Array `db:"user_ids" json:"userIDs" check:"ids"`
}

// GetMore 获取多人的最新订单
func GetMore(args *ArgsGetMore) (dataList []FieldsGPS, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, user_id, country, city, map_type, longitude, latitude FROM user_gps WHERE user_id = ANY($1) GROUP BY user_id, id ORDER BY id DESC", args.UserIDs)
	return
}

// ArgsCreate 添加新的定位参数
type ArgsCreate struct {
	//所属用户
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//所属城市
	City int `db:"city" json:"city" check:"city"`
	//地图制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// Create 添加新的定位
func Create(args *ArgsCreate) (err error) {
	//检查地图类型
	if err = checkMapType(args.MapType); err != nil {
		return
	}
	//定位如果为空，禁止提交
	if args.Longitude < 1 || args.Latitude < 1 {
		err = errors.New("lo or la is empty cannot insert data")
		return
	}
	//检查该记录是否可以写入？
	var data FieldsGPS
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT create_at, map_type, longitude, latitude FROM user_gps WHERE user_id = $1 ORDER BY id DESC LIMIT 1", args.UserID)
	//如果存在数据
	if err == nil {
		//当前时间
		nowTime := CoreFilter.GetNowTime()
		//极限记录时间，如果符合则跳过位移判断
		var recordLimitMaxTime int64
		recordLimitMaxTime, err = getUserGPSRecordLimitMaxTime()
		if nowTime.Unix()-data.CreateAt.Unix() < recordLimitMaxTime {
			//检查位移是否符合记录条件
			var recordLimitDistance float64
			recordLimitDistance, err = getUserGPSRecordLimitDistance()
			if err != nil {
				err = errors.New("get config recordLimitDistance, " + err.Error())
				return
			}
			if math.Abs(args.Longitude-data.Longitude) < recordLimitDistance || math.Abs(args.Latitude-data.Latitude) < recordLimitDistance {
				return
			}
		}
		//检查时间
		var recordLimitTime int64
		recordLimitTime, err = getUserGPSRecordLimitTime()
		if err != nil {
			err = errors.New("get config recordLimitTime, " + err.Error())
			return
		}
		if nowTime.Unix()-data.CreateAt.Unix() < recordLimitTime {
			return
		}
	}
	//创建新的记录
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_gps (user_id, country, city, map_type, longitude, latitude) VALUES (:user_id, :country, :city, :map_type,:longitude,:latitude)", args)
	return
}

// ArgsDeleteByUser 删除指定用户的所有定位参数
type ArgsDeleteByUser struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// DeleteByUser 删除指定用户的所有定位
func DeleteByUser(args *ArgsDeleteByUser) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_gps", "user_id = :user_id", args)
	return
}
