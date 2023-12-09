package MapRoom

import (
	"github.com/lib/pq"
	"time"
)

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID"`
	//任务ID
	// 可能为0
	MissionID int64 `db:"mission_id" json:"missionID"`
	//状态
	// 0 进入; 1 退出; 2 退房中(核对清理状态); 3 清理中; 4 呼叫中; 5 已经应答并发出任务; 6 任务处理完成; 7 应急呼叫器按下; 8 应急呼叫器处理
	Status int `db:"status" json:"status"`
	//入驻人员列
	Infos pq.Int64Array `db:"infos" json:"infos"`
	//备注
	Des string `db:"des" json:"des"`
}
