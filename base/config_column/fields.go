package BaseConfigColumn

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// 列头核心表
type FieldsColumn struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//来源系统
	// 0 系统层 / 1 组织层 / 2 用户层
	// 系统层影响所有系统配置设计，该设计全系统通用，但用户层可自定义覆盖设定
	// 组织层用于声明组织内部的所有列头，用于覆盖系统层的设计
	// 用户层可直接覆盖系统层或组织层的设定
	System int `db:"system" json:"system"`
	//来源ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//标识码
	// 在来源系统内，该数据必须唯一，前端可识别具体是哪个页面的哪个组件
	// 不同层级可声明一个标识码，系统将反馈最大层级的数据
	Mark string `db:"mark" json:"mark"`
	//保存数据集
	// 前后顺序将按照该顺序一致
	Data FieldsChildList `db:"data" json:"data"`
}

// 子表
type FieldsChildList []FieldsChild

// sql底层处理器
func (t FieldsChildList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsChildList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsChild struct {
	//头标识码
	Mark string `db:"mark" json:"mark"`
	//头名称
	Name string `db:"name" json:"name"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// sql底层处理器
func (t FieldsChild) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsChild) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
