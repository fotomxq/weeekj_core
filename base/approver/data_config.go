package BaseApprover

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type DataConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//关联的模块标识码
	// erp_project
	ModuleCode string `db:"module_code" json:"moduleCode" check:"des" min:"1" max:"50"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
	//审批分叉标识码
	// 用于识别模块内，不同的审批流程
	ForkCode string `db:"fork_code" json:"forkCode" check:"des" min:"1" max:"50"`
	//审批流配置
	Items DataConfigItems `json:"items"`
}

type DataConfigItems []DataConfigItem

// Value sql底层处理器
func (t DataConfigItems) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *DataConfigItems) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type DataConfigItem struct {
	//审批顺序
	FlowOrder int `db:"flow_order" json:"flowOrder" check:"intThan0" empty:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//审批人用户ID
	// 用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// Value sql底层处理器
func (t DataConfigItem) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *DataConfigItem) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
