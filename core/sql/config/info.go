package CoreSQLConfig

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//FieldsInfoType 扩展结构
type FieldsInfoType struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//值
	Val string `db:"val" json:"val"`
}

//Value sql底层处理器
func (t FieldsInfoType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsInfoType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsInfosType []FieldsInfoType

//Value sql底层处理器
func (t FieldsInfosType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsInfosType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
