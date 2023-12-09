package ServiceHealthSelf

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 查看检疫日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 必填
	OrgID int64 `json:"orgID" check:"id"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//健康码状态
	// 0 正常（绿）; 1 警告（黄）; 2 危险（红）
	HealthStatus int `db:"health_status" json:"healthStatus" check:"intThan0" empty:"true"`
	//行程卡状态
	// 0 正常（绿）; 1 警告（黄）; 2 危险（红）
	TravelStatus int `db:"travel_status" json:"travelStatus" check:"intThan0" empty:"true"`
	//核酸结果
	// 0 正常（阴性）; 1 异常（阳性）
	NAReportStatus int `db:"na_report_status" json:"naReportStatus" check:"intThan0" empty:"true"`
	//总的检查结果
	// 0 正常（绿）; 1 警告（黄）; 2 危险（红）
	Result int `db:"result" json:"result" check:"intThan0" empty:"true"`
}

// GetLogList 查看检疫日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.OrgBindID > -1 {
		where = where + " AND org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.HealthStatus > -1 {
		where = where + " AND health_status = :health_status"
		maps["health_status"] = args.HealthStatus
	}
	if args.TravelStatus > -1 {
		where = where + " AND travel_status = :travel_status"
		maps["travel_status"] = args.TravelStatus
	}
	if args.NAReportStatus > -1 {
		where = where + " AND na_report_status = :na_report_status"
		maps["na_report_status"] = args.NAReportStatus
	}
	if args.Result > -1 {
		where = where + " AND result = :result"
		maps["result"] = args.Result
	}
	var rawList []FieldsLog
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"service_health_self_log",
		"id",
		"SELECT org_id, org_bind_id, user_id FROM service_health_self_log WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getLog(v.OrgID, v.OrgBindID, v.UserID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// CheckLog 检查某个来源是否核对通过
func CheckLog(orgID int64, orgBindID int64, userID int64) (data FieldsLog, b bool) {
	data = getLog(orgID, orgBindID, userID)
	if data.ID < 1 {
		return
	}
	b = data.HealthStatus == 0 && data.TravelStatus == 0
	return
}

// ArgsAppendLog 添加新的记录参数
type ArgsAppendLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//健康码状态
	// 0 正常（绿）; 1 警告（黄）; 2 危险（红）
	HealthStatus int `db:"health_status" json:"healthStatus" check:"intThan0" empty:"true"`
	//健康码附加文件
	HealthFileID int64 `db:"health_file_id" json:"healthFileID" check:"id" empty:"true"`
	//行程卡状态
	// 0 正常（绿）; 1 警告（黄）; 2 危险（红）
	TravelStatus int `db:"travel_status" json:"travelStatus" check:"intThan0" empty:"true"`
	//行程卡附加文件
	TravelFileID int64 `db:"travel_file_id" json:"travelFileID" check:"id" empty:"true"`
	//体温
	// 小数点保留2位数x100
	BodyTemperature int `db:"body_temperature" json:"bodyTemperature" check:"intThan0" empty:"true"`
	//核酸报告截图
	NAReportFileID int64 `db:"na_report_file_id" json:"naReportFileID" check:"id" empty:"true"`
	//核酸结果
	// 0 正常（阴性）; 1 异常（阳性）
	NAReportStatus int `db:"na_report_status" json:"naReportStatus" check:"intThan0" empty:"true"`
}

// AppendLog 添加新的记录
func AppendLog(args *ArgsAppendLog) (err error) {
	switch args.HealthStatus {
	case 0:
	case 1:
	case 2:
	default:
		err = errors.New("health status error")
		return
	}
	switch args.TravelStatus {
	case 0:
	case 1:
	case 2:
	default:
		err = errors.New("travel status error")
		return
	}
	switch args.NAReportStatus {
	case 0:
	case 1:
	default:
		err = errors.New("na report status error")
		return
	}
	var result int
	if args.HealthStatus == 0 && args.TravelStatus == 0 && args.NAReportStatus == 0 {
		result = 0
	} else {
		if args.HealthStatus == 2 || args.TravelStatus == 2 || args.NAReportStatus == 1 {
			result = 2
		} else {
			if args.HealthStatus == 1 || args.TravelStatus == 1 {
				result = 1
			}
		}
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_health_self_log (org_id, org_bind_id, user_id, health_status, health_file_id, travel_status, travel_file_id, body_temperature, na_report_file_id, na_report_status, result) VALUES (:org_id,:org_bind_id,:user_id,:health_status,:health_file_id,:travel_status,:travel_file_id,:body_temperature,:na_report_file_id,:na_report_status,:result)", map[string]interface{}{
		"org_id":            args.OrgID,
		"org_bind_id":       args.OrgBindID,
		"user_id":           args.UserID,
		"health_status":     args.HealthStatus,
		"health_file_id":    args.HealthFileID,
		"travel_status":     args.TravelStatus,
		"travel_file_id":    args.TravelFileID,
		"body_temperature":  args.BodyTemperature,
		"na_report_file_id": args.NAReportFileID,
		"na_report_status":  args.NAReportStatus,
		"result":            result,
	})
	if err != nil {
		return
	}
	return
}

func getLog(orgID int64, orgBindID int64, userID int64) (data FieldsLog) {
	cacheMark := getLogCacheMark(orgID, orgBindID, userID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, user_id, health_status, health_file_id, travel_status, travel_file_id, body_temperature, na_report_file_id, na_report_status, result FROM service_health_self_log WHERE org_id = $1 AND org_bind_id = $2 AND user_id = $3", orgID, orgBindID, userID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 10800)
	return
}

// 缓冲
func getLogCacheMark(orgID int64, orgBindID int64, userID int64) string {
	return fmt.Sprint("service:health:self:log:org:", orgID, ".", orgBindID, ".", userID)
}

func deleteLogCache(orgID int64, orgBindID int64, userID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getLogCacheMark(orgID, orgBindID, userID))
}
