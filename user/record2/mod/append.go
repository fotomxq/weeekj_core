package UserRecord2Mod

import (
	"fmt"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
)

// argsAppendData 插入数据参数
type argsAppendData struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//系统来源
	System string `db:"system" json:"system"`
	//影响ID
	ModID int64 `db:"mod_id" json:"modID"`
	//操作内容标识码
	Mark string `db:"mark" json:"mark"`
	//操作内容概述
	Des string `db:"des" json:"des"`
}

// AppendData 插入数据
func AppendData(orgID int64, orgBindID int64, userID int64, system string, modID int64, mark string, des ...interface{}) {
	//重组消息
	var desStr string
	for _, v := range des {
		desStr = desStr + fmt.Sprint(v)
	}
	//通知写入数据
	CoreNats.PushDataNoErr("user_record2_append", "/user/record2/append", "", 0, "", argsAppendData{
		OrgID:     orgID,
		OrgBindID: orgBindID,
		UserID:    userID,
		System:    system,
		ModID:     modID,
		Mark:      mark,
		Des:       desStr,
	})
}
