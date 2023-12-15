package TMSTransport

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAnalysis "github.com/fotomxq/weeekj_core/v5/core/sql/analysis"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserGPS "github.com/fotomxq/weeekj_core/v5/user/gps"
	"time"
)

// ArgsGetTransportGPSList 获取定位列表参数
type ArgsGetTransportGPSList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id" empty:"true"`
}

// GetTransportGPSList 获取定位列表
func GetTransportGPSList(args *ArgsGetTransportGPSList) (dataList []FieldsTransportGPS, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.TransportID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "transport_id = :transport_id"
		maps["transport_id"] = args.TransportID
	}
	if where == "" {
		where = "true"
	}
	tableName := "tms_transport_gps"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, transport_id, map_type, longitude, latitude FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if args.TransportID > -1 {
		if err == nil && len(dataList) > 0 {
			if args.Pages.Desc {
				if time.Now().Unix()-dataList[0].CreateAt.Unix() > 900 {
					updateTransportGPSByUserGPS(args.TransportID)
				}
			}
		} else {
			if dataCount < 1 {
				updateTransportGPSByUserGPS(args.TransportID)
			}
		}
	}
	return
}

// ArgsGetTransportGPSLast 获取指定配送单的最近定位数据参数
type ArgsGetTransportGPSLast struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id"`
}

// GetTransportGPSLast 获取指定配送单的最近定位数据
func GetTransportGPSLast(args *ArgsGetTransportGPSLast) (data FieldsTransportGPS, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, transport_id, map_type, longitude, latitude FROM tms_transport_gps WHERE (org_id = $1 OR $1 < 1) AND transport_id = $2 ORDER BY id DESC LIMIT 1", args.OrgID, args.TransportID)
	if err != nil || data.ID < 1 {
		updateTransportGPSByUserGPS(args.TransportID)
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, transport_id, map_type, longitude, latitude FROM tms_transport_gps WHERE (org_id = $1 OR $1 < 1) AND transport_id = $2 ORDER BY id DESC LIMIT 1", args.OrgID, args.TransportID)
	}
	return
}

// ArgsGetTransportGPSGroup 分组反馈数据参数
// 根据时间长度，超出3天的按照天划分，否则按照小时划分；超出3个月的按照月划分；超出24个月的按照年划分
type ArgsGetTransportGPSGroup struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id"`
	//最早的时间
	MinTime string `db:"min_time" json:"minTime" check:"isoTime"`
}

type DataGetTransportGPSGroup struct {
	//创建时间
	CreateAt string `db:"d" json:"createAt"`
	//地图制式
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// GetTransportGPSGroup 分组反馈数据
func GetTransportGPSGroup(args *ArgsGetTransportGPSGroup) (dataList []DataGetTransportGPSGroup, err error) {
	var minTime time.Time
	minTime, err = CoreFilter.GetTimeByISO(args.MinTime)
	if err != nil {
		return
	}
	where := "(org_id = :org_id OR :org_id < 1) AND (transport_id = :transport_id OR :transport_id < 1) AND create_at >= :min_time"
	maps := map[string]interface{}{
		"org_id":       args.OrgID,
		"transport_id": args.TransportID,
		"min_time":     minTime,
	}
	var timeType string
	if minTime.Unix() <= CoreFilter.GetNowTimeCarbon().SubYears(2).Time.Unix() {
		timeType = "year"
	} else {
		if minTime.Unix() <= CoreFilter.GetNowTimeCarbon().SubMonths(3).Time.Unix() {
			timeType = "month"
		} else {
			if minTime.Unix() <= CoreFilter.GetNowTimeCarbon().SubDays(3).Time.Unix() {
				timeType = "day"
			} else {
				timeType = "hour"
			}
		}
	}
	tableName := "tms_transport_gps"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", timeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", MAX(map_type) as map_type, AVG(longitude) as longitude, AVG(latitude) as latitude FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	if err != nil || len(dataList) < 1 {
		updateTransportGPSByUserGPS(args.TransportID)
		err = CoreSQL.GetList(
			Router2SystemConfig.MainDB.DB,
			&dataList,
			"SELECT "+timeField+", MAX(map_type) as map_type, AVG(longitude) as longitude, AVG(latitude) as latitude FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
			maps,
		)
	}
	return
}

// argsAppendTransportGPS 创建新的定位数据参数
type argsAppendTransportGPS struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id"`
	//地图制式
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// appendTransportGPS 创建新的定位数据参数
// 方法跟随用户定位进行联动
func appendTransportGPS(args *argsAppendTransportGPS) (err error) {
	//最近10分钟存在数据，跳出
	var dataGPS FieldsTransportGPS
	dataGPS, err = GetTransportGPSLast(&ArgsGetTransportGPSLast{
		OrgID:       -1,
		TransportID: args.TransportID,
	})
	if err == nil && dataGPS.CreateAt.Unix() >= CoreFilter.GetNowTimeCarbon().SubMinutes(10).Time.Unix() {
		return
	}
	//添加数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_gps (org_id, transport_id, map_type, longitude, latitude) VALUES (:org_id,:transport_id,:map_type,:longitude,:latitude)", args)
	return
}

// 从用户GPS信号给予配送单信号
func updateTransportGPSByUserGPS(transportID int64) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("update transport gps by user gps, ", r)
		}
	}()
	//最近10分钟存在数据，跳出
	var dataGPS FieldsTransportGPS
	err := Router2SystemConfig.MainDB.Get(&dataGPS, "SELECT id, create_at, org_id, transport_id, map_type, longitude, latitude FROM tms_transport_gps WHERE transport_id = $1 ORDER BY id DESC LIMIT 1", transportID)
	if err == nil && dataGPS.CreateAt.Unix() >= CoreFilter.GetNowTimeCarbon().SubMinutes(10).Time.Unix() {
		return
	}
	//获取配送单
	var data FieldsTransport
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id, bind_id FROM tms_transport WHERE delete_at < to_timestamp(1000000) AND id = $1 AND status != 3 AND bind_id > 0", transportID)
	if err != nil || data.ID < 1 {
		return
	}
	//查询配送员定位数据
	bindGPS, err := GetBindGPSLast(&ArgsGetBindGPSLast{
		OrgID:  data.OrgID,
		BindID: data.BindID,
	})
	if err == nil && bindGPS.ID > 0 {
		//写入数据
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_gps (org_id, transport_id, map_type, longitude, latitude) VALUES (:org_id,:transport_id,:map_type,:longitude,:latitude)", map[string]interface{}{
			"org_id":       bindGPS.OrgID,
			"transport_id": transportID,
			"map_type":     bindGPS.MapType,
			"longitude":    bindGPS.Longitude,
			"latitude":     bindGPS.Latitude,
		})
		return
	}
	//没有找到，则通过用户GPS定位查询
	//找到配送员数据
	bindData, err := OrgCore.GetBind(&OrgCore.ArgsGetBind{
		ID:     data.BindID,
		OrgID:  data.OrgID,
		UserID: -1,
	})
	if err != nil {
		return
	}
	//找到用户定位
	userGPS, err := UserGPS.GetLast(&UserGPS.ArgsGetLast{
		UserID: bindData.UserID,
	})
	if err != nil {
		return
	}
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_gps (org_id, transport_id, map_type, longitude, latitude) VALUES (:org_id,:transport_id,:map_type,:longitude,:latitude)", map[string]interface{}{
		"org_id":       bindData.OrgID,
		"transport_id": transportID,
		"map_type":     userGPS.MapType,
		"longitude":    userGPS.Longitude,
		"latitude":     userGPS.Latitude,
	})
	return
}
