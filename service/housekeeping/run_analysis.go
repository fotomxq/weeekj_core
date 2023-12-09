package ServiceHousekeeping

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runAnalysis() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("service housekeeping analysis run, ", r)
		}
	}()
	type vAnalysisType struct {
		//总数
		Count int64 `db:"count" json:"count"`
		//公里数
		KM int64 `db:"km" json:"km"`
		//总耗时
		OverTime int64 `db:"over_time" json:"overTime"`
		//评级
		// 1-5 级别
		Level int `db:"level" json:"level"`
	}
	type vCountType struct {
		//总数
		Count int64 `db:"count" json:"count"`
	}
	//遍历所有绑定成员
	limit := 100
	step := 0
	for {
		var bindList []FieldsBind
		if err := Router2SystemConfig.MainDB.Select(&bindList, "SELECT id, org_id, bind_id FROM service_housekeeping_bind WHERE delete_at < to_timestamp(1000000) ORDER BY id LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(bindList) < 1 {
			break
		}
		for _, vBind := range bindList {
			//未完成任务量
			vUnFinishCount, err := CoreSQL.GetAllCountMap(Router2SystemConfig.MainDB.DB, "service_housekeeping_log", "id", "org_id = :org_id AND bind_id = :bind_id AND delete_at < to_timestamp(1000000) AND finish_at < to_timestamp(1000000)", map[string]interface{}{
				"org_id":  vBind.OrgID,
				"bind_id": vBind.BindID,
			})
			if err != nil {
				//
			}
			//更新数量
			if _, err := CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_housekeeping_bind SET un_finish_count = :un_finish_count WHERE id = :id", map[string]interface{}{
				"id":              vBind.ID,
				"un_finish_count": vUnFinishCount,
			}); err != nil {
				CoreLog.Error("service housekeeping analysis run, update bind analysis, ", err)
				continue
			}
		}
		step += limit
	}
}
