package CoreSQLGPS

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//一组坐标系
type FieldsPoints []FieldsPoint

//sql底层处理器
func (t FieldsPoints) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsPoints) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

//坐标设计
type FieldsPoint struct {
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

//sql底层处理器
func (t FieldsPoint) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsPoint) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}