package BaseQiniuClient

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base qiniu run, ", r)
		}
	}()
	//延迟10秒启动
	time.Sleep(time.Second * 5)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 1h", func() {
		if runConfigLock {
			return
		}
		runConfigLock = true
		ak, sk, err := getKey()
		if err != nil {
			CoreLog.Error("get qiniu ak and sk failed, ", err)
			return
		}
		if qiniuAK != ak {
			qiniuAK = ak
		}
		if qiniuSK != sk {
			qiniuSK = sk
		}
		runConfigLock = false
	}); err != nil {
		CoreLog.Error("base qiniu run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
