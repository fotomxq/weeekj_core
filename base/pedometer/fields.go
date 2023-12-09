package PedometerCore

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsPedometerType 主要结构
type FieldsPedometerType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//过期时间
	// 超期后将自动删除数据
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//创建来源和创建来源ID
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//值
	Count int `db:"count" json:"count"`
}

// FieldsPedometerConfigType 默认配置表
type FieldsPedometerConfigType struct {
	//系统来源
	Mark string `db:"mark" json:"mark"`
	//值
	Count int `db:"count" json:"count"`
	//默认超期时间
	DefaultExpire string `db:"default_expire" json:"defaultExpire"`
	//最小值和最大值
	MinCount int `db:"min_count" json:"minCount"`
	MaxCount int `db:"max_count" json:"maxCount"`
	//递增还是递减
	IsAdd bool `db:"is_add" json:"isAdd"`
}
