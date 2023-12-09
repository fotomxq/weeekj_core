package BaseToken2

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//会话过期处理
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsExpire)
}

// 会话过期处理
func subNatsExpire(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	if action != "base_token2" {
		return
	}
	//获取ID
	data := getByID(id)
	if data.ID < 1 {
		return
	}
	//如果5秒内到期，则删除
	if data.ExpireAt.Unix() > CoreFilter.GetNowTime().Unix()-5 {
		return
	}
	//删除会话
	DeleteToken(id)
}
