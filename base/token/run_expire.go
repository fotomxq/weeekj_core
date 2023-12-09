package BaseToken

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
)

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("token run, ", r)
		}
	}()
	DeleteByExpire()
}
