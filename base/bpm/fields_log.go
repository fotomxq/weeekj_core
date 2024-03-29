package BaseBPM

import "time"

// FieldsLog 流程节点日志
// 每个日志不可变更，只能创建迭代
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreatedAt time.Time `db:"create_at" json:"createdAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//管理单元
	UnitID int64 `db:"unit_id" json:"unitID"`
	//操作用户
	UserID int64 `db:"user_id" json:"userID"`
	//操作组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//BPM ID
	BPMID int64 `db:"bpm_id" json:"bpmId" check:"id"`
	//当前节点ID
	NodeID string `db:"node_id" json:"nodeId" check:"des" min:"1" max:"300"`
	//当前节点序号
	NodeNum int `db:"node_num" json:"nodeNum" check:"int64Than0"`
	//节点存储内容
	NodeContent string `db:"node_content" json:"nodeContent"`
}
