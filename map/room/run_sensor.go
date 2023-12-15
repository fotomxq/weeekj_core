package MapRoom

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	IOTBind "github.com/fotomxq/weeekj_core/v5/iot/bind"
	IOTSensor "github.com/fotomxq/weeekj_core/v5/iot/sensor"
)

// 将统计数据整合迁移
// 只统计上一个小时的数据，如果存在则跳过设备
func runSensor() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("map room sensor run, ", r)
		}
	}()
	//批量遍历房间
	limit := 100
	step := 0
	for {
		roomList := getRunSensorRoomList(limit, step)
		if len(roomList) < 1 {
			break
		}
		for _, vRoom := range roomList {
			//获取该房间绑定的所有设备
			bindList, err := IOTBind.GetBindFrom(&IOTBind.ArgsGetBindFrom{
				OrgID: vRoom.OrgID,
				FromInfo: CoreSQLFrom.FieldsFrom{
					System: "room",
					ID:     vRoom.ID,
					Mark:   "*",
					Name:   "",
				},
			})
			if err != nil {
				continue
			}
			if len(bindList) < 1 {
				continue
			}
			//遍历设备，查询设备的遥感数据集合
			// 之前2小时到1小时的时间
			minTime := CoreFilter.GetNowTimeCarbon().SubMinutes(10)
			maxTime := CoreFilter.GetNowTimeCarbon()
			for _, vBind := range bindList {
				//获取该绑定设备的统计数据
				sensorList, err := IOTSensor.GetListTime(&IOTSensor.ArgsGetListTime{
					TimeBetween: CoreSQLTime.FieldsCoreTime{
						MinTime: minTime.Time,
						MaxTime: maxTime.Time,
					},
					DeviceID: vBind.DeviceID,
					Mark:     "",
				})
				if err != nil {
					continue
				}
				if len(sensorList) < 1 {
					continue
				}
				//计算平均值
				var sensorDataCount int64
				var sensorData int64
				var sensorFCount float64
				var sensorF float64
				var sensorS string
				for _, vSensor := range sensorList {
					sensorDataCount += vSensor.Data
					sensorFCount += vSensor.DataF
					sensorS = vSensor.DataS
				}
				if len(sensorList) > 0 {
					sensorData = sensorDataCount / int64(len(sensorList))
					sensorF = sensorFCount / float64(len(sensorList))
				}
				//检查是否存在数据？
				var data FieldsSensor
				if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM map_room_iot_sensor WHERE device_id = $1 AND org_id = $2 AND create_at <= $3 AND create_at >= $4", vBind.DeviceID, vBind.OrgID, minTime, maxTime); err == nil && data.ID > 0 {
					//存在则准备修改数据
					_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE map_room_iot_sensor SET data = :data, data_f = :data_f, data_s = :data_s WHERE id = :id", map[string]interface{}{
						"id":     data.ID,
						"data":   sensorData,
						"data_f": sensorF,
						"data_s": sensorS,
					})
					if err != nil {
						CoreLog.Error("map room sensor run, update data, ", err)
						continue
					}
				} else {
					_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO map_room_iot_sensor (create_at, device_id, org_id, room_id, mark, data, data_f, data_s) VALUES (:create_at,:device_id,:org_id,:room_id,:mark,:data,:data_f,:data_s)", map[string]interface{}{
						"create_at": maxTime.Time,
						"device_id": sensorList[0].DeviceID,
						"org_id":    vBind.OrgID,
						"room_id":   vRoom.ID,
						"mark":      sensorList[0].Mark,
						"data":      sensorData,
						"data_f":    sensorF,
						"data_s":    sensorS,
					})
					if err != nil {
						CoreLog.Error("map room sensor run, insert new data, ", err)
						continue
					}
				}
				//消峰处理
				time.Sleep(time.Millisecond * 200)
			}
			//消峰处理
			time.Sleep(time.Millisecond * 200)
		}
		step += limit
	}
	//删除1个月以上的数据
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "map_room_iot_sensor", "create_at < :create_at", map[string]interface{}{
		"create_at": CoreFilter.GetNowTimeCarbon().SubMonth().Time,
	})
}
