package ToolsCommunication

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("tools communication expire run, ", r)
		}
	}()
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "tools_communication_from", "expire_at < NOW()", nil)
	limit := 100
	step := 0
	for {
		var dataList []FieldsRoom
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM tools_communication_room WHERE delete_at < to_timestamp(1000000) AND expire_at < NOW() LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "tools_communication_from", "room_id = :room_id", map[string]interface{}{
				"room_id": v.ID,
			})
			_, _ = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "tools_communication_room", "id = :id", map[string]interface{}{
				"id": v.ID,
			})
		}
		step += limit
	}
}
