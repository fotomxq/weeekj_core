package BasePageStyle

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsPage 页面
type FieldsPage struct {
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
	//系统
	System string `db:"system" json:"system"`
	//页面识别码
	// 同一个系统和商户ID下唯一
	Page string `db:"page" json:"page"`
	//标题
	Title string `db:"title" json:"title"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//组件列
	ComponentList FieldsPageComponentList `db:"component_list" json:"componentList"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type FieldsPageComponentList []FieldsPageComponent

// Value sql底层处理器
func (t FieldsPageComponentList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsPageComponentList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsPageComponent struct {
	//组件ID
	ComponentID int64 `db:"component_id" json:"componentID"`
	//组件标识码
	ComponentMark string `db:"component_mark" json:"componentMark"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Value sql底层处理器
func (t FieldsPageComponent) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsPageComponent) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
