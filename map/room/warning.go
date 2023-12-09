package MapRoom

import (
	"errors"
	"fmt"
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetWarningList 获取应急呼叫日志列表参数
type ArgsGetWarningList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//是否需要完成
	NeedIsFinish bool `db:"need_is_finish" json:"needIsFinish" check:"bool"`
	IsFinish     bool `db:"is_finish" json:"isFinish" check:"bool"`
}

// GetWarningList 获取应急呼叫日志列表
func GetWarningList(args *ArgsGetWarningList) (dataList []FieldsWarning, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.RoomID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "room_id = :room_id"
		maps["room_id"] = args.RoomID
	}
	if args.DeviceID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "device_id = :device_id"
		maps["device_id"] = args.DeviceID
	}
	if args.NeedIsFinish {
		if where != "" {
			where = where + " AND "
		}
		if args.IsFinish {
			where = where + "finish_at >= to_timestamp(1000000)"
		} else {
			where = where + "finish_at < to_timestamp(1000000)"
		}
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"map_room_warning",
		"id",
		fmt.Sprint("SELECT id, create_at, finish_at, org_id, room_id, device_id, call_type FROM map_room_warning WHERE ", where),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "finish_at"},
	)
	return
}

// ArgsGetWarningByRooms 获取一批房间的预警情况参数
type ArgsGetWarningByRooms struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//房间ID
	RoomIDs pq.Int64Array `db:"room_ids" json:"roomIDs" check:"ids"`
	//呼叫类型
	// 0 紧急呼叫; 1 普通呼叫
	CallType int `db:"call_type" json:"callType"`
}

// GetWarningByRooms 获取一批房间的预警情况
// 只会反馈存在未完成应急呼叫的数据包
func GetWarningByRooms(args *ArgsGetWarningByRooms) (dataList []FieldsWarning, err error) {
	for _, v := range args.RoomIDs {
		appendData := getWarningByRoom(v, args.CallType)
		if appendData.ID < 1 || !CoreFilter.EqID2(args.OrgID, appendData.OrgID) || CoreSQL.CheckTimeHaveData(appendData.FinishAt) {
			continue
		}
		dataList = append(dataList, appendData)
	}
	if len(dataList) < 1 {
		err = errors.New("data is empty")
		return
	}
	return
}

// CheckWarningByRoomID 判断房间是否存在紧急呼叫？
func CheckWarningByRoomID(roomID int64, callType int) (b bool) {
	appendData := getWarningByRoom(roomID, callType)
	if appendData.ID < 1 || CoreSQL.CheckTimeHaveData(appendData.FinishAt) {
		return
	}
	return true
}

// ArgsAppendWarning 推送一个新的应急呼叫日志参数
type ArgsAppendWarning struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id" empty:"true"`
	//是否需要推送MQTT
	NeedMQTT bool `json:"needMQTT" check:"bool"`
	//呼叫类型
	// 0 紧急呼叫; 1 普通呼叫
	CallType int `db:"call_type" json:"callType"`
}

// AppendWarning 推送一个新的应急呼叫日志
// 1分钟内重复请求创建不会被记录到系统，不会抛出错误信息
func AppendWarning(args *ArgsAppendWarning) (haveData bool, err error) {
	//如果1小时内有数据，则跳出
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM map_room_warning WHERE org_id = $1 AND room_id = $2 AND device_id = $3 AND create_at >= $4 AND finish_at < to_timestamp(1000000) AND call_type = $5 ORDER BY id DESC LIMIT 1", args.OrgID, args.RoomID, args.DeviceID, CoreFilter.GetNowTimeCarbon().SubMinutes(1).Time, args.CallType)
	if err == nil && id > 0 {
		haveData = true
		return
	}
	//写入新数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO map_room_warning (org_id, room_id, device_id, call_type) VALUES (:org_id,:room_id,:device_id,:call_type)", map[string]interface{}{
		"org_id":    args.OrgID,
		"room_id":   args.RoomID,
		"device_id": args.DeviceID,
		"call_type": args.CallType,
	})
	if err != nil {
		return
	}
	//根据呼叫类型处理
	switch args.CallType {
	case 0:
		//清理缓冲
		deleteWarningCache(args.RoomID)
		//推送mqtt
		if args.NeedMQTT {
			pushRoomEmergencyCall(args.OrgID, args.DeviceID, args.RoomID)
		}
		//推送nats
		pushNatsUpdateStatus(args.RoomID, "ew", "on")
		//记录统计
		var count int64
		_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM map_room_log WHERE status = 7 AND org_id = $1 AND room_id = $2", args.OrgID, args.RoomID)
		AnalysisAny2.AppendData("re", "map_room_service_warning_btn_count", time.Time{}, args.OrgID, 0, args.RoomID, 0, 0, count)
		//获取房间入住人员信息
		infoIDs := getAppendWarningRoom(args.RoomID)
		if len(infoIDs) > 0 {
			for _, v := range infoIDs {
				AnalysisAny2.AppendData("re", "map_room_service_info_warning_btn_count", time.Time{}, args.OrgID, 0, v, 0, 0, count)
			}
		}
	case 1:
	}
	//反馈
	return
}

// UnWarning 解除房间的紧急呼叫
func UnWarning(args *ArgsAppendWarning) (err error) {
	//获取最近的紧急呼叫作为解除项
	data := getWarningByRoom(args.RoomID, args.CallType)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || CoreSQL.CheckTimeHaveData(data.FinishAt) {
		return
	}
	//更新紧急呼叫
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE map_room_warning SET finish_at = NOW() WHERE id = :id", map[string]interface{}{
		"id": data.ID,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteWarningCache(data.RoomID)
	//推送nats
	pushNatsUpdateStatus(data.RoomID, "ew", "off")
	//反馈
	return
}

// 根据房间获取最新的应急呼叫设置
func getWarningByRoom(roomID int64, callType int) (data FieldsWarning) {
	cacheMark := getWarningCacheMark(roomID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, finish_at, org_id, room_id, device_id, call_type FROM map_room_warning WHERE room_id = $1 AND call_type = $2 ORDER BY id DESC LIMIT 1", roomID, callType)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}
