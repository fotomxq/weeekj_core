package BaseCache

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	"github.com/robfig/cron"
	"time"
)

// Run 将自动遍历，删除过期数据
// Deprecated
func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("cache run, ", r)
		}
	}()
	time.Sleep(time.Second * 10)
	//启动定时器
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 1s", func() {
		if runCacheLock {
			return
		}
		runCacheLock = true
		//找到过期数据，重组
		nowTime := CoreFilter.GetNowTime().Unix()
		cacheLock.Lock()
		var newCacheData []DataCache
		for _, v := range cacheData {
			if v.ExpireTime < nowTime {
				continue
			}
			newCacheData = append(newCacheData, v)
		}
		cacheData = newCacheData
		cacheLock.Unlock()
		runCacheLock = false
	}); err != nil {
		CoreLog.Error("cache run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
