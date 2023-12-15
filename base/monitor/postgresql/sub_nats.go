package BaseMonitorPostgresql

import (
	"fmt"
	BaseSystemMission "github.com/fotomxq/weeekj_core/v5/base/system_mission"
	CoreCache "github.com/fotomxq/weeekj_core/v5/core/cache"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func subNats() {
	//任务调度，分析数据
	BaseSystemMission.ReginSub(&runSysM, subNatsRun)
}

func subNatsRun() {
	//日志
	logAppend := "base monitor postgresql run, "
	//捕捉异常
	defer func() {
		//进度控制
		runSysM.Bind.UpdateNextAtFutureSec(runSec)
		//跳出处理
		if r := recover(); r != nil {
			runSysM.Update(fmt.Sprint("发生错误: ", r), "run.error", 0)
			CoreLog.Error(logAppend, r)
		}
	}()
	//拷贝处理处理
	var copyLogList []CoreSQL2.WaitLogType
	CoreSQL2.WaitLogLock.Lock()
	copy(copyLogList, CoreSQL2.WaitLog)
	CoreSQL2.WaitLog = []CoreSQL2.WaitLogType{}
	CoreSQL2.WaitLogLock.Unlock()
	//如果日志不存在，跳出
	if copyLogList == nil || len(copyLogList) < 1 {
		copyLogList = []CoreSQL2.WaitLogType{}
	}
	//进度监听
	runSysM.Start("开始分析", "start", int64(len(copyLogList)))
	//遍历处理
	var analysisData FieldsAnalysis
	//获取当前统计数据
	err := Router2SystemConfig.MainCache.GetStruct(cacheAnalysisKey, &analysisData)
	if err != nil {
		analysisData = FieldsAnalysis{}
	}
	//获取链接池总数
	_ = Router2SystemConfig.MainDB.Get(&analysisData.ConnectMaxCount, "show max_connections")
	_ = Router2SystemConfig.MainDB.Get(&analysisData.ConnectCount, "select count(*) from pg_stat_activity")
	//遍历日志数据
	var copyLogListAny []any
	for k := 0; k < len(copyLogList); k++ {
		v := copyLogList[k]
		//拆分统计
		if v.IsBegin {
			analysisData.AllBeginCount += 1
			switch v.Action {
			case "select":
				analysisData.SelectBeginCount += 1
			case "get":
				analysisData.GetBeginCount += 1
			case "insert":
				analysisData.InsertBeginCount += 1
			case "update":
				analysisData.UpdateBeginCount += 1
			case "delete":
				analysisData.DeleteBeginCount += 1
			case "analysis":
				analysisData.AnalysisBeginCount += 1
			}
		} else {
			analysisData.AllCount += 1
			switch v.Action {
			case "select":
				analysisData.SelectCount += 1
			case "get":
				analysisData.GetCount += 1
			case "insert":
				analysisData.InsertCount += 1
			case "update":
				analysisData.UpdateCount += 1
			case "delete":
				analysisData.DeleteCount += 1
			case "analysis":
				analysisData.AnalysisCount += 1
			}
		}
		if v.Err != nil {
			analysisData.AllErrCount += 1
			switch v.Action {
			case "select":
				analysisData.SelectErrCount += 1
			case "get":
				analysisData.GetErrCount += 1
			case "insert":
				analysisData.InsertErrCount += 1
			case "update":
				analysisData.UpdateErrCount += 1
			case "delete":
				analysisData.DeleteErrCount += 1
			case "analysis":
				analysisData.AnalysisErrCount += 1
			}
		}
		//追加到列队
		copyLogListAny = append(copyLogListAny, FieldsData{
			Action:     v.Action,
			Msg:        v.Msg,
			IsBegin:    v.IsBegin,
			StartAt:    v.StartAt,
			EndAt:      v.EndAt,
			RunSec:     v.RunSec,
			ResultSize: v.ResultSize,
			Err:        v.Err,
		})
	}
	//将日志全部加入列队中
	if len(copyLogListAny) > 0 {
		Router2SystemConfig.MainCache.AppendList(cacheDataKey, copyLogListAny...)
		//检查长度，如果超出10000条自动删除开头的一部分
		cacheLen, _ := Router2SystemConfig.MainCache.GetListLen(cacheDataKey)
		if cacheLen > 10000 {
			for k := 0; k < cacheLen-10000; k++ {
				Router2SystemConfig.MainCache.DeleteListFirst(cacheDataKey)
			}
		}
	}
	//记录统计数据
	Router2SystemConfig.MainCache.SetStruct(cacheAnalysisKey, analysisData, CoreCache.CacheTime1Year)
	//进度处理
	runSysM.Finish()
}
