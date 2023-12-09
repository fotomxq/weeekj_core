package OrgMission

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetAutoList 获取自动化列表
type ArgsGetAutoList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//创建人
	CreateBindID int64 `json:"createBindID" check:"id" empty:"true"`
	//执行人
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//其他执行人
	OtherBindID int64 `json:"otherBindID" check:"id" empty:"true"`
	//级别
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags int64 `db:"tags" json:"tags" check:"id" empty:"true"`
	//开始时间范围
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetAutoList 获取任务列表
func GetAutoList(args *ArgsGetAutoList) (dataList []FieldsAuto, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.CreateBindID > -1 {
		where = where + " AND create_bind_id = :create_bind_id"
		maps["create_bind_id"] = args.CreateBindID
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.OtherBindID > -1 {
		where = where + " AND :other_bind_id = ANY(other_bind_ids)"
		maps["other_bind_id"] = args.OtherBindID
	}
	if args.Level > -1 {
		where = where + " AND level = :level"
		maps["level"] = args.Level
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.Tags > -1 {
		where = where + " AND :tags = ANY(tags)"
		maps["tags"] = args.Tags
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_mission_auto"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, time_type, time_n, skip_holiday, start_hour, start_minute, end_hour, end_minute, org_id, create_bind_id, bind_id, other_bind_ids, title, des, des_files, start_at, next_at, end_at, tip_id, level, sort_id, tags, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "start_at", "end_at", "level"},
	)
	return
}

// ArgsCreateAuto 创建新的自动化参数
type ArgsCreateAuto struct {
	//时间类型
	// 0 每天重复 day / 1 每周重复 week / 2 每月重复 month / 3 临时1次 once
	// 4 每隔N天重复 day_n / 5 每隔N周重复 week_n / 6 每隔N月重复 month_n
	TimeType int `db:"time_type" json:"timeType"`
	//扩展N
	// 重复时间内，数组的第一个值作为相隔N；
	// 重复周内，数组代表指定的星期1-7
	TimeN pq.Int64Array `db:"time_n" json:"timeN"`
	//是否跳过节假日
	SkipHoliday bool `db:"skip_holiday" json:"skipHoliday" check:"bool"`
	//开始时间
	StartHour   int `db:"start_hour" json:"startHour"`
	StartMinute int `db:"start_minute" json:"startMinute"`
	//结束时间
	EndHour   int `db:"end_hour" json:"endHour"`
	EndMinute int `db:"end_minute" json:"endMinute"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//创建人
	CreateBindID int64 `db:"create_bind_id" json:"createBindID" check:"id"`
	//执行人
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//其他执行人
	OtherBindIDs pq.Int64Array `db:"other_bind_ids" json:"otherBindIDs" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//执行时间
	StartAt time.Time `db:"start_at" json:"startAt" check:"isoTime" empty:"true"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt" check:"isoTime" empty:"true"`
	//是否需提醒
	// -1 不需要; 0 需要等待提醒中; >0 已经触发提醒的ID
	TipID int64 `db:"tip_id" json:"tipID"`
	//级别
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateAuto 创建新的自动化
func CreateAuto(args *ArgsCreateAuto) (data FieldsAuto, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_mission_auto", "INSERT INTO org_mission_auto (time_type, time_n, skip_holiday, start_hour, start_minute, end_hour, end_minute, org_id, create_bind_id, bind_id, other_bind_ids, title, des, des_files, start_at, end_at, next_at, tip_id, level, sort_id, tags, params) VALUES (:time_type, :time_n, :skip_holiday, :start_hour, :start_minute, :end_hour, :end_minute,:org_id,:create_bind_id,:bind_id,:other_bind_ids,:title,:des,:des_files,:start_at,:end_at,:start_at,:tip_id,:level,:sort_id,:tags,:params)", args, &data)
	return
}

// ArgsUpdateAuto 修改自动化参数
type ArgsUpdateAuto struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//任意一种形式包含此人
	OperateBindID int64 `db:"operate_bind_id" json:"operateBindID" check:"id" empty:"true"`
	//时间类型
	// 0 每天重复 day / 1 每周重复 week / 2 每月重复 month / 3 临时1次 once
	// 4 每隔N天重复 day_n / 5 每隔N周重复 week_n / 6 每隔N月重复 month_n
	TimeType int `db:"time_type" json:"timeType"`
	//扩展N
	// 重复时间内，数组的第一个值作为相隔N；
	// 重复周内，数组代表指定的星期1-7
	TimeN pq.Int64Array `db:"time_n" json:"timeN"`
	//是否跳过节假日
	SkipHoliday bool `db:"skip_holiday" json:"skipHoliday" check:"bool"`
	//开始时间
	StartHour   int `db:"start_hour" json:"startHour"`
	StartMinute int `db:"start_minute" json:"startMinute"`
	//结束时间
	EndHour   int `db:"end_hour" json:"endHour"`
	EndMinute int `db:"end_minute" json:"endMinute"`
	//执行人
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//其他执行人
	OtherBindIDs pq.Int64Array `db:"other_bind_ids" json:"otherBindIDs" check:"ids" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//执行时间
	StartAt time.Time `db:"start_at" json:"startAt" check:"isoTime" empty:"true"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt" check:"isoTime" empty:"true"`
	//是否需提醒
	// -1 不需要; 0 需要等待提醒中; >0 已经触发提醒的ID
	TipID int64 `db:"tip_id" json:"tipID"`
	//级别
	Level int `db:"level" json:"level" check:"intThan0" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateAuto 修改自动化
func UpdateAuto(args *ArgsUpdateAuto) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_mission_auto SET update_at = NOW(), time_type = :time_type, time_n = :time_n, skip_holiday = :skip_holiday, start_hour = :start_hour, start_minute = :start_minute, end_hour = :end_hour, end_minute = :end_minute, bind_id = :bind_id, other_bind_ids = :other_bind_ids, title = :title, des = :des, des_files = :des_files, start_at = :start_at, next_at = :start_at, end_at = :end_at, tip_id = :tip_id, level = :level, sort_id = :sort_id, tags = :tags, params = :params WHERE id = :id AND org_id = :org_id AND (:operate_bind_id < 1 OR create_bind_id = :operate_bind_id)", args)
	return
}

// ArgsDeleteAuto 删除自动化参数
type ArgsDeleteAuto struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//任意一种形式包含此人
	OperateBindID int64 `db:"operate_bind_id" json:"operateBindID" check:"id" empty:"true"`
}

// DeleteAuto 删除自动化
func DeleteAuto(args *ArgsDeleteAuto) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_mission_auto", "id = :id AND org_id = :org_id AND (:operate_bind_id < 1 OR create_bind_id = :operate_bind_id)", args)
	return
}
