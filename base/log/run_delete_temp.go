package BaseLog

import (
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 自动删除临时文件
func runDeleteTemp() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base log delete temp run, ", r)
		}
	}()
	downloadFileSrc := fmt.Sprint(Router2SystemConfig.RootDir, CoreFile.Sep+"temp"+CoreFile.Sep, "log")
	if CoreFile.IsFolder(downloadFileSrc) {
		_ = CoreFile.DeleteF(downloadFileSrc)
	}
}
