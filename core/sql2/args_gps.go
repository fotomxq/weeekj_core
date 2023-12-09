package CoreSQL2

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ArgsGPS struct {
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// GetFields 获取SQL数据
func (t ArgsGPS) GetFields() FieldsGPS {
	return FieldsGPS{
		Longitude: t.Longitude,
		Latitude:  t.Latitude,
	}
}

// FieldsGPS 坐标设计
type FieldsGPS struct {
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
}

// Value sql底层处理器
func (t FieldsGPS) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsGPS) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
