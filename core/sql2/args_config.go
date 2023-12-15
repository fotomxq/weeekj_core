package CoreSQL2

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

type ArgsConfigList []ArgsConfig

func (t ArgsConfigList) GetFields() FieldsConfigList {
	data := FieldsConfigList{}
	for _, v := range t {
		data = append(data, v.GetField())
	}
	return data
}

// ArgsConfig 扩展结构
type ArgsConfig struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//值
	Val string `db:"val" json:"val"`
}

func (t *ArgsConfig) GetField() FieldsConfig {
	return FieldsConfig{
		Mark: t.Mark,
		Val:  t.Val,
	}
}

// FieldsConfig 扩展结构
type FieldsConfig struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
	//值
	Val string `db:"val" json:"val"`
}

// Value sql底层处理器
func (t FieldsConfig) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfig) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsConfigList 简化得扩展结构
type FieldsConfigList []FieldsConfig

// Value sql底层处理器
func (t FieldsConfigList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// CheckVal 检查一组数据内的值是否匹配
func (t FieldsConfigList) CheckVal(mark string, val string) (nowVal string, b bool) {
	for _, v := range t {
		if v.Mark == mark {
			if v.Val == val {
				return val, true
			} else {
				return val, false
			}
		}
	}
	return
}

// GetVal 抽取指定的值
func (t FieldsConfigList) GetVal(mark string) (val string, b bool) {
	for _, v := range t {
		if v.Mark == mark {
			val = v.Val
			b = true
			return
		}
	}
	return
}

func (t FieldsConfigList) GetValNoErr(mark string) (val string) {
	for _, v := range t {
		if v.Mark == mark {
			val = v.Val
			return
		}
	}
	return
}

func (t FieldsConfigList) GetValNoBool(mark string) (val string) {
	for _, v := range t {
		if v.Mark == mark {
			val = v.Val
			return
		}
	}
	return
}

// GetValFloat64 抽取指定的值
func (t FieldsConfigList) GetValFloat64(mark string) (val float64, b bool) {
	var valStr string
	valStr, b = t.GetVal(mark)
	if !b {
		return
	}
	var err error
	val, err = CoreFilter.GetFloat64ByString(valStr)
	if err != nil {
		return
	}
	b = true
	return
}

// GetValInt 抽取指定的值
func (t FieldsConfigList) GetValInt(mark string) (val int, b bool) {
	var valStr string
	valStr, b = t.GetVal(mark)
	if !b {
		return
	}
	var err error
	val, err = CoreFilter.GetIntByString(valStr)
	if err != nil {
		return
	}
	b = true
	return
}

// GetValInt64 抽取指定的值
func (t FieldsConfigList) GetValInt64(mark string) (val int64, b bool) {
	var valStr string
	valStr, b = t.GetVal(mark)
	if !b {
		return
	}
	var err error
	val, err = CoreFilter.GetInt64ByString(valStr)
	if err != nil {
		return
	}
	b = true
	return
}

func (t FieldsConfigList) GetValInt64NoBool(mark string) (val int64) {
	val, _ = t.GetValInt64(mark)
	return
}

// GetValBool 抽取指定的值
func (t FieldsConfigList) GetValBool(mark string) (val bool, b bool) {
	var valStr string
	valStr, b = t.GetVal(mark)
	if !b {
		return
	}
	var err error
	val, err = CoreFilter.GetBoolByInterface(valStr)
	if err != nil {
		return
	}
	b = true
	return
}

// Set 写入数据
func Set(data FieldsConfigList, mark string, val interface{}) FieldsConfigList {
	valStr, err := CoreFilter.GetStringByInterface(val)
	if err != nil {
		valStr = ""
	}
	for k, v := range data {
		if v.Mark == mark {
			data[k].Val = valStr
			return data
		}
	}
	data = append(data, FieldsConfig{
		Mark: mark,
		Val:  valStr,
	})
	return data
}
