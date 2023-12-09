package ToolsHolidaySeason

import "time"

type FieldsHolidaySeason struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//对应的日期
	DateAt time.Time `db:"date_at" json:"dateAt"`
	//节假日类型
	// 0 工作日、1 周末、2 节日、3 调休
	Status int `db:"status" json:"status"`
	//是否房价
	IsHoliday bool `db:"is_holiday" json:"isHoliday"`
	//名称
	// eg: 周二
	Name string `db:"name" json:"name"`
	//薪资倍数
	Wage int `db:"wage" json:"wage"`
	//是否强制修改
	// 如果启动，则只认准修改数据，将不同步API
	IsForce bool `db:"is_force" json:"isForce"`
}