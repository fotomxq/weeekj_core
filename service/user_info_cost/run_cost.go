package ServiceUserInfoCost

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	IOTBind "github.com/fotomxq/weeekj_core/v5/iot/bind"
	IOTSensor "github.com/fotomxq/weeekj_core/v5/iot/sensor"
	MapRoom "github.com/fotomxq/weeekj_core/v5/map/room"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

func runCost() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("service user info cost run, ", r)
		}
	}()
	limit := 100
	step := 0
	for {
		//获取配置列表
		var dataList []FieldsConfig
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, org_id, name, room_bind_mark, sensor_mark, count_type, each_unit, each_price, params FROM service_user_info_cost_config WHERE delete_at < to_timestamp(1000000) LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, vConfig := range dataList {
			runCostConfig(&vConfig)
		}
		//下一页
		step += limit
	}
}

// 统计上一个小时的数据，如果存在则跳过
func runCostConfig(configData *FieldsConfig) {
	//获取上一个小时的开始结束时间
	startAt := CoreFilter.GetNowTimeCarbon().SubHour().StartOfHour()
	endAt := CoreFilter.GetNowTimeCarbon().SubHour().EndOfHour()
	//根据配置的room_bind_mark获取绑定关系结构
	var page int64 = 1
	for {
		bindList, _, err := IOTBind.GetBindList(&IOTBind.ArgsGetBindList{
			Pages: CoreSQLPages.ArgsDataList{
				Page: page,
				Max:  100,
				Sort: "id",
				Desc: false,
			},
			OrgID:    configData.OrgID,
			DeviceID: -1,
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "room",
				ID:     -1,
				Mark:   configData.RoomBindMark,
				Name:   "",
			},
			IsRemove: false,
		})
		if err != nil {
			break
		}
		if len(bindList) < 1 {
			break
		}
		for _, vBind := range bindList {
			//检查上一个小时的数据是否存在，如果存在则跳过
			var prevData FieldsCost
			err = Router2SystemConfig.MainDB.Get(&prevData, "SELECT id FROM service_user_info_cost WHERE org_id = $1 AND config_id = $2 AND room_id = $3 AND sensor_mark = $4 AND create_at >= $5 AND create_at < $6", configData.OrgID, configData.ID, vBind.FromInfo.ID, configData.SensorMark, startAt, endAt)
			if err == nil && prevData.ID > 0 {
				continue
			}
			//找到房间居住人
			var infoID int64
			roomData, err := MapRoom.GetRoomID(&MapRoom.ArgsGetRoomID{
				ID:    vBind.FromInfo.ID,
				OrgID: vBind.OrgID,
			})
			if err == nil && roomData.ID > 0 {
				if len(roomData.Infos) == 1 {
					infoID = roomData.Infos[0]
				}
			}
			//甄别统计形式
			// 根据绑定关系，查询该设备的传感器值
			switch configData.CountType {
			case 0:
				//0 合并计算，将时间阶段内的所有遥感数据合并进行统计计算
				sensorList, err := IOTSensor.GetAnalysis(&IOTSensor.ArgsGetAnalysis{
					TimeBetween: CoreSQLTime.FieldsCoreTime{
						MinTime: startAt.Time,
						MaxTime: endAt.Time,
					},
					TimeType:  "hour",
					DeviceID:  vBind.DeviceID,
					Mark:      configData.SensorMark,
					IsHistory: false,
				})
				if err != nil {
					continue
				}
				if len(sensorList) < 1 {
					continue
				}
				var sumData float64
				var sumType int
				for _, v := range sensorList {
					if sumType < 1 {
						if v.Data > 0 {
							sumType = 1
						} else {
							if v.DataF > 0 {
								sumType = 2
							}
						}
					} else {
						switch sumType {
						case 1:
							sumData = sumData + float64(v.Data)
						case 2:
							sumData = sumData + v.DataF
						}
					}
				}
				_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_user_info_cost (create_at, org_id, room_id, info_id, config_id, room_bind_mark, sensor_mark, unit, currency, price) VALUES (:create_at,:org_id,:room_id,:info_id,:config_id,:room_bind_mark,:sensor_mark,:unit,:currency,:price)", map[string]interface{}{
					"create_at":      startAt.Time,
					"org_id":         configData.OrgID,
					"room_id":        vBind.FromInfo.ID,
					"info_id":        infoID,
					"room_bind_mark": configData.RoomBindMark,
					"sensor_mark":    configData.SensorMark,
					"unit":           sumData,
					"currency":       configData.Currency,
					"price":          int64((float64(configData.EachPrice) / configData.EachUnit) * sumData),
				})
				if err != nil {
					CoreLog.Error("service user info cost run, insert sum data, ", err)
					continue
				}
			case 1:
				//1 平均值计算，将时间段内的数据平均化计算
				sensorList, err := IOTSensor.GetAnalysisAvg(&IOTSensor.ArgsGetAnalysis{
					TimeBetween: CoreSQLTime.FieldsCoreTime{
						MinTime: startAt.Time,
						MaxTime: endAt.Time,
					},
					TimeType:  "hour",
					DeviceID:  vBind.DeviceID,
					Mark:      configData.SensorMark,
					IsHistory: false,
				})
				if err != nil {
					continue
				}
				if len(sensorList) < 1 {
					continue
				}
				var sumData float64
				var sumType int
				for _, v := range sensorList {
					if sumType < 1 {
						if v.Data > 0 {
							sumType = 1
						} else {
							if v.DataF > 0 {
								sumType = 2
							}
						}
					} else {
						switch sumType {
						case 1:
							sumData = sumData + float64(v.Data)
						case 2:
							sumData = sumData + v.DataF
						}
					}
				}
				sumData = sumData / float64(len(sensorList))
				_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_user_info_cost (create_at, org_id, room_id, info_id, config_id, room_bind_mark, sensor_mark, unit, currency, price) VALUES (:create_at,:org_id,:room_id,:info_id,:config_id,:room_bind_mark,:sensor_mark,:unit,:currency,:price)", map[string]interface{}{
					"create_at":      startAt.Time,
					"org_id":         configData.OrgID,
					"room_id":        vBind.FromInfo.ID,
					"info_id":        infoID,
					"room_bind_mark": configData.RoomBindMark,
					"sensor_mark":    configData.SensorMark,
					"unit":           sumData,
					"currency":       configData.Currency,
					"price":          int64((float64(configData.EachPrice) / configData.EachUnit) * sumData),
				})
				if err != nil {
					CoreLog.Error("service user info cost run, insert avg data, ", err)
					continue
				}
			}
		}
		page += 1
	}
}
