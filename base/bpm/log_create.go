package BaseBPM

type ArgsCreateLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//管理单元
	UnitID int64 `db:"unit_id" json:"unitID"`
	//操作用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//操作组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//BPM ID
	BPMID int64 `db:"bpm_id" json:"bpmId" check:"id"`
	//当前节点ID
	NodeID string `db:"node_id" json:"nodeId" check:"des" min:"1" max:"300"`
	//当前节点序号
	NodeNumber int `db:"node_number" json:"nodeNumber" check:"int64Than0"`
	//节点存储内容
	NodeContent string `db:"node_content" json:"nodeContent"`
}

func CreateLog(args *ArgsCreateLog) (id int64, err error) {
	//创建数据
	id, err = logDB.Insert().SetFields([]string{"org_id", "unit_id", "user_id", "org_bind_id", "bpm_id", "node_id", "node_number", "node_content"}).Add(map[string]any{
		"org_id":       args.OrgID,
		"unit_id":      args.UnitID,
		"user_id":      args.UserID,
		"org_bind_id":  args.OrgBindID,
		"bpm_id":       args.BPMID,
		"node_id":      args.NodeID,
		"node_number":  args.NodeNumber,
		"node_content": args.NodeContent,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}
