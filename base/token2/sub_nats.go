package BaseToken2

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//会话过期处理
	CoreNats.SubDataByteNoErr("base_expire_tip_expire", "/base/expire_tip/expire", subNatsExpire)
	//清理所有过期数据
	CoreNats.SubDataByteNoErr("base_expire_tip_expire_clear", "/base/expire_tip/expire_clear", subNatsClearExpire)
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

// 清理全部过期会话
func subNatsClearExpire(_ *nats.Msg, _ string, _ int64, _ string, _ []byte) {
	clearExpireToken()
}
