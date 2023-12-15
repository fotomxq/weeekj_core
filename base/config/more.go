package BaseConfig

import CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"

// GetDataInt 扩展支持，直接转化对应的数据值
func GetDataInt(mark string) (data int, err error) {
	var configData string
	configData, err = GetData(&ArgsGetData{
		Mark: mark,
	})
	if err != nil {
		return
	}
	return CoreFilter.GetIntByString(configData)
}

func GetDataInt64(mark string) (data int64, err error) {
	var configData string
	configData, err = GetData(&ArgsGetData{
		Mark: mark,
	})
	if err != nil {
		return
	}
	return CoreFilter.GetInt64ByString(configData)
}

func GetDataInt64NoErr(mark string) (data int64) {
	var err error
	data, err = GetDataInt64(mark)
	if err != nil {
		return
	}
	return
}

func GetDataString(mark string) (data string, err error) {
	var configData string
	configData, err = GetData(&ArgsGetData{
		Mark: mark,
	})
	if err != nil {
		return
	}
	return configData, err
}

func GetDataStringNoErr(mark string) (data string) {
	data, _ = GetDataString(mark)
	return
}

func GetDataFloat64(mark string) (data float64, err error) {
	var configData string
	configData, err = GetData(&ArgsGetData{
		Mark: mark,
	})
	if err != nil {
		return
	}
	return CoreFilter.GetFloat64ByString(configData)
}

func GetDataBool(mark string) (data bool, err error) {
	var configData string
	configData, err = GetData(&ArgsGetData{
		Mark: mark,
	})
	if err != nil {
		return
	}
	return configData == "true", err
}

func GetDataBoolNoErr(mark string) (data bool) {
	var err error
	data, err = GetDataBool(mark)
	if err != nil {
		data = false
		return
	}
	return
}
