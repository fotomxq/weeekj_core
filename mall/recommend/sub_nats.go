package MallRecommend

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//为用户组装数据
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "商城推荐服务用户",
		Description:  "",
		EventSubType: "all",
		Code:         "mall_recommend_user",
		EventType:    "nats",
		EventURL:     "/mall/recommend/user",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("mall_recommend_user", "/mall/recommend/user", subNatsUser)
}

// 为用户组装数据
func subNatsUser(_ *nats.Msg, _ string, userID int64, _ string, _ []byte) {
	//相同用户禁止频繁提交数据
	blockerUser.CheckWait(userID, "", func(modID int64, _ string) {
		putUserData(modID)
	})
}
