package IOTMission

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetMissionList 获取任务列表参数
type ArgsGetMissionList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备分组
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//状态
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//任务动作
	Action string `db:"action" json:"action" check:"mark" empty:"true"`
	//任务动作连接方案
	ConnectType string `db:"connect_type" json:"connectType" check:"mark" empty:"true"`
	//是否为历史
	IsHistory bool `db:"is_history" json:"isHistory" check:"bool"`
}

// GetMissionList 获取任务列表
func GetMissionList(args *ArgsGetMissionList) (dataList []FieldsMission, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.GroupID > 0 {
		where = where + "group_id = :group_id"
		maps["group_id"] = args.GroupID
	}
	if args.DeviceID > 0 {
		where = where + "device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	if args.Status > 0 {
		where = where + "status = :status"
		maps["status"] = args.Status
	}
	if args.Action != "" {
		where = where + "action = :action"
		maps["action"] = args.Action
	}
	if args.ConnectType != "" {
		where = where + "connect_type = :connect_type"
		maps["connect_type"] = args.ConnectType
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_mission"
	if args.IsHistory {
		tableName = "iot_mission_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, org_id, group_id, device_id, status, action, connect_type FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at"},
	)
	return
}

// ArgsGetWaitMissionByDevice 获取设备待处理任务参数
type ArgsGetWaitMissionByDevice struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
}

// GetWaitMissionByDevice 获取设备待处理任务
func GetWaitMissionByDevice(args *ArgsGetWaitMissionByDevice) (data FieldsMission, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, org_id, group_id, device_id, status, params_data, report_data, action, connect_type FROM iot_mission WHERE device_id = $1 ORDER BY id DESC LIMIT 1", args.DeviceID)
	return
}

// ArgsCreateMission 创建新的任务参数
type ArgsCreateMission struct {
	//组织ID
	// -1则忽略
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//发送请求数据集合
	ParamsData []byte `db:"params_data" json:"paramsData"`
	//任务动作
	Action string `db:"action" json:"action" check:"mark"`
}

// CreateMission 创建新的任务
func CreateMission(args *ArgsCreateMission) (data FieldsMission, errCode string, err error) {
	//检查设备是否支持动作
	var deviceData IOTDevice.FieldsDevice
	deviceData, err = IOTDevice.GetDeviceByID(&IOTDevice.ArgsGetDeviceByID{
		ID:    args.DeviceID,
		OrgID: args.OrgID,
	})
	if err != nil || deviceData.ID < 1 {
		errCode = "device_not_exist"
		err = errors.New(fmt.Sprint("device not exist, ", err))
		return
	}
	var groupData IOTDevice.FieldsGroup
	groupData, err = IOTDevice.GetGroupByID(&IOTDevice.ArgsGetGroupByID{
		ID: deviceData.GroupID,
	})
	if err != nil || groupData.ID < 1 {
		errCode = "group_not_exist"
		err = errors.New(fmt.Sprint("device group not exist, ", err))
		return
	}
	var actionList []IOTDevice.FieldsAction
	actionList, err = IOTDevice.GetActionMore(&IOTDevice.ArgsGetActionMore{
		IDs:        groupData.Action,
		HaveRemove: false,
	})
	var actionData IOTDevice.FieldsAction
	for _, v := range actionList {
		if v.Mark == args.Action {
			actionData = v
			break
		}
	}
	if actionData.ID < 1 {
		errCode = "action_not_support"
		err = errors.New("not support action, action: " + args.Action)
		return
	}
	//任务过期时间
	expireAt := CoreFilter.GetNowTime().Add(time.Second * time.Duration(actionData.ExpireTime))
	//检查组织是否可以执行该任务？
	if args.OrgID > 0 {
		var operateData IOTDevice.FieldsOperate
		operateData, err = IOTDevice.CheckOperate(&IOTDevice.ArgsCheckOperate{
			DeviceID: deviceData.ID,
			OrgID:    args.OrgID,
		})
		if err != nil || operateData.ID < 1 {
			errCode = "operate"
			err = errors.New(fmt.Sprint("check operate, ", err))
			return
		}
		allowMission := false
		for _, v := range operateData.Permissions {
			if v == "all" || v == "mission" {
				allowMission = true
				break
			}
		}
		if !allowMission {
			errCode = "not_support_mission"
			err = errors.New("org not support mission")
			return
		}
		isFind := false
		for _, v := range operateData.Action {
			if v == actionData.ID {
				isFind = true
				break
			}
		}
		if !isFind {
			errCode = "not_support_action"
			err = errors.New("org not support action")
			return
		}
	}
	//创建任务
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "iot_mission", "INSERT INTO iot_mission (expire_at, org_id, group_id, device_id, status, params_data, report_data, action, connect_type) VALUES (:expire_at, :org_id, :group_id, :device_id, 0, :params_data, :report_data, :action, :connect_type)", map[string]interface{}{
		"expire_at":    expireAt,
		"org_id":       args.OrgID,
		"group_id":     groupData.ID,
		"device_id":    args.DeviceID,
		"params_data":  args.ParamsData,
		"report_data":  []byte{},
		"action":       args.Action,
		"connect_type": actionData.ConnectType,
	}, &data)
	if err != nil {
		errCode = "insert"
	}
	//记录阻拦器
	RunMissionBlocker.NewEdit()
	//反馈
	return
}

// ArgsUpdateMissionStatus 更新任务状态完成参数
type ArgsUpdateMissionStatus struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//任务状态
	// 0 wait 等待发起 / 1 send 已经发送 / 3 failed 已经失败 / 4 cancel 取消
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
}

// UpdateMissionStatus 更新任务状态完成
func UpdateMissionStatus(args *ArgsUpdateMissionStatus) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_mission SET update_at = NOW(), status = :status WHERE id = :id", args)
	return
}

// ArgsUpdateMissionFinish 更新任务状态完成参数
type ArgsUpdateMissionFinish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//回收数据
	// 回收数据如果过大，将不会被存储到本地
	ReportData []byte `db:"report_data" json:"reportData"`
}

// UpdateMissionFinish 更新任务状态完成
func UpdateMissionFinish(args *ArgsUpdateMissionFinish) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE iot_mission SET update_at = NOW(), status = 2, report_data = :report_data WHERE id = :id", args)
	if err != nil {
		return
	}
	_ = CreateLog(&ArgsCreateLog{
		MissionID: args.ID,
		Status:    2,
		Mark:      "finish",
		Content:   "任务执行完成并回执",
	})
	return
}

// ArgsDeleteMission 删除任务参数
type ArgsDeleteMission struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteMission 删除任务
func DeleteMission(args *ArgsDeleteMission) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_mission", "id = :id AND (status = 0 OR status = 1)", args)
	if err != nil {
		return
	}
	_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_mission_log", "mission_id = :mission_id", map[string]interface{}{
		"mission_id": args.ID,
	})
	return
}
