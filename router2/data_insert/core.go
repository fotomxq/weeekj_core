package Router2DataInsert

import BaseQiniu "gitee.com/weeekj/weeekj_core/v5/base/qiniu"

//快速数据填充模块

// 获取头像地址
func getURLByFileID(fileID int64) string {
	if fileID < 1 {
		return ""
	}
	urlStr, _ := BaseQiniu.GetPublicURLStr(fileID)
	return urlStr
}
