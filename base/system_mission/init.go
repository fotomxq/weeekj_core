package BaseSystemMission

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreHighf "github.com/fotomxq/weeekj_core/v5/core/highf"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	"github.com/golang-module/carbon"
	"github.com/nats-io/nats.go"
	"time"
)

// Mission 初始化方法
type Mission struct {
	//组织ID
	// 如果为0则为系统服务
	OrgID int64 `db:"org_id" json:"orgID"`
	//任务名称
	Name string `db:"name" json:"name"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//计划执行时间
	NextTime string `db:"next_time" json:"nextTime"`
	//当前ID
	nowID int64
	//当前结构体
	nowData FieldsMission
	//更新拦截器
	updateBlock CoreHighf.BlockerWait
	//是否禁止拦截器
	noBlock bool
	//拦截器时间
	blockSec int
	//本次开始时间
	nowStartAt carbon.Carbon
	//本次结束时间
	nowEndAt carbon.Carbon
	//挂靠模式，本模块主动触发形式
	Bind MissionBind
	//禁止block后，内部需要做一定拦截处理，否则可能面临高请求异常
	// 记录上次拦截时间，自动间隔3秒再给数据库更新数据
	blockAutoTimer int64
	//是否完成初始化？
	isStart bool
}

type MissionBind struct {
	//触发消息地址
	NatsMsg string
	//下一轮执行时间
	nextAt time.Time
	//是否正在执行
	isRun bool
}

// init 初始化
func (t *Mission) init() {
	//初始化ID
	err := createMission(&argsCreateMission{
		OrgID:    t.OrgID,
		Name:     t.Name,
		Mark:     t.Mark,
		NextTime: t.NextTime,
	})
	if err != nil {
		CoreLog.Error("core system mission: ", t.Name, ", mark: ", t.Mark, ", create mission failed, err: ", err)
	}
	//通过mark获取数据
	t.nowData = getMissionByMark(t.Mark, t.OrgID)
	t.nowID = t.nowData.ID
	//如果尚未初始化，则继续
	if t.isStart {
		return
	}
	t.isStart = true
	//订阅nats
	CoreNats.SubDataByteNoErr("/base/system_mission/stop", func(_ *nats.Msg, _ string, id int64, _ string, _ []byte) {
		if id != t.nowID {
			return
		}
		t.updateData()
	})
	//拦截器初始化
	if t.blockSec < 0 {
		t.blockSec = 3
	}
	t.updateBlock.Init(t.blockSec)
}

// 修改拦截器时间
func (t *Mission) UpdateBlockTime(sec int) {
	t.blockSec = sec
	if sec < 1 {
		t.noBlock = true
		t.updateBlock.Init(1)
	} else {
		t.noBlock = false
		t.updateBlock.Init(sec)
	}
}

// Do 触发器
func (t *Mission) Do() {

}

// Start 开始任务
func (t *Mission) Start(nowTip string, location string, allCount int64) {
	if t.nowID < 1 {
		t.init()
	}
	t.nowStartAt = CoreFilter.GetNowTimeCarbon()
	_ = startMission(t.nowID, nowTip, location, allCount)
	t.updateData()
}

// Update 更新执行情况
func (t *Mission) Update(nowTip string, location string, runCount int64) {
	isEnd := false
	if runCount >= t.nowData.AllCount {
		if runCount > t.nowData.AllCount {
			t.UpdateAddTotal(runCount - t.nowData.AllCount)
		}
		isEnd = true
	}
	if t.noBlock && !isEnd {
		nowAtUnix := CoreFilter.GetNowTime().Unix()
		if t.blockAutoTimer+3 > nowAtUnix {
			return
		}
		t.blockAutoTimer = nowAtUnix
	}
	if t.noBlock {
		err := updateMission(t.nowID, nowTip, location, runCount, 0)
		if err != nil {
			CoreLog.Error("core system mission: ", t.Name, ", mark: ", t.Mark, ", update mission failed, err: ", err)
		}
		t.updateData()
	} else {
		t.updateBlock.CheckWait(0, "", func(_ int64, _ string) {
			err := updateMission(t.nowID, nowTip, location, runCount, 0)
			if err != nil {
				CoreLog.Error("core system mission: ", t.Name, ", mark: ", t.Mark, ", update mission failed, err: ", err)
			}
			t.updateData()
		})
	}
}

// UpdateTotal 更新总量情况
func (t *Mission) UpdateTotal(allCount int64) {
	err := updateMissionTotal(t.nowID, allCount)
	if err != nil {
		CoreLog.Error("core system mission: ", t.Name, ", mark: ", t.Mark, ", update mission total failed, err: ", err)
	}
}

// UpdateAddTotal 更新总量情况
func (t *Mission) UpdateAddTotal(allCount int64) {
	err := updateMissionAddTotal(t.nowID, allCount)
	if err != nil {
		CoreLog.Error("core system mission: ", t.Name, ", mark: ", t.Mark, ", update mission total failed, err: ", err)
	}
}

// Finish 完成一轮执行
func (t *Mission) Finish() {
	data := getMission(t.nowID)
	t.nowEndAt = CoreFilter.GetNowTimeCarbon()
	runSec := t.nowEndAt.DiffInSecondsWithAbs(t.nowStartAt)
	err := updateMission(t.nowID, "完成处理", "end", data.AllCount, runSec)
	if err != nil {
		CoreLog.Error("core system mission: ", t.Name, ", mark: ", t.Mark, ", finish mission failed, err: ", err)
	}
	t.updateData()
	t.Bind.Finish()
}

// IsStop 是否停止
func (t *Mission) IsStop() bool {
	if t.nowID < 1 {
		return false
	}
	return CoreSQL.CheckTimeHaveData(t.nowData.StopAt)
}

// Stop 停止任务
func (t *Mission) Stop() {
	_ = stopMission(t.nowID)
	t.updateData()
}

// Pause 暂停任务
func (t *Mission) Pause() {
	_ = pauseMission(t.nowID)
	t.updateData()
}

// updateData 更新结构体
func (t *Mission) updateData() {
	t.nowData = getMission(t.nowID)
}
