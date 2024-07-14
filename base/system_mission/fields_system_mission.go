package BaseSystemMission

import "time"

// FieldsMission 任务记录
type FieldsMission struct {
	//ID
	ID int64 `db:"id" json:"id" unique:"true"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//组织ID
	// 如果为0则为系统服务
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true" index:"true"`
	//任务名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300" empty:"true"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//开始时间
	StartAt time.Time `db:"start_at" json:"startAt" default:"0"`
	//当前执行提示
	NowTip string `db:"now_tip" json:"nowTip" check:"des" min:"1" max:"1000" empty:"true"`
	//停止时间
	StopAt time.Time `db:"stop_at" json:"stopAt" default:"0"`
	//暂停时间
	PauseAt time.Time `db:"pause_at" json:"pauseAt" default:"0"`
	//暂停位置
	Location string `db:"location" json:"location" check:"des" min:"1" max:"1000" empty:"true"`
	//总数量
	AllCount int64 `db:"all_count" json:"allCount"`
	//已经执行数量
	RunCount int64 `db:"run_count" json:"runCount"`
	//总消耗时间秒
	RunAllSec int64 `db:"run_all_sec" json:"runAllSec"`
	//计划执行时间
	NextTime string `db:"next_time" json:"nextTime"`
}
