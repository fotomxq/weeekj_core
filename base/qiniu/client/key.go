package BaseQiniuClient

import (
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
)

// 获取key数据
func getKey() (string, string, error) {
	ak, err := BaseConfig.GetDataString("FileQiniuAK")
	if err != nil {
		return "", "", errors.New("config load FileQiniuAK, " + err.Error())
	}
	sk, err := BaseConfig.GetDataString("FileQiniuSK")
	if err != nil {
		return "", "", errors.New("config load FileQiniuSK, " + err.Error())
	}
	return ak, sk, nil
}
