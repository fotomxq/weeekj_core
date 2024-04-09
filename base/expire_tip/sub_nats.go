package BaseExpireTip

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//通知写入了新的数据
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "过期模块通知",
		Description:  "数据需要发生过期，等待处理",
		EventSubType: "sub",
		Code:         "analysis_expire_tip_new",
		EventType:    "nats",
		EventURL:     "/base/expire_tip/new",
		//TODO:待补充
		EventParams: "<<id>>:int64:过期服务ID;<<data>>:string:'待补充'",
	})
	CoreNats.SubDataByteNoErr("analysis_expire_tip_new", "/base/expire_tip/new", func(_ *nats.Msg, _ string, id int64, _ string, data []byte) {
		hash := gjson.GetBytes(data, "hash").String()
		expireAtStr := gjson.GetBytes(data, "expireAt").String()
		expireAt, err := CoreFilter.GetTimeByISO(expireAtStr)
		if err != nil {
			//异常数据，直接从数据库获取一次
			rawData, err := getID(id)
			if err != nil {
				return
			}
			hash = rawData.Hash
			expireAt = rawData.ExpireAt
		}
		appendHaveNewData(id, hash, expireAt)
	})
}
