package BaseSystemMission

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	"time"
)

// IsStart 检查是否正在执行
func (t *MissionBind) IsStart() bool {
	return t.isRun
}

// Start 开始执行
func (t *MissionBind) Start() {
	t.isRun = true
	CoreNats.PushDataNoErr(t.NatsMsg, "", 0, "", nil)
}

// NeedStart 是否需要开始执行
func (t *MissionBind) NeedStart() bool {
	if t.IsStart() {
		return false
	}
	if CoreSQL.CheckTimeHaveData(t.nextAt) && CoreSQL.CheckTimeThanNow(t.nextAt) {
		return false
	}
	return true
}

func (t *MissionBind) Finish() {
	t.isRun = false
}

// UpdateNextAt 更新下一次执行时间
func (t *MissionBind) UpdateNextAt(nextAt time.Time) {
	t.nextAt = nextAt
}

// UpdateNextAtFutureSec 更新下一次执行时间未来几秒
func (t *MissionBind) UpdateNextAtFutureSec(sec int) {
	t.UpdateNextAt(CoreFilter.GetNowTimeCarbon().AddSeconds(sec).Time)
}

// UpdateNextAtFutureHour 更新下一次执行时间到明天指定时间
func (t *MissionBind) UpdateNextAtFutureHour(hour, minute, sec int) {
	t.UpdateNextAt(CoreFilter.GetNowTimeCarbon().AddDay().SetHour(hour).SetMinute(minute).SetSecond(sec).Time)
}

// UpdateNextAtFutureDay 更新下一次执行时间到未来某一天
func (t *MissionBind) UpdateNextAtFutureDay(day, hour, minute, sec int) {
	t.UpdateNextAt(CoreFilter.GetNowTimeCarbon().AddDays(day).SetHour(hour).SetMinute(minute).SetSecond(sec).Time)
}

// UpdateNextAtFutureMonth 更新下一次执行时间到未来某月
func (t *MissionBind) UpdateNextAtFutureMonth(month, day, hour, minute, sec int) {
	t.UpdateNextAt(CoreFilter.GetNowTimeCarbon().AddMonths(month).AddDays(day).SetHour(hour).SetMinute(minute).SetSecond(sec).Time)
}
