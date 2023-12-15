package ServiceUserInfoMod

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"time"
)

// PushInfoOut 推送标记老人离开
func PushInfoOut(infoID int64, outAt time.Time) {
	CoreNats.PushDataNoErr("/service/user_info/post_update", "out", infoID, "", map[string]interface{}{
		"atTime": CoreFilter.GetTimeToDefaultTime(outAt),
	})
}

// PushInfoReturn 还原档案
func PushInfoReturn(infoID int64) {
	CoreNats.PushDataNoErr("/service/user_info/post_update", "return", infoID, "", nil)
}

// ArgsAppendLog 添加日志参数
type ArgsAppendLog struct {
	//档案ID
	InfoID int64 `db:"info_id" json:"infoID"`
	//组织ID
	// 允许平台方的0数据，该数据可能来源于其他领域
	OrgID int64 `db:"org_id" json:"orgID"`
	//修改的位置
	// 1. 字段
	// 2. 或扩展参数指定的内容，例如params.[mark]
	// 3. 其他内容采用.形式跨越记录
	// 4. room.in 入驻房间变更
	ChangeMark string `db:"change_mark" json:"changeMark"`
	ChangeDes  string `db:"change_des" json:"changeDes"`
	//修改前描述
	OldDes string `db:"old_des" json:"oldDes"`
	//修改后描述
	NewDes string `db:"new_des" json:"newDes"`
}

// AppendLog 添加日志
func AppendLog(args ArgsAppendLog) {
	CoreNats.PushDataNoErr("/service/user_info/append_log", "", 0, "", args)
}
