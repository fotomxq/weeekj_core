package MapRoom

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID" check:"id" empty:"true"`
	//状态
	// 0 进入; 1 退出; 2 退房中(核对清理状态); 3 清理中; 4 呼叫中; 5 已经应答并发出任务; 6 任务处理完成
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//入驻人员
	InfoID int64 `db:"info_id" json:"infoID" check:"id" empty:"true"`
	//任务ID
	MissionID int64 `db:"mission_id" json:"missionID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
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
	if args.Status > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "status = :status"
		maps["status"] = args.Status
	}
	if args.InfoID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + ":info_id = ANY(infos)"
		maps["info_id"] = args.InfoID
	}
	if args.MissionID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mission_id = :mission_id"
		maps["mission_id"] = args.MissionID
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "map_room_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, room_id, mission_id, status, infos, des FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// 添加日志数据
type argsAppendLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//房间ID
	RoomID int64 `db:"room_id" json:"roomID"`
	//任务ID
	// 可能为0
	MissionID int64 `db:"mission_id" json:"missionID"`
	//状态
	// 0 进入; 1 退房中; 2 清理中; 3 呼叫中; 4 已经应答并发出任务; 5 任务处理完成; 6 闲置; 7 不可用
	Status int `db:"status" json:"status"`
	//入驻人员列
	Infos pq.Int64Array `db:"infos" json:"infos"`
	//备注
	Des string `db:"des" json:"des"`
}

func appendLog(args *argsAppendLog) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO map_room_log (org_id, room_id, mission_id, status, infos, des) VALUES (:org_id,:room_id,:mission_id,:status,:infos,:des)", args)
	return
}
