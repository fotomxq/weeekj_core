package OrgTime

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/lib/pq"
	"time"
)

//FieldsWorkTime 上下班处理，提供给用户
type FieldsWorkTime struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//过期时间
	// 过期后自动失效，用于排临时班
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//掌管该数据的组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分组ID列
	Groups pq.Int64Array `db:"groups" json:"groups"`
	//绑定人列
	Binds pq.Int64Array `db:"binds" json:"binds"`
	//名称
	Name string `db:"name" json:"name"`
	//当前上下班状态
	// 根据配置时间自动调整，外部读取即可使用
	IsWork bool `db:"is_work" json:"isWork"`
	//时间配置组
	Configs FieldsConfigs `db:"configs" json:"configs"`
	//轮动任务
	RotConfig FieldsConfigRot `db:"rot_config" json:"rotConfig"`
}

//FieldsConfigs 配置组
type FieldsConfigs struct {
	//每年月时间
	// 1-12月
	Month []int `db:"month" json:"month"`
	//每月时间
	// 检查1-31的天
	MonthDay []int `db:"month_day" json:"monthDay"`
	//每月第几个周
	// 1-6周
	MonthWeek []int `db:"month_week" json:"monthWeek"`
	//每周7天时间
	// 检查1\2\3\4\5\6\0
	Week []int `db:"week" json:"week"`
	//具体的分配时间组合
	WorkTime FieldsWorkTimeTimes `db:"work_time" json:"workTime"`
	//自动跳过节假日
	AllowHoliday bool `db:"allow_holiday" json:"allowHoliday"`
}

//Value sql底层处理器
func (t FieldsConfigs) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

//FieldsWorkTimeTimes 上下班时间 24小时制
type FieldsWorkTimeTimes []FieldsWorkTimeTime

//Value sql底层处理器
func (t FieldsWorkTimeTimes) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsWorkTimeTimes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsWorkTimeTime struct {
	//上班时间
	StartHour   int `db:"start_hour" json:"startHour"`
	StartMinute int `db:"start_minute" json:"startMinute"`
	//下班时间
	EndHour   int `db:"end_hour" json:"endHour"`
	EndMinute int `db:"end_minute" json:"endMinute"`
}

//Value sql底层处理器
func (t FieldsWorkTimeTime) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsWorkTimeTime) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

//FieldsConfigRot 轮动任务
type FieldsConfigRot struct {
	//当前轮动到的位置
	NowKey int `db:"now_key" json:"nowKey"`
	//切换间隔日
	// 如果给1，则每个轮动间隔1天；必须大于0，否则自动修正为1
	DiffDay int `db:"diff_day" json:"diffDay"`
	//具体的分配时间组合
	WorkTime FieldsWorkTimeTimes `db:"work_time" json:"workTime"`
}

//Value sql底层处理器
func (t FieldsConfigRot) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigRot) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
