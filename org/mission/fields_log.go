package OrgMission

import "time"

//FieldsLog 操作日志
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//任务ID
	MissionID int64 `db:"mission_id" json:"missionID"`
	//操作内容标识码
	// 可用于其他语言处理
	ContentMark string `db:"content_mark" json:"contentMark"`
	//操作内容概述
	Content string `db:"content" json:"content"`
}
