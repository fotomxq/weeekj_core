package IOTMission

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// PickMissionList 提取一组等待发送的任务ID
func PickMissionList() (dataList []FieldsMission, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, connect_type FROM iot_mission WHERE status = 0 ORDER BY id LIMIT 100")
	return
}

// ArgsPickMission 提取一个任务参数
type ArgsPickMission struct {
	//任务ID
	ID int64 `json:"id"`
}

// PickMission 提取一个任务
func PickMission(args *ArgsPickMission) (data FieldsMission, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, org_id, group_id, device_id, status, params_data, action, connect_type FROM iot_mission WHERE id = $1 AND expire_at >= NOW()", args.ID)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_mission SET status = 1 WHERE id = :id", args)
	if err != nil {
		return
	}
	return
}
