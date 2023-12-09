package BaseRank

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 过期处理
func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("rank expire run, ", r)
		}
	}()
	//获取过期数据，转移
	limit := 100
	step := 0
	for {
		var dataList []FieldsRank
		err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, expire_at, pick_min, pick_at, service_mark, mission_mark, mission_data FROM core_rank WHERE expire_at <= NOW() LIMIT $1 OFFSET $2", limit, step)
		if err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			err = overRankByData(&v, []byte{})
			if err != nil {
				CoreLog.Error("rank expire run, delete and insert data, ", err)
			}
		}
		step += limit
	}
}
