package OrgTime

import (
	"fmt"
	BaseSystemMission "github.com/fotomxq/weeekj_core/v5/base/system_mission"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func subNats() {
	//任务调度，更新考勤情况
	BaseSystemMission.ReginSub(&runUpdateSysM, subNatsRunUpdate)
}

func subNatsRunUpdate() {
	//日志
	logAppend := "org time run, "
	//捕捉异常
	defer func() {
		//进度控制
		runUpdateSysM.Bind.UpdateNextAtFutureSec(1)
		//跳出处理
		if r := recover(); r != nil {
			runUpdateSysM.Update(fmt.Sprint("发生错误: ", r), "run.error", 0)
			CoreLog.Error(logAppend, r)
		}
	}()
	//确保数据全部在内存内
	loadAllConfigToMem()
	//锁定机制
	runUpdateCacheLock.Lock()
	defer runUpdateCacheLock.Unlock()
	//跟踪器
	runUpdateSysM.Start("开始更新", "start", int64(len(allConfigList)))
	//遍历数据处理
	for k, v := range allConfigList {
		//跟踪器
		runUpdateSysM.Update(fmt.Sprint("正在处理", v.Name), fmt.Sprint("config:", v.ID), int64(k))
		//检查是否上班
		checkIsWork := checkIsWorkByData(&v)
		if checkIsWork == v.IsWork {
			continue
		}
		//切换轮动机制
		if len(v.RotConfig.WorkTime) > 0 && checkIsWork == false {
			//更换轮动次序
			v.RotConfig.NowKey += 1
			if v.RotConfig.NowKey < 0 {
				v.RotConfig.NowKey = 0
			}
			if v.RotConfig.NowKey >= len(v.RotConfig.WorkTime) {
				v.RotConfig.NowKey = 0
			}
			//标记切换上下班状态
			if _, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_work_time SET is_work = :is_work, rot_config = :rot_config WHERE id = :id", map[string]interface{}{
				"id":         v.ID,
				"is_work":    checkIsWork,
				"rot_config": v.RotConfig,
			}); err != nil {
				CoreLog.Error(logAppend, "update failed, ", err)
				continue
			} else {
				//修改内存数据
				allConfigList[k].IsWork = checkIsWork
				allConfigList[k].RotConfig = v.RotConfig
			}
		} else {
			//标记切换上下班状态
			if _, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_work_time SET is_work = :is_work WHERE id = :id", map[string]interface{}{
				"id":      v.ID,
				"is_work": checkIsWork,
			}); err != nil {
				CoreLog.Error(logAppend, "update failed, ", err)
				continue
			} else {
				//修改内存数据
				allConfigList[k].IsWork = checkIsWork
			}
		}
		//推送nats通知
		pushNatsWork(v.OrgID, v.Binds, v.Groups, checkIsWork)
		//删除缓冲
		deleteConfigCache(v.ID)
	}
	//收尾处理
	runUpdateSysM.Finish()
}
