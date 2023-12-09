package MapRoom

import (
	"errors"
	"fmt"
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	ClassSort "gitee.com/weeekj/weeekj_core/v5/class/sort"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	OrgMission "gitee.com/weeekj/weeekj_core/v5/org/mission"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceUserInfoMod "gitee.com/weeekj/weeekj_core/v5/service/user_info/mod"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateRoom 更新房间基本信息参数
type ArgsUpdateRoom struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//房间编号
	Code string `db:"code" json:"code" check:"mark" empty:"true"`
	//入驻人员列
	Infos pq.Int64Array `db:"infos" json:"infos" check:"ids" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述图
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateRoom 更新房间基本信息
func UpdateRoom(args *ArgsUpdateRoom) (errCode string, err error) {
	//旧的房间数据
	oldData := getRoomID(args.ID)
	if oldData.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//检查入住人员是否在其他房间也入住
	var count int64
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND infos && $2 AND id != $3", args.OrgID, args.Infos, args.ID)
	if err == nil && count > 0 {
		errCode = "info_other_room"
		err = errors.New("info have other room")
		return
	}
	//检查房间的code是否唯一
	if args.Code != "" {
		var data FieldsRoom
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM map_room WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND code = $2", args.OrgID, args.Code)
		if err == nil && data.ID != args.ID {
			errCode = "code"
			err = errors.New(fmt.Sprint("mark is exist, ", err))
			return
		}
	}
	//更新房间信息
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_room SET update_at = NOW(), sort_id = :sort_id, tags = :tags, infos = :infos, code = :code, name = :name, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, params = :params WHERE id = :id", args)
	if err != nil {
		errCode = "err_insert"
		return
	}
	var data FieldsRoom
	data = getRoomID(args.ID)
	if data.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	//删除缓冲
	deleteRoomCache(data.ID)
	//推送nats
	pushNatsUpdateStatus(data.ID, "update_info", "")
	//更新统计
	pushNatsUpdateAnalysis(data.OrgID)
	//如果入驻人员不一致，则标记进入和退出
	for _, v := range oldData.Infos {
		isFind := false
		for _, v2 := range data.Infos {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			ServiceUserInfoMod.AppendLog(ServiceUserInfoMod.ArgsAppendLog{
				InfoID:     v,
				OrgID:      data.OrgID,
				ChangeMark: "room.out",
				ChangeDes:  "离开房间",
				OldDes:     "",
				NewDes:     fmt.Sprint("离开", oldData.Name, "[", oldData.ID, "]房间"),
			})
		}
	}
	for _, v := range data.Infos {
		isFind := false
		for _, v2 := range oldData.Infos {
			if v == v2 {
				isFind = true
				break
			}
		}
		if !isFind {
			ServiceUserInfoMod.AppendLog(ServiceUserInfoMod.ArgsAppendLog{
				InfoID:     v,
				OrgID:      data.OrgID,
				ChangeMark: "room.out",
				ChangeDes:  "离开房间",
				OldDes:     "",
				NewDes:     fmt.Sprint("入驻", data.Name, "[", data.ID, "]房间"),
			})
		}
	}
	//反馈
	return
}

// ArgsUpdateStatus 修改房间入驻状态参数
type ArgsUpdateStatus struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//状态
	// 0 空闲; 1 有人; 2 退房; 3 不可用; 4 清理中
	Status int `db:"status" json:"status"`
	//入驻人员列
	// 只有入驻时，从0->1时生效；如果为2则清理该人员信息
	Infos pq.Int64Array `db:"infos" json:"infos"`
}

// UpdateStatus 修改房间入驻状态
func UpdateStatus(args *ArgsUpdateStatus) (errCode string, err error) {
	var data FieldsRoom
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id, sort_id, tags, status, params FROM map_room WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.ID)
	if err != nil || data.ID < 1 {
		errCode = "room_not_exist"
		err = errors.New(fmt.Sprint("room not exist, ", err))
		return
	}
	switch args.Status {
	case 1:
		if len(args.Infos) < 1 {
			errCode = "infos_empty"
			err = errors.New("infos is empty")
			return
		}
		if data.Status != 0 {
			errCode = "room_used"
			err = errors.New("room have other infos use")
			return
		}
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_room SET status = :status, infos = :infos WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
		if err == nil {
			_ = appendLog(&argsAppendLog{
				OrgID:     data.OrgID,
				RoomID:    data.ID,
				MissionID: 0,
				Status:    0,
				Infos:     args.Infos,
				Des:       "入驻房间",
			})
		} else {
			errCode = "update"
			return
		}
	case 2:
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_room SET status = :status, infos = :infos WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
			"status": args.Status,
			"infos":  pq.Int64Array{},
		})
		if err == nil {
			_ = appendLog(&argsAppendLog{
				OrgID:     data.OrgID,
				RoomID:    data.ID,
				MissionID: 0,
				Status:    1,
				Infos:     data.Infos,
				Des:       "退房",
			})
		} else {
			errCode = "update"
			return
		}
	default:
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_room SET status = :status WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
			"id":     args.ID,
			"org_id": args.OrgID,
			"status": args.Status,
		})
		if err == nil {
			switch args.Status {
			case 0:
				_ = appendLog(&argsAppendLog{
					OrgID:     data.OrgID,
					RoomID:    data.ID,
					MissionID: 0,
					Status:    6,
					Infos:     []int64{},
					Des:       "房屋闲置",
				})
			case 3:
				_ = appendLog(&argsAppendLog{
					OrgID:     data.OrgID,
					RoomID:    data.ID,
					MissionID: 0,
					Status:    7,
					Infos:     []int64{},
					Des:       "房屋不可用",
				})
			case 4:
				_ = appendLog(&argsAppendLog{
					OrgID:     data.OrgID,
					RoomID:    data.ID,
					MissionID: 0,
					Status:    2,
					Infos:     []int64{},
					Des:       "清理房间中",
				})
			}
		} else {
			errCode = "update"
			return
		}
	}
	data = getRoomID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//删除缓冲
	deleteRoomCache(data.ID)
	//推送nats
	pushNatsUpdateStatus(data.ID, "update_status", "")
	//推送统计
	pushNatsUpdateAnalysis(data.OrgID)
	//反馈
	return
}

// ArgsUpdateServiceSortBindGroup 设置分类的默认分配ID参数
type ArgsUpdateServiceSortBindGroup struct {
	//分类ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//设置组织分组ID
	// 如果留空将无法自动分配任务
	GroupID int64 `db:"group_id" json:"groupID" check:"id" empty:"true"`
}

// UpdateServiceSortBindGroup 设置分类的默认分配ID
func UpdateServiceSortBindGroup(args *ArgsUpdateServiceSortBindGroup) (err error) {
	err = Sort.UpdateParamsAdd(&ClassSort.ArgsUpdateParams{
		ID:     args.ID,
		BindID: args.OrgID,
		Params: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "auto_org_group_id",
				Val:  fmt.Sprint(args.GroupID),
			},
		},
	})
	if err != nil {
		return
	}
	var data FieldsRoom
	data = getRoomID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//删除缓冲
	deleteRoomCache(data.ID)
	//推送nats
	pushNatsUpdateStatus(args.ID, "update_service_sort", "")
	//推送统计
	pushNatsUpdateAnalysis(data.OrgID)
	//反馈
	return
}

// ArgsUpdateServiceSortMissionExpire 设置分类的请求时间长度参数
type ArgsUpdateServiceSortMissionExpire struct {
	//分类ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//时间长度
	// 以秒为单位
	ExpireTime int64 `db:"expire_time" json:"expireTime" check:"int64Than0"`
}

// UpdateServiceSortMissionExpire 设置分类的请求时间长度
func UpdateServiceSortMissionExpire(args *ArgsUpdateServiceSortMissionExpire) (err error) {
	err = Sort.UpdateParamsAdd(&ClassSort.ArgsUpdateParams{
		ID:     args.ID,
		BindID: args.OrgID,
		Params: []CoreSQLConfig.FieldsConfigType{
			{
				Mark: "mission_expire_time",
				Val:  fmt.Sprint(args.ExpireTime),
			},
		},
	})
	if err != nil {
		return
	}
	//重新获取数据
	var data FieldsRoom
	data = getRoomID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//删除缓冲
	deleteRoomCache(data.ID)
	//推送nats
	pushNatsUpdateStatus(data.ID, "update_service_expire", "")
	//推送统计
	pushNatsUpdateAnalysis(data.OrgID)
	//反馈
	return
}

// ArgsUpdateServiceStatus 修改房间呼叫状态参数
type ArgsUpdateServiceStatus struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//服务呼叫状态
	// 0 无呼叫; 1 正在呼叫; 2 已经应答并处置
	// 处置完成后将回归0状态
	ServiceStatus int `db:"service_status" json:"serviceStatus" check:"intThan0" empty:"true"`
	//服务工作人员
	// 如果没有指定，将自动分配
	ServiceBindID int64 `db:"service_bind_id" json:"serviceBindID" check:"id" empty:"true"`
	//联动行政任务
	// 任务完成后将自动清除为0，否则将一直挂起
	// 如果没有指定，将自动生成
	ServiceMissionID int64 `db:"service_mission_id" json:"serviceMissionID" check:"id" empty:"true"`
}

// UpdateServiceStatus 修改房间呼叫状态
func UpdateServiceStatus(args *ArgsUpdateServiceStatus) (errCode string, err error) {
	var data FieldsRoom
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_id, sort_id, tags, service_status, service_bind_id, service_mission_id, params FROM map_room WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	if err != nil || data.ID < 1 {
		errCode = "room_not_exist"
		err = errors.New(fmt.Sprint("room not exist, ", err))
		return
	}
	//预备数据
	status := 0
	des := ""
	//根据类型区别处理
	switch args.ServiceStatus {
	case 1:
		if data.ServiceStatus == 1 {
			errCode = "service_replace"
			err = errors.New("service status is 1")
			return
		}
		/** 如果已应答，则强制进入呼叫状态即可，而不是拦截处理
		if data.ServiceStatus == 2 {
			errCode = "go_go_go"
			err = errors.New("service status is 2")
			return
		}
		*/
		if err != nil {
			errCode = "update"
			return
		}
		status = 4
		des = "呼叫服务"
		_, _ = AppendWarning(&ArgsAppendWarning{
			OrgID:    args.OrgID,
			RoomID:   args.ID,
			DeviceID: 0,
			NeedMQTT: false,
			CallType: 1,
		})
	case 2:
		if data.ServiceStatus == 2 {
			errCode = "service_replace"
			err = errors.New("service status is 2")
			return
		}
		status = 5
		des = "应答服务"
		if args.ServiceMissionID == 0 {
			var defaultGroupID int64
			defaultGroupID, err = Sort.GetParamInt64(&ClassSort.ArgsGetParam{
				ID:     data.SortID,
				BindID: data.OrgID,
				Mark:   "auto_org_group_id",
			})
			if err != nil || defaultGroupID < 1 {
				//errCode = "no_org_group"
				//return
				//退出switch处理模块
				break
			}
			var defaultExpireTime int64
			defaultExpireTime, err = Sort.GetParamInt64(&ClassSort.ArgsGetParam{
				ID:     data.SortID,
				BindID: data.OrgID,
				Mark:   "mission_expire_time",
			})
			if err != nil {
				defaultExpireTime = 1800
			}
			if args.ServiceBindID < 1 {
				var bindData OrgCore.FieldsBind
				bindData, err = OrgCore.GetBindLast(&OrgCore.ArgsGetBindLast{
					OrgID:   data.OrgID,
					GroupID: defaultGroupID,
					Mark:    "service",
					Params:  []CoreSQLConfig.FieldsConfigType{},
				})
				if err != nil {
					errCode = "no_bind_last"
					return
				} else {
					args.ServiceBindID = bindData.ID
				}
			}
			var missionData OrgMission.FieldsMission
			missionData, err = OrgMission.CreateMission(&OrgMission.ArgsCreateMission{
				OrgID:        data.OrgID,
				CreateBindID: args.ServiceBindID,
				BindID:       args.ServiceBindID,
				OtherBindIDs: []int64{},
				Title:        "房间需要提供服务支持",
				Des:          "房间居住人员或工作人员发起了服务请求，请前往或联系房间住户，完成此任务",
				DesFiles:     []int64{},
				StartAt:      CoreFilter.GetNowTime(),
				EndAt:        CoreFilter.GetNowTimeCarbon().AddSeconds(int(defaultExpireTime)).Time,
				TipID:        -1,
				ParentID:     0,
				Level:        0,
				SortID:       0,
				Tags:         []int64{},
				Params:       []CoreSQLConfig.FieldsConfigType{},
			})
			if err != nil {
				errCode = "create_mission"
				return
			}
			args.ServiceMissionID = missionData.ID
		}
		if args.ServiceBindID < 1 {
			err = errors.New("bind not find")
			errCode = "bind_not_find"
			return
		}
		if args.ServiceMissionID < 1 {
			//err = errors.New("mission not find")
			//errCode = "mission_not_find"
			//return
			//不一定非要具备任务
		}
	default:
		args.ServiceBindID = 0
		args.ServiceMissionID = 0
		status = 6
		des = "完成处理"
		_ = UnWarning(&ArgsAppendWarning{
			OrgID:    args.OrgID,
			RoomID:   args.ID,
			DeviceID: 0,
			NeedMQTT: false,
			CallType: 1,
		})
	}
	//更新
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_room SET service_status = :service_status, service_bind_id = :service_bind_id, service_mission_id = :service_mission_id WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":                 args.ID,
		"org_id":             args.OrgID,
		"service_status":     args.ServiceStatus,
		"service_bind_id":    args.ServiceBindID,
		"service_mission_id": args.ServiceMissionID,
	})
	if err != nil {
		return
	}
	//记录日志
	_ = appendLog(&argsAppendLog{
		OrgID:     data.OrgID,
		RoomID:    data.ID,
		MissionID: args.ServiceMissionID,
		Status:    status,
		Infos:     []int64{},
		Des:       des,
	})
	//清理缓冲
	deleteRoomCache(data.ID)
	//重新获取数据
	data, err = GetRoomID(&ArgsGetRoomID{
		ID:    args.ID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	//推送nats
	switch args.ServiceStatus {
	case 0:
		pushNatsUpdateStatus(data.ID, "service_status", "off")
	case 1:
		pushNatsUpdateStatus(data.ID, "service_status", "on")
	case 2:
		pushNatsUpdateStatus(data.ID, "service_status", "off")
		pushNatsServiceStatus("no", data.ID)
	}
	//推送统计
	pushNatsUpdateAnalysis(data.OrgID)
	//记录统计
	if data.ServiceStatus == 1 {
		var count int64
		_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM map_room_log WHERE status = 4 AND org_id = $1 AND room_id = $2", data.OrgID, data.ID)
		AnalysisAny2.AppendData("re", "map_room_service_count", time.Time{}, data.OrgID, 0, data.ID, 0, 0, count)
		for _, v := range data.Infos {
			AnalysisAny2.AppendData("re", "map_room_service_info_count", time.Time{}, data.OrgID, 0, v, 0, 0, count)
		}
		_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM map_room_log WHERE status = 6 AND org_id = $1 AND room_id = $2", data.OrgID, data.ID)
		AnalysisAny2.AppendData("re", "map_room_service_finish_count", time.Time{}, data.OrgID, 0, data.ID, 0, 0, count)
		for _, v := range data.Infos {
			AnalysisAny2.AppendData("re", "map_room_service_info_finish_count", time.Time{}, data.OrgID, 0, v, 0, 0, count)
		}
	}
	//清理缓冲
	deleteRoomCache(data.ID)
	//反馈
	return
}

// ArgsUpdateRoomColor 更新房间颜色参数
type ArgsUpdateRoomColor struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//背景色
	BgColor string `db:"bg_color" json:"bgColor" check:"color" empty:"true"`
}

// UpdateRoomColor 更新房间颜色
func UpdateRoomColor(args *ArgsUpdateRoomColor) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE map_room SET bg_color = :bg_color WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":       args.ID,
		"org_id":   args.OrgID,
		"bg_color": args.BgColor,
	})
	if err != nil {
		return
	}
	//清理缓冲
	deleteRoomCache(args.ID)
	//推送nats
	pushNatsUpdateStatus(args.ID, "bg_color", "")
	//反馈
	return
}
