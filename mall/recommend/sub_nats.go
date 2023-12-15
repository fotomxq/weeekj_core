package MallRecommend

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//为用户组装数据
	CoreNats.SubDataByteNoErr("/mall/recommend/user", subNatsUser)
}

// 为用户组装数据
func subNatsUser(_ *nats.Msg, _ string, userID int64, _ string, _ []byte) {
	//相同用户禁止频繁提交数据
	blockerUser.CheckWait(userID, "", func(modID int64, _ string) {
		putUserData(modID)
	})
}
