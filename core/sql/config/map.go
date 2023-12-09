package CoreSQLConfig

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

//扩展结构
type FieldsMap map[string]string

//sql底层处理器
func (t FieldsMap) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
