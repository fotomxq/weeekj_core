package OrgTime

import (
	BaseSystemMission "github.com/fotomxq/weeekj_core/v5/base/system_mission"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"sync"
	"time"
)

var (
	//是否需要重新拉取数据？
	runUpdateCacheLock      sync.Mutex
	needUpdateConfigMemData = false
	allConfigList           []FieldsWorkTime
	//调度任务
	runUpdateSysM = BaseSystemMission.Mission{
		OrgID:    0,
		Name:     "组织考勤更新服务",
		Mark:     "org_time.run.update",
		NextTime: "每1秒",
		Bind: BaseSystemMission.MissionBind{
			NatsCode: "org_time_run_update",
			NatsMsg:  "/org/time/run_update",
		},
	}
	//OpenSub 是否启动订阅
	OpenSub = false
)

func Init() {
	if OpenSub {
		subNats()
		BaseSystemMission.ReginWait(&runUpdateSysM, time.Time{})
	}
}

// 将数据加载到内存
func loadAllConfigToMem() {
	//锁定机制
	runUpdateCacheLock.Lock()
	defer runUpdateCacheLock.Unlock()
	//遍历数据
	if needUpdateConfigMemData || len(allConfigList) < 1 {
		//将数据集合写入0
		allConfigList = []FieldsWorkTime{}
		//遍历数据
		limit := 1000
		step := 0
		for {
			var rawList []FieldsWorkTime
			if err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_work_time WHERE expire_at >= NOW() OR expire_at < to_timestamp(1000000) ORDER BY id DESC LIMIT $1 OFFSET $2", limit, step); err != nil {
				break
			}
			if len(rawList) < 1 {
				if step < 1 {
					time.Sleep(time.Second * 10)
				}
				break
			}
			for _, v := range rawList {
				vData := getConfigByID(v.ID)
				if vData.ID < 1 {
					continue
				}
				allConfigList = append(allConfigList, v)
			}
			step += limit
		}
	}
	if len(allConfigList) > 0 {
		needUpdateConfigMemData = false
	}
}
