package CoreSQLConfig

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
)

// FieldsConfigType 简化得扩展结构
type FieldsConfigType struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//值
	Val string `db:"val" json:"val"`
}

// Value sql底层处理器
func (t FieldsConfigType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// FieldsConfigsType 简化得扩展结构
type FieldsConfigsType []FieldsConfigType

// Value sql底层处理器
func (t FieldsConfigsType) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigsType) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

// CheckVal 检查一组数据内的值是否匹配
func (t FieldsConfigsType) CheckVal(mark string, val string) (nowVal string, b bool) {
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
func (t FieldsConfigsType) GetVal(mark string) (val string, b bool) {
	for _, v := range t {
		if v.Mark == mark {
			val = v.Val
			b = true
			return
		}
	}
	return
}

func (t FieldsConfigsType) GetValNoErr(mark string) (val string) {
	for _, v := range t {
		if v.Mark == mark {
			val = v.Val
			return
		}
	}
	return
}

func (t FieldsConfigsType) GetValNoBool(mark string) (val string) {
	for _, v := range t {
		if v.Mark == mark {
			val = v.Val
			return
		}
	}
	return
}

// GetValFloat64 抽取指定的值
func (t FieldsConfigsType) GetValFloat64(mark string) (val float64, b bool) {
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
func (t FieldsConfigsType) GetValInt(mark string) (val int, b bool) {
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

// GetValIntNoBool 抽取指定的值
func (t FieldsConfigsType) GetValIntNoBool(mark string) (val int) {
	val, _ = t.GetValInt(mark)
	return
}

// GetValInt64 抽取指定的值
func (t FieldsConfigsType) GetValInt64(mark string) (val int64, b bool) {
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

func (t FieldsConfigsType) GetValInt64NoBool(mark string) (val int64) {
	val, _ = t.GetValInt64(mark)
	return
}

// GetValBool 抽取指定的值
func (t FieldsConfigsType) GetValBool(mark string) (val bool, b bool) {
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
func Set(data FieldsConfigsType, mark string, val interface{}) FieldsConfigsType {
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
	data = append(data, FieldsConfigType{
		Mark: mark,
		Val:  valStr,
	})
	return data
}
