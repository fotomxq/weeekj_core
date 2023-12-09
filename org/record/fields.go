package OrgRecord

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

//FieldsRecordType 操作日志
type FieldsRecordType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//模块标识码
	FromMark string `db:"from_mark" json:"fromMark"`
	//修改内容ID
	FromID int64 `db:"from_id" json:"fromID"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//操作内容标识码
	// 可用于其他语言处理
	ContentMark string `db:"content_mark" json:"contentMark"`
	//操作内容概述
	Content string `db:"content" json:"content"`
	//变动内容列
	ChangeData FieldsRecordChangeList `db:"change_data" json:"changeData"`
}

type FieldsRecordChangeList []FieldsRecordChange

//Value sql底层处理器
func (t FieldsRecordChangeList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsRecordChangeList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsRecordChange struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//变动前
	Before string `db:"before" json:"before"`
	//变动后
	After string `db:"after" json:"after"`
}

//Value sql底层处理器
func (t FieldsRecordChange) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsRecordChange) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
