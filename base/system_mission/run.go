package BaseSystemMission

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
	"time"
)

//主动触发器
/**
1. 本模块可以挂靠函数，根据定时规划自动触发
*/

var (
	waitMission []*Mission
)

// ReginWait 注册新的挂靠
func ReginWait(mission *Mission, nextAt time.Time) {
	mission.Bind.UpdateNextAt(nextAt)
	for k, v := range waitMission {
		if v.Mark == mission.Mark {
			waitMission[k] = mission
			return
		}
	}
	waitMission = append(waitMission, mission)
}

// ReginSub 快速订阅
func ReginSub(mission *Mission, handle func()) {
	CoreNats.SubDataByteNoErr(mission.Bind.NatsMsg, func(_ *nats.Msg, _ string, _ int64, _ string, _ []byte) {
		handle()
	})
}

// Run 调度程序保护器
func Run() {
	time.Sleep(time.Second * 5)
	for {
		runChild()
		time.Sleep(time.Minute * 1)
	}
}

func runChild() {
	//日志
	appendLog := "base system mission run, "
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error(appendLog, r)
		}
	}()
	for {
		//等待1秒继续
		time.Sleep(time.Second * 1)
		//预备跳过
		if waitMission == nil || len(waitMission) < 1 {
			continue
		}
		//遍历所有服务
		for k := 0; k < len(waitMission); k++ {
			//获取当前节点
			v := waitMission[k]
			//跳过正在执行的服务
			if v.Bind.IsStart() {
				continue
			}
			//检查时间是否已经达到
			if !v.Bind.NeedStart() {
				continue
			}
			//触发执行
			v.Bind.Start()
		}
	}
}
