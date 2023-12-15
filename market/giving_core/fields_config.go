package MarketGivingCore

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsConfig 赠送内容配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//推荐后奖励配置
	MarketConfigID int64 `db:"market_config_id" json:"marketConfigID"`
	//领取周期类型
	// 0 不限制; 1 一次性; 2 每天限制; 3 每周限制; 4 每月限制; 5 每季度限制; 6 每年限制
	LimitTimeType int `db:"limit_time_type" json:"limitTimeType"`
	//领取次数
	LimitCount int `db:"limit_count" json:"limitCount"`
	//领取积分
	UserIntegral int64 `db:"user_integral" json:"userIntegral"`
	//领取用户订阅
	UserSubs FieldsConfigUserSubs `db:"user_subs" json:"userSubs"`
	//领取票据
	UserTickets FieldsConfigUserTickets `db:"user_tickets" json:"userTickets"`
	//奖励金储蓄标识码
	DepositConfigMark string `db:"deposit_config_mark" json:"depositConfigMark"`
	//奖励金额
	Price int64 `db:"price" json:"price"`
	//奖励次数
	Count int64 `db:"count" json:"count"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type FieldsConfigUserSubs []FieldsConfigUserSub

type FieldsConfigUserSub struct {
	//订阅配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//赠送数量
	Count int64 `db:"count" json:"count"`
	//赠送时间长度
	CountTime int64 `db:"count_time" json:"countTime"`
}

func (t FieldsConfigUserSubs) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigUserSubs) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

func (t FieldsConfigUserSub) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigUserSub) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsConfigUserTickets []FieldsConfigUserTicket

type FieldsConfigUserTicket struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//票据数量
	Count int64 `db:"count" json:"count"`
}

func (t FieldsConfigUserTickets) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigUserTickets) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

func (t FieldsConfigUserTicket) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigUserTicket) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
