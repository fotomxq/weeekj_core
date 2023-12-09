package TMSUserRunning

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteMission 删除任务参数
type ArgsDeleteMission struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// DeleteMission 删除任务
func DeleteMission(args *ArgsDeleteMission) (err error) {
	var newLog string
	newLog, err = getLogData("取消跑腿任务", []int64{})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE tms_user_running_mission SET delete_at = NOW(), logs = logs || :logs WHERE id = :id AND (:user_id < 1 OR user_id = :user_id)", map[string]interface{}{
		"id":      args.ID,
		"user_id": args.UserID,
		"logs":    newLog,
	})
	if err != nil {
		return
	}
	//通知成功
	pushNatsStatusUpdate("cancel", args.ID, "跑腿员取消服务")
	//反馈
	return
}
