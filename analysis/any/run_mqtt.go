package AnalysisAny

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// 自动推送MQTT数据包
// Deprecated 废弃
func runMQTT() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("analysis any run, auto push mqtt, ", r)
		}
	}()
	//批量获取配置
	limit := 5
	step := 0
	for {
		var configList []FieldsConfig
		err := Router2SystemConfig.MainDB.Select(&configList, "SELECT id, mqtt_org, mqtt_user, mqtt_user, last_mqtt, last_hash FROM analysis_any_config WHERE delete_at < to_timestamp(1000000) AND last_mqtt < $1 LIMIT $2 OFFSET $3", CoreFilter.GetNowTimeCarbon().SubSeconds(10).Time, limit, step)
		if err != nil || len(configList) < 1 {
			break
		}
		for _, vConfig := range configList {
			//获取数据
			var data FieldsAny
			err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, user_id, bind_id, config_id, hash, data, data_val FROM analysis_any WHERE config_id = $1 ORDER BY id DESC LIMIT 1", vConfig.ID)
			if err != nil || data.ID < 1 {
				continue
			}
			if data.Hash == vConfig.LastHash {
				continue
			}
			if data.OrgID > 0 {
				//_ = IOTMQTT.PushUpdateData(data.OrgID, "analysis_any", 0)
			}
			//更新推送时间
			_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any_config SET last_mqtt = NOW(), last_hash = :last_hash WHERE id = :id", map[string]interface{}{
				"id":        vConfig.ID,
				"last_hash": data.Hash,
			})
		}
		step += limit
	}
}
