package BaseQiniu

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
)

// 为文件构建新的名称
func getFileName(data []byte, fileType string) (string, error) {
	var nowTime = CoreFilter.GetNowTime().Format("2006-01-02_15-04-05")
	fileSha1, err := CoreFilter.GetSha1(data)
	if err != nil {
		return "", err
	}
	return nowTime + "_" + string(fileSha1) + "." + fileType, nil
}
