package BaseWeixin

import (
	BaseWeixinWXXMessage "github.com/fotomxq/weeekj_core/v5/base/weixin/wxx/message"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
)

// Run 打开启动服务
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("weixin message template run error, ", r)
		}
	}()
	BaseWeixinWXXMessage.Run()
}
