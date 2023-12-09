package AnalysisAny2

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// AppendData 添加统计数据
// 只保留相同参数的最新值
// action 支持：re 替代数据；add 叠加数据; reduce 减少
func AppendData(action, mark string, createAt time.Time, orgID, userID, bindID, param1, param2 int64, data int64) {
	if action == "" {
		action = "re"
	}
	appendData(waitAppendDataType{
		Action:   action,
		Mark:     mark,
		CreateAt: createAt,
		OrgID:    orgID,
		UserID:   userID,
		BindID:   bindID,
		Param1:   param1,
		Param2:   param2,
		Data:     data,
	})
}

// 写入数据
func appendData(data waitAppendDataType) {
	appendLog := "analysis any2 append data, "
	//获取配置
	configData, err := getConfigByMark(data.Mark, false)
	if err != nil {
		CoreLog.Warn(appendLog, "get config mark: ", data.Mark, ", ", err)
		return
	}
	//根据配置计算要查询的时间范围
	nowAt := CoreFilter.GetNowTimeCarbon()
	if data.CreateAt.Unix() > 1000000 {
		nowAt = nowAt.CreateFromGoTime(data.CreateAt)
	}
	var findAtMin time.Time
	var findAtMax time.Time
	switch configData.Particle {
	case 0:
		findAtMin = nowAt.StartOfHour().Time
		findAtMax = nowAt.EndOfHour().Time
	case 1:
		findAtMin = nowAt.StartOfDay().Time
		findAtMax = nowAt.EndOfDay().Time
	case 2:
		findAtMin = nowAt.StartOfWeek().Time
		findAtMax = nowAt.EndOfWeek().Time
	case 3:
		findAtMin = nowAt.StartOfMonth().Time
		findAtMax = nowAt.EndOfMonth().Time
	case 4:
		findAtMin = nowAt.StartOfYear().Time
		findAtMax = nowAt.EndOfYear().Time
	default:
		findAtMin = nowAt.StartOfHour().Time
		findAtMax = nowAt.EndOfHour().Time
	}
	//查询最近加入的数据ID
	var id int64
	_ = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM analysis_any2 WHERE config_id = $1 AND org_id = $2 AND user_id = $3 AND bind_id = $4 AND params1 = $5 AND params2 = $6 AND create_at >= $7 AND create_at <= $8 ORDER BY id DESC LIMIT 1", configData.ID, data.OrgID, data.UserID, data.BindID, data.Param1, data.Param2, findAtMin, findAtMax)
	//更新数据
	if id > 0 {
		//根据模式类型写入数据
		switch data.Action {
		case "re":
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any2 SET data = :data WHERE id = :id", map[string]interface{}{
				"id":   id,
				"data": data.Data,
			})
		case "add":
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any2 SET data = data + :data WHERE id = :id", map[string]interface{}{
				"id":   id,
				"data": data.Data,
			})
		case "reduce":
			_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any2 SET data = data - :data WHERE id = :id", map[string]interface{}{
				"id":   id,
				"data": data.Data,
			})
		}
		if err != nil {
			CoreLog.Warn(appendLog, "update data mark: ", data.Mark, ", data: ", data.Data, ", err: ", err)
			return
		}
	} else {
		//写入数据
		switch data.Action {
		case "re":
			_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_any2 (create_at, org_id, user_id, bind_id, params1, params2, config_id, data) VALUES (:create_at,:org_id,:user_id,:bind_id,:params1,:params2,:config_id,:data)", map[string]interface{}{
				"create_at": nowAt.Time,
				"org_id":    data.OrgID,
				"user_id":   data.UserID,
				"bind_id":   data.BindID,
				"params1":   data.Param1,
				"params2":   data.Param2,
				"config_id": configData.ID,
				"data":      data.Data,
			})
		case "add":
			_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_any2 (create_at, org_id, user_id, bind_id, params1, params2, config_id, data) VALUES (:create_at,:org_id,:user_id,:bind_id,:params1,:params2,:config_id,:data)", map[string]interface{}{
				"create_at": nowAt.Time,
				"org_id":    data.OrgID,
				"user_id":   data.UserID,
				"bind_id":   data.BindID,
				"params1":   data.Param1,
				"params2":   data.Param2,
				"config_id": configData.ID,
				"data":      data.Data,
			})
		case "reduce":
			_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_any2 (create_at, org_id, user_id, bind_id, params1, params2, config_id, data) VALUES (:create_at,:org_id,:user_id,:bind_id,:params1,:params2,:config_id,:data)", map[string]interface{}{
				"create_at": nowAt.Time,
				"org_id":    data.OrgID,
				"user_id":   data.UserID,
				"bind_id":   data.BindID,
				"params1":   data.Param1,
				"params2":   data.Param2,
				"config_id": configData.ID,
				"data":      0 - data.Data,
			})
		}
		if err != nil {
			CoreLog.Warn(appendLog, "append data mark: ", data.Mark, ", data: ", data.Data, ", err: ", err)
			return
		}
	}
	//清理缓冲
	clearAnyCache(configData.ID)
}
