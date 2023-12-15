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

// ArgsGetBindGPSList 获取定位列表参数
type ArgsGetBindGPSList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// GetBindGPSList 获取定位列表
func GetBindGPSList(args *ArgsGetBindGPSList) (dataList []FieldsBindGPS, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if where == "" {
		where = "true"
	}
	tableName := "tms_transport_bind_gps"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, bind_id, map_type, longitude, latitude FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if args.BindID > -1 {
		if err == nil && len(dataList) > 0 {
			if args.Pages.Desc {
				if time.Now().Unix()-dataList[0].CreateAt.Unix() > 900 {
					updateBindGPSByUserGPS(args.BindID)
				}
			}
		} else {
			if dataCount < 1 {
				updateBindGPSByUserGPS(args.BindID)
			}
		}
	}
	return
}

// ArgsGetBindGPSLast 获取指定人员的最近定位数据参数
type ArgsGetBindGPSLast struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
}

// GetBindGPSLast 获取指定人员的最近定位数据
func GetBindGPSLast(args *ArgsGetBindGPSLast) (data FieldsBindGPS, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, bind_id, map_type, longitude, latitude FROM tms_transport_bind_gps WHERE ($1 < 1 OR org_id = $1) AND bind_id = $2 ORDER BY id DESC LIMIT 1", args.OrgID, args.BindID)
	if err != nil || data.ID < 1 {
		updateBindGPSByUserGPS(args.BindID)
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, bind_id, map_type, longitude, latitude FROM tms_transport_bind_gps WHERE ($1 < 1 OR org_id = $1) AND bind_id = $2 ORDER BY id DESC LIMIT 1", args.OrgID, args.BindID)
	}
	return
}

// ArgsGetBindGPSGroup 分组反馈数据参数
// 根据时间长度，超出3天的按照天划分，否则按照小时划分；超出3个月的按照月划分；超出24个月的按照年划分
type ArgsGetBindGPSGroup struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//最早的时间
	MinTime string `db:"min_time" json:"minTime" check:"isoTime"`
}

type DataGetBindGPSGroup struct {
	//创建时间
	CreateAt string `db:"d" json:"createAt"`
	//地图制式
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// GetBindGPSGroup 分组反馈数据
func GetBindGPSGroup(args *ArgsGetBindGPSGroup) (dataList []DataGetBindGPSGroup, err error) {
	var minTime time.Time
	minTime, err = CoreFilter.GetTimeByISO(args.MinTime)
	if err != nil {
		return
	}
	where := "(org_id = :org_id OR :org_id < 1) AND (bind_id = :bind_id OR :bind_id < 1) AND create_at >= :min_time"
	maps := map[string]interface{}{
		"org_id":   args.OrgID,
		"bind_id":  args.BindID,
		"min_time": minTime,
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
	tableName := "tms_transport_bind_gps"
	timeField := CoreSQLAnalysis.GetAnalysisQueryField("create_at", timeType, "d")
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT "+timeField+", MAX(map_type) as map_type, AVG(longitude) as longitude, AVG(latitude) as latitude FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
		maps,
	)
	if err != nil || len(dataList) < 1 {
		updateBindGPSByUserGPS(args.BindID)
		err = CoreSQL.GetList(
			Router2SystemConfig.MainDB.DB,
			&dataList,
			"SELECT "+timeField+", MAX(map_type) as map_type, AVG(longitude) as longitude, AVG(latitude) as latitude FROM "+tableName+" WHERE "+where+" GROUP BY d ORDER BY d",
			maps,
		)
	}
	return
}

// argsAppendBindGPS 创建新的定位数据参数
type argsAppendBindGPS struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID"`
	//地图制式
	// WGS-84 / GCJ-02 / BD-09
	MapType int `db:"map_type" json:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// appendBindGPS 创建新的定位数据
// 方法跟随用户定位进行联动
func appendBindGPS(args *argsAppendBindGPS) (err error) {
	//最近10分钟存在数据，跳出
	var dataGPS FieldsBindGPS
	dataGPS, err = GetBindGPSLast(&ArgsGetBindGPSLast{
		OrgID:  -1,
		BindID: args.BindID,
	})
	if err == nil && dataGPS.CreateAt.Unix() >= CoreFilter.GetNowTimeCarbon().SubMinutes(10).Time.Unix() {
		return
	}
	//添加数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_bind_gps (org_id, bind_id, map_type, longitude, latitude) VALUES (:org_id,:bind_id,:map_type,:longitude,:latitude)", args)
	return
}

// 从用户GPS信号给予配送人员信号
func updateBindGPSByUserGPS(bindID int64) {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("update transport bind gps by user gps, ", r)
		}
	}()
	//最近10分钟存在数据，跳出
	var dataGPS FieldsBindGPS
	err := Router2SystemConfig.MainDB.Get(&dataGPS, "SELECT id, create_at, org_id, bind_id, map_type, longitude, latitude FROM tms_transport_bind_gps WHERE bind_id = $1 ORDER BY id DESC LIMIT 1", bindID)
	if err == nil && dataGPS.CreateAt.Unix() >= CoreFilter.GetNowTimeCarbon().SubMinutes(10).Time.Unix() {
		return
	}
	//获取配送员数据
	bindData, err := OrgCore.GetBind(&OrgCore.ArgsGetBind{
		ID:     dataGPS.BindID,
		OrgID:  dataGPS.OrgID,
		UserID: -1,
	})
	if err != nil {
		return
	}
	//获取用户数据
	userGPS, err := UserGPS.GetLast(&UserGPS.ArgsGetLast{
		UserID: bindData.UserID,
	})
	if err != nil {
		return
	}
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_bind_gps (org_id, bind_id, map_type, longitude, latitude) VALUES (:org_id,:bind_id,:map_type,:longitude,:latitude)", map[string]interface{}{
		"org_id":    bindData.OrgID,
		"bind_id":   bindData.ID,
		"map_type":  userGPS.MapType,
		"longitude": userGPS.Longitude,
		"latitude":  userGPS.Latitude,
	})
	return
}
