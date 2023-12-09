package RouterSystem

import (
	"fmt"
	CoreFile "gitee.com/weeekj/weeekj_core/v5/core/file"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"time"
)

// Close 系统总的控制器
// 关闭系统
func Close() {
	//请求关闭的文件路径
	closeSysSrc := fmt.Sprint(CoreFile.BaseSrc, CoreFile.Sep, "close_sys")
	//启动维护程序
	for {
		//检查是否存在关闭请求
		if CoreFile.IsFile(closeSysSrc) {
			if err := CoreFile.DeleteF(closeSysSrc); err != nil {
				CoreLog.Error("cannot delete system close file, ", err)
			}
			CoreSystemClose.Close()
			CoreLog.Info("system is close")
			time.Sleep(time.Second * 1)
			return
		}
		//等待1秒
		time.Sleep(time.Second * 1)
	}
}
