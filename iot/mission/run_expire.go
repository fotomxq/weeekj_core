package IOTMission

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runExpire() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("iot device mission run error, ", r)
		}
	}()
	//过期处理
	_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_mission SET update_at = NOW(), status = 3 WHERE (status = 0 OR status = 1) AND expire_at < NOW()", nil)
	limit := 100
	step := 0
	for {
		var dataList []FieldsMission
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM iot_mission WHERE (status = 0 OR status = 1) AND expire_at < NOW() LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			tx := Router2SystemConfig.MainDB.MustBegin()
			result, err := tx.Exec("INSERT INTO iot_mission_history(create_at, update_at, expire_at, org_id, group_id, device_id, status, params_data, report_data, action, connect_type) SELECT create_at, update_at, expire_at, org_id, group_id, device_id, status, params_data, report_data, action, connect_type FROM iot_mission WHERE id = $1", v.ID)
			if err = CoreSQL.LastRowsAffected(tx, result, err); err != nil {
				CoreLog.Error("iot device mission run, insert history, ", err)
				continue
			}
			result, err = tx.Exec("DELETE FROM iot_mission WHERE id = $1", v.ID)
			if err = CoreSQL.LastRowsAffected(tx, result, err); err != nil {
				CoreLog.Error("iot device mission run, delete old data, ", err)
				continue
			}
			result, err = tx.Exec("INSERT INTO iot_mission_log_history(create_at, mission_id, status, mark, content) SELECT create_at, mission_id, status, mark, content FROM iot_mission_log WHERE mission_id = $1", v.ID)
			if err = CoreSQL.LastRowsAffected(tx, result, err); err != nil {
				CoreLog.Error("iot device mission run, insert log history, ", err)
				continue
			}
			result, err = tx.Exec("DELETE FROM iot_mission_log WHERE mission_id = $1", v.ID)
			if err = CoreSQL.LastRowsAffected(tx, result, err); err != nil {
				CoreLog.Error("iot device mission run, delete old log data, ", err)
				continue
			}
			if err = tx.Commit(); err != nil {
				CoreLog.Error("iot device mission run, insert history, ", err)
				continue
			}
		}
		step += limit
	}
}
