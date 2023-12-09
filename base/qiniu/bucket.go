package BaseQiniu

import (
	"errors"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	"strings"
)

// 获取桶名称
// 自动检查是否支持？
// param bucketConfigName string 配置名
func checkBucketName(name string) bool {
	bucketList, err := getBucketNameList()
	if err != nil {
		return false
	}
	//分割结束后检查是否在列表？
	for _, v := range bucketList {
		if v == "" {
			continue
		}
		if v == name {
			return true
		}
	}
	return false
}

// 获取桶列表
func getBucketNameList() ([]string, error) {
	//从配置QiNiuBucketList获取数据，并检查是否在范围内？
	bucketListStr, err := BaseConfig.GetDataString("FileQiNiuBucketList")
	if err != nil {
		return []string{}, errors.New("get config, " + err.Error())
	}
	//分割
	bucketList := strings.Split(bucketListStr, "|")
	return bucketList, nil
}

// 获取桶对应URL
func getBucketURL(bucketName string) (string, error) {
	bucketNameList, err := getBucketNameList()
	if err != nil {
		return "", err
	}
	bucketURLList, err := BaseConfig.GetDataString("FileQiNiuURLList")
	if err != nil {
		return "", errors.New("get config, " + err.Error())
	}
	bucketURLarr := strings.Split(bucketURLList, "|")
	for k, v := range bucketNameList {
		if v != bucketName {
			continue
		}
		for k2, v2 := range bucketURLarr {
			if k == k2 {
				if v2 == "" {
					return "", errors.New("bucket not have url")
				}
				return v2, nil
			}
		}
	}
	return "", errors.New("cannot find bucket url")
}
