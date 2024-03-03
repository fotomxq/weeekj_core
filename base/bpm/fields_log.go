package BaseBPM

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	//BPM ID
	BPMID int64 `db:"bpm_id" json:"bpmId" check:"id"`
	//当前节点ID
	NodeID string `db:"node_id" json:"nodeId" check:"des" min:"1" max:"300"`
	//当前节点序号
	NodeNum int `db:"node_num" json:"nodeNum" check:"int64Than0"`
	//节点存储内容
	NodeContent string `db:"node_content" json:"nodeContent"`
}
