package BaseQiniu

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
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

// 获取call back
func getCallBack() (string, error) {
	appAPI, err := BaseConfig.GetDataString("AppAPI")
	if err != nil {
		return "", errors.New("config load FileQiNiuCallBackURL, " + err.Error())
	}
	return fmt.Sprint(appAPI, "/v2/base/file/public/qiniu/callback"), nil
}

// 获取是否需要回调函数
func getCallbackOn() (bool, error) {
	b, err := BaseConfig.GetDataBool("FileQiNiuCallBackURLON")
	if err != nil {
		return false, errors.New("config load FileQiNiuCallBackURLON, " + err.Error())
	}
	return b, nil
}
