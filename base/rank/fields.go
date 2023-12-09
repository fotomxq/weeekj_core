package BaseRank

import (
	"time"
)

//列队数据表
type FieldsRank struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//提取最短间隔 s
	PickMin int64 `db:"pick_min" json:"pickMin"`
	//下一次允许提取的时间
	// 会根据任务类型，做一些超时锁定，避免被连续提取
	// 此处为下一次提取的时间
	PickAt time.Time `db:"pick_at" json:"pickAt"`
	//服务来源
	ServiceMark string `db:"service_mark" json:"serviceMark"`
	//任务标识码
	MissionMark string `db:"mission_mark" json:"missionMark"`
	//任务内容
	MissionData []byte `db:"mission_data" json:"missionData"`
}

//任务结束列队
type FieldsRankOver struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//结束时间
	OverAt time.Time `db:"over_at" json:"overAt"`
	//完成结果
	Result []byte `db:"result" json:"result"`
	//服务来源
	ServiceMark string `db:"service_mark" json:"serviceMark"`
	//任务标识码
	MissionMark string `db:"mission_mark" json:"missionMark"`
	//任务内容
	MissionData []byte `db:"mission_data" json:"missionData"`
}
