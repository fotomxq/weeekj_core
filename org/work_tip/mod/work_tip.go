package OrgWorkTipMod

import CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"

// ArgsAppendTip 添加一个数据
type ArgsAppendTip struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//消息内容
	Msg string `db:"msg" json:"msg"`
	//系统
	System string `db:"system" json:"system"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID"`
}

func AppendTip(args *ArgsAppendTip) {
	CoreNats.PushDataNoErr("/org/work_tip", "new", 0, "", args)
}
