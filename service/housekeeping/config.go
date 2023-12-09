package ServiceHousekeeping

import (
	ClassConfig "gitee.com/weeekj/weeekj_core/v5/class/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetConfigList 获取配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(title ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_housekeeping_config"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, org_id, sort, title, start_at, end_at, limit_count, now_count FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "sort"},
	)
	if err != nil {
		return
	}
	return
}

// ArgsGetConfigData 获取重组的数据列表参数
type ArgsGetConfigData struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DataGetConfigData 获取重组的数据列表数据
type DataGetConfigData struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//描述标题
	Title string `db:"title" json:"title"`
	//预约时间
	StartAt string `db:"start_at" json:"startAt"`
	//是否可以预约
	IsCan bool `db:"is_can" json:"isCan"`
}

// GetConfigData 获取重组的数据列表
func GetConfigData(args *ArgsGetConfigData) (dataList []DataGetConfigData, err error) {
	housekeepingShowDay, _ := OrgCore.Config.GetConfigValInt(&ClassConfig.ArgsGetConfig{
		BindID:    args.OrgID,
		Mark:      "HousekeepingShowDay",
		VisitType: "admin",
	})
	if housekeepingShowDay < 1 {
		housekeepingShowDay = 7
	}
	if housekeepingShowDay > 30 {
		housekeepingShowDay = 30
	}
	nowAt := CoreFilter.GetNowTimeCarbon()
	var rawList []FieldsConfig
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id, title, start_at, end_at, limit_count FROM service_housekeeping_config WHERE org_id = $1 ORDER BY sort", args.OrgID)
	if err == nil && len(rawList) > 0 {
		for kDay := 0; kDay < housekeepingShowDay; kDay++ {
			for k := 0; k < len(rawList); k++ {
				v := rawList[k]
				//生成开始和结束时间
				configStartAt := CoreFilter.GetCarbonByTime(v.StartAt)
				nowStartAt := nowAt.AddDays(kDay).SetHour(configStartAt.Hour()).SetMinute(configStartAt.Minute()).SetSecond(configStartAt.Second())
				//检查开始时间是否是过去？跳过处理
				if nowStartAt.Time.Unix() < nowAt.Time.Unix() {
					continue
				}
				//生成结束时间
				configEndAt := CoreFilter.GetCarbonByTime(v.EndAt)
				nowEndAt := nowAt.AddDays(kDay).SetHour(configEndAt.Hour()).SetMinute(configEndAt.Minute()).SetSecond(configEndAt.Second())
				//检查isCan，根据预约时间判断，判断当前预约个数是否达到上限？如果达到则标记无法预约、反之允许
				var count int
				_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_housekeeping_log WHERE delete_at < to_timestamp(1000000) AND finish_at < to_timestamp(1000000) AND need_at >= $1 AND need_at <= $2", nowStartAt.Time, nowEndAt.Time)
				//写入数据集
				dataList = append(dataList, DataGetConfigData{
					ID:      v.ID,
					Title:   v.Title,
					StartAt: nowStartAt.Time.Format("2006-01-02 15:04:05"),
					IsCan:   count < v.LimitCount,
				})
			}
		}
	}
	return
}

// ArgsCreateConfig 创建配置参数
type ArgsCreateConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//排序
	Sort int `db:"sort" json:"sort" check:"intThan0" empty:"true"`
	//描述标题
	Title string `db:"title" json:"title" check:"title"`
	//服务时间范围
	StartAt string `db:"start_at" json:"startAt" check:"defaultTime"`
	EndAt   string `db:"end_at" json:"endAt" check:"defaultTime"`
	//服务限制
	LimitCount int `db:"limit_count" json:"limitCount" check:"intThan0"`
	//未来预约的数量
	NowCount int `db:"now_count" json:"nowCount" check:"intThan0"`
}

// CreateConfig 创建配置
func CreateConfig(args *ArgsCreateConfig) (err error) {
	var startAt, endAt time.Time
	startAt, err = CoreFilter.GetTimeByDefault(args.StartAt)
	if err != nil {
		return
	}
	endAt, err = CoreFilter.GetTimeByDefault(args.EndAt)
	if err != nil {
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_housekeeping_config (org_id, sort, title, start_at, end_at, limit_count, now_count) VALUES(:org_id, :sort, :title, :start_at, :end_at, :limit_count, :now_count)", map[string]interface{}{
		"org_id":      args.OrgID,
		"sort":        args.Sort,
		"title":       args.Title,
		"start_at":    startAt,
		"end_at":      endAt,
		"limit_count": args.LimitCount,
		"now_count":   args.NowCount,
	})
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//排序
	Sort int `db:"sort" json:"sort" check:"intThan0" empty:"true"`
	//描述标题
	Title string `db:"title" json:"title" check:"title"`
	//服务时间范围
	StartAt string `db:"start_at" json:"startAt" check:"defaultTime"`
	EndAt   string `db:"end_at" json:"endAt" check:"defaultTime"`
	//服务限制
	LimitCount int `db:"limit_count" json:"limitCount" check:"intThan0"`
	//未来预约的数量
	NowCount int `db:"now_count" json:"nowCount" check:"intThan0"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	var startAt, endAt time.Time
	startAt, err = CoreFilter.GetTimeByDefault(args.StartAt)
	if err != nil {
		return
	}
	endAt, err = CoreFilter.GetTimeByDefault(args.EndAt)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_config SET update_at = NOW(), sort = :sort, title = :title, start_at = start_at, end_at = :end_at, limit_count = :limit_count, now_count = :now_count WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":          args.ID,
		"org_id":      args.OrgID,
		"sort":        args.Sort,
		"title":       args.Title,
		"start_at":    startAt,
		"end_at":      endAt,
		"limit_count": args.LimitCount,
		"now_count":   args.NowCount,
	})
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteConfig 删除配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "service_housekeeping_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsCheckConfigTime 检查配置预约参数
type ArgsCheckConfigTime struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

type DataCheckConfigTime struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//服务时间范围
	StartAt time.Time `db:"start_at" json:"startAt"`
	EndAt   time.Time `db:"end_at" json:"endAt"`
	//服务限制
	LimitCount int `db:"limit_count" json:"limitCount"`
	//未来预约的数量
	NowCount int `db:"now_count" json:"nowCount"`
}

// CheckConfigTime 检查配置预约
func CheckConfigTime(args *ArgsCheckConfigTime) (data DataCheckConfigTime, b bool) {
	//获取可预约数据集
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, start_at, end_at, limit_count, now_count FROM service_housekeeping_config WHERE id = $1", args.ID)
	if err != nil {
		return
	}
	if data.LimitCount <= data.NowCount {
		return
	}
	//生成开始和结束时间
	nowAt := CoreFilter.GetNowTimeCarbon()
	configStartAt := CoreFilter.GetCarbonByTime(data.StartAt)
	nowStartAt := nowAt.SetHour(configStartAt.Hour()).SetMinute(configStartAt.Minute()).SetSecond(configStartAt.Second())
	data.StartAt = nowStartAt.Time
	configEndAt := CoreFilter.GetCarbonByTime(data.EndAt)
	nowEndAt := nowAt.SetHour(configEndAt.Hour()).SetMinute(configEndAt.Minute()).SetSecond(configEndAt.Second())
	data.EndAt = nowEndAt.Time
	//检查时间是否过期？
	if data.StartAt.Unix() < nowAt.Time.Unix() {
		return
	}
	//检查预约数量是否超出限制？
	var count int
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM service_housekeeping_log WHERE delete_at < to_timestamp(1000000) AND finish_at < to_timestamp(1000000) AND need_at >= $1 AND need_at <= $2", nowStartAt.Time, nowEndAt.Time)
	b = count < data.LimitCount
	//反馈数据
	return
}

// 增加次数
func addConfig(id int64) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_config SET now_count = now_count + 1 WHERE id = :id", map[string]interface{}{
		"id": id,
	})
	return
}

// reduceConfig 减少次数
// 该方法不会使用，因为nowCount改为累计总的记录
func reduceConfig(id int64) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_config SET now_count = now_count - 1 WHERE id = :id AND now_count > 0", map[string]interface{}{
		"id": id,
	})
	return
}
