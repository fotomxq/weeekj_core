package BaseExpireTip

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//通知写入了新的数据
	CoreNats.SubDataByteNoErr("/base/expire_tip/new", func(msg *nats.Msg, action string, id int64, mark string, data []byte) {
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
