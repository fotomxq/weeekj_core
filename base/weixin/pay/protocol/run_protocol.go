package BaseWeixinPayProtocol

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func runNext() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("base weixin pay protocol run, ", r)
		}
	}()
	//获取合同列表
	limit := 100
	step := 0
	for {
		var dataList []FieldsProtocol
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM core_weixin_pay_protocol LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			//检查配置对应的数据是否到期？如果还有24小时到期，则触发订单请求，否则跳过处理
			// 触发请求后，如果扣款请求失败，则禁止继续请求
			switch v.ConfigSystem {
			}
		}
		step += limit
	}
}
