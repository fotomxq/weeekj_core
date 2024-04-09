package UserLogin2

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

/**
TODO：替代当前V4的所有用户包的设计，改为此设计统一反馈
用户登录和汇总数据包
1. 统一获取用户登录后的数据包
2. 配合第二代token设计，自动填充和获取数据集
3. 解决用户锁定商户的问题，不需要再单独的表投放数据(token v2管理数据)
4. 二维码扫码登录的设计
*/

var (
	OpenSub = false
)

func Init() {
	if OpenSub {
		_ = BaseService.SetService(&BaseService.ArgsSetService{
			ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
			Name:         "用户登录2新增通知",
			Description:  "",
			EventSubType: "all",
			Code:         "user_login2_new",
			EventType:    "nats",
			EventURL:     "/user/login2/new",
			//TODO:待补充
			EventParams: "",
		})
	}
}
