package IOTMQTT

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSystemClose "gitee.com/weeekj/weeekj_core/v5/core/system_close"
	IOTMission "gitee.com/weeekj/weeekj_core/v5/iot/mission"
	"github.com/robfig/cron"
)

func Run() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device run error, ", r)
		}
		runTimer.Stop()
	}()
	//等待
	//time.Sleep(time.Second * 10)
	//重连处理机制
	runTimer = cron.New()
	if err := runTimer.AddFunc("@every 3s", func() {
		if runConnectLock {
			return
		}
		runConnectLock = true
		if !runConnect() {
			runConnectLock = false
		}
	}); err != nil {
		CoreLog.Error("iot device mqtt run, cron time, ", err)
	}
	if OpenBaseMission {
		//设置阻拦器并启动任务推送
		IOTMission.RunMissionBlocker.SetExpire(10)
		if err := runTimer.AddFunc("@every 1s", func() {
			if runMissionLock {
				return
			}
			if !MQTTIsConnect {
				return
			}
			if !MQTTClient.Client.IsConnected() {
				return
			}
			runMissionLock = true
			runMission()
			runMissionLock = false
		}); err != nil {
			CoreLog.Error("iot device mqtt run, cron time, ", err)
		}
	}
	if err := runTimer.AddFunc("@every 3s", func() {
		if runUpdateDataLock {
			return
		}
		runUpdateDataLock = true
		runUpdateData()
		runUpdateDataLock = false
	}); err != nil {
		CoreLog.Error("iot device mqtt run, cron time, ", err)
	}
	if err := runTimer.AddFunc("@every 1s", func() {
		if runWaitSendLock {
			return
		}
		runWaitSendLock = true
		runWaitSend()
		runWaitSendLock = false
	}); err != nil {
		CoreLog.Error("iot device mqtt run, cron time, ", err)
	}
	runTimer.Start()
	//卡住进程
	CoreSystemClose.Wait()
	//退出时间
	runTimer.Stop()
}
