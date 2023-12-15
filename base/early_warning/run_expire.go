package BaseEarlyWarning

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
)

// 清理过期数据
func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("early warning run, ", r)
		}
	}()
	//处理数据
	if err := clearSendExpire(); err != nil {
		CoreLog.Error("early warning run, cannot clear send expire data, ", err)
	}
}
