package CoreFilter

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

//转化模块
// 将某些类型转为其他类型
// 总结最常用的转化方案

func GetBoolByInterface(data interface{}) (bool, error) {
	if value, ok := data.([]byte); ok {
		if len(value) < 1 {
			return false, errors.New("data is empty")
		}
		return value[0] == 0x01, nil
	} else if value, ok := data.(string); ok {
		return value == "1" || value == "true", nil
	} else if value, ok := data.(bool); ok {
		return value, nil
	} else {
		return false, errors.New("type error")
	}
}

func GetBoolByInterfaceNoErr(data interface{}) bool {
	b, _ := GetBoolByInterface(data)
	return b
}

func GetStringByInt(data int) string {
	return strconv.Itoa(data)
}

func GetStringByInt64(data int64) string {
	return strconv.FormatInt(data, 10)
}

func GetStringByUint(data uint) string {
	return fmt.Sprint(data)
}

func GetStringByUint64(data uint64) string {
	return strconv.FormatUint(data, 10)
}

func GetStringByFloat64(data float64) string {
	return strconv.FormatFloat(data, 'f', -1, 64)
}

func GetStringByInterface(data interface{}) (string, error) {
	if value, ok := data.([]byte); ok {
		return fmt.Sprintf("%x", value), nil
	} else if value, ok := data.(string); ok {
		return value, nil
	} else if value, ok := data.(int64); ok {
		return GetStringByInt64(value), nil
	} else if value, ok := data.(int); ok {
		return GetStringByInt(value), nil
	} else if value, ok := data.(float64); ok {
		return GetStringByFloat64(value), nil
	} else if value, ok := data.(bool); ok {
		valStr := "true"
		if !value {
			valStr = "false"
		}
		return valStr, nil
	} else {
		return "", errors.New("type error")
	}
}

func GetIntByString(data string) (int, error) {
	return strconv.Atoi(data)
}

func GetIntByStringNoErr(data string) int {
	r, _ := strconv.Atoi(data)
	return r
}

func GetUIntByString(data string) (uint, error) {
	res, err := strconv.ParseUint(data, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(res), nil
}

func GetUIntByInt(data int) uint {
	return uint(data)
}

func GetIntByFloat64(data float64) int {
	return int(data)
}

func GetIntByInterface(data interface{}) (int, error) {
	if value, ok := data.([]byte); ok {
		return int(binary.BigEndian.Uint32(value)), nil
	} else if value, ok := data.(int); ok {
		return value, nil
	} else if value, ok := data.(int64); ok {
		return int(value), nil
	} else if value, ok := data.(float64); ok {
		return GetIntByFloat64(value), nil
	} else if value, ok := data.(string); ok {
		return GetIntByString(value)
	} else {
		return 0, errors.New("type error")
	}
}

func GetInt64ByFloat64(data float64) int64 {
	return int64(data)
}

func GetInt64ByUint8(data []uint8) (int64, error) {
	vF, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0, err
	}
	return GetInt64ByFloat64(vF), nil
}

func GetInt64ByUint8NoErr(data []uint8) int64 {
	r, _ := GetInt64ByUint8(data)
	return r
}

func GetInt64ByString(data string) (int64, error) {
	return strconv.ParseInt(data, 10, 64)
}

func GetInt64ByStringNoErr(data string) int64 {
	d, err := GetInt64ByString(data)
	if err != nil {
		return 0
	}
	return d
}

func GetInt64ByInterface(data interface{}) (int64, error) {
	if value, ok := data.([]byte); ok {
		return int64(binary.BigEndian.Uint64(value)), nil
	} else if value, ok := data.(int64); ok {
		return value, nil
	} else if value, ok := data.(int); ok {
		return int64(value), nil
	} else if value, ok := data.(float64); ok {
		return GetInt64ByFloat64(value), nil
	} else if value, ok := data.(string); ok {
		return GetInt64ByString(value)
	} else {
		return 0, errors.New("type error")
	}
}

func GetFloat64ByInt(data int) float64 {
	return float64(data)
}

func GetFloat64ByString(data string) (float64, error) {
	return strconv.ParseFloat(data, 64)
}

func GetFloat64ByStringNoErr(data string) float64 {
	d, _ := GetFloat64ByString(data)
	return d
}

func GetFloat64ByInt64(data int64) float64 {
	return float64(data)
}

func GetFloat64ByUint8(data []uint8) (float64, error) {
	return strconv.ParseFloat(string(data), 64)
}

func GetFloat64ByInterface(data interface{}) (float64, error) {
	if value, ok := data.([]byte); ok {
		return GetFloat64ByString(fmt.Sprintf("%x", value))
	} else if value, ok := data.(float64); ok {
		return value, nil
	} else if value, ok := data.(int64); ok {
		return GetFloat64ByInt64(value), nil
	} else if value, ok := data.(int); ok {
		return GetFloat64ByInt(value), nil
	} else if value, ok := data.(string); ok {
		return GetFloat64ByString(value)
	} else if value, ok := data.([]uint8); ok {
		return GetFloat64ByUint8(value)
	} else {
		return 0, errors.New("type error")
	}
}

// GetStructToMap 将struct数组解析为map结构
func GetStructToMap(structData interface{}, mapData *map[string]interface{}) error {
	rawByte, err := json.Marshal(structData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawByte, mapData)
	if err != nil {
		return err
	}
	return nil
}
