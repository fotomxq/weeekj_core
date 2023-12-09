package OrgMission

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsAuto 任务自动生成模块
type FieldsAuto struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//时间类型
	// 0 每天重复 day / 1 每周重复 week / 2 每月重复 month / 3 临时1次 once
	// 4 每隔N天重复 day_n / 5 每隔N周重复 week_n / 6 每隔N月重复 month_n
	TimeType int `db:"time_type" json:"timeType"`
	//扩展N
	// 重复时间内，数组的第一个值作为相隔N；
	// 重复周内，数组代表指定的星期1-7
	TimeN pq.Int64Array `db:"time_n" json:"timeN"`
	//是否跳过节假日
	SkipHoliday bool `db:"skip_holiday" json:"skipHoliday"`
	//开始时间
	StartHour   int `db:"start_hour" json:"startHour"`
	StartMinute int `db:"start_minute" json:"startMinute"`
	//结束时间
	EndHour   int `db:"end_hour" json:"endHour"`
	EndMinute int `db:"end_minute" json:"endMinute"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//创建人
	CreateBindID int64 `db:"create_bind_id" json:"createBindID"`
	//执行人
	BindID int64 `db:"bind_id" json:"bindID"`
	//其他执行人
	OtherBindIDs pq.Int64Array `db:"other_bind_ids" json:"otherBindIDs"`
	//标题
	Title string `db:"title" json:"title"`
	//描述
	Des string `db:"des" json:"des"`
	//文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//执行时间
	StartAt time.Time `db:"start_at" json:"startAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//下一次启动时间
	NextAt time.Time `db:"next_at" json:"nextAt"`
	//是否需提醒
	// -1 不需要; 0 需要等待提醒中; >0 已经触发提醒的ID
	TipID int64 `db:"tip_id" json:"tipID"`
	//级别
	Level int `db:"level" json:"level"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//标签列
	Tags pq.Int64Array `db:"tags" json:"tags"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
