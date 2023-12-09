package CoreSQLConfig

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

//按照顺序排序的专用方法
// json转化后，按照提交的顺序进行转化，避免json无序问题

//FieldsConfigSortType 简化得扩展结构
type FieldsConfigSortType struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//值
	Val string `db:"val" json:"val"`
}

//Value sql底层处理器
func (t FieldsConfigSortType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigSortType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

//FieldsConfigSortsType 简化得扩展结构
type FieldsConfigSortsType []FieldsConfigSortType

//Value sql底层处理器
func (t FieldsConfigSortsType) Value() (driver.Value, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'['})
	l := len(t)
	for i, k := range t {
		_, _ = fmt.Fprintf(buf, "{\"mark\":\"%s\",\"val\":\"%v\"}", k.Mark, k.Val)
		if i < l-1 {
			buf.WriteByte(',')
		}
	}
	buf.Write([]byte{']'})
	return buf.Bytes(), nil
}

func (t *FieldsConfigSortsType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}