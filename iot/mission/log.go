package IOTMission

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取指定任务的日志数据参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//任务ID
	MissionID int64 `db:"mission_id" json:"missionID" check:"id" empty:"true"`
	//状态
	Status int `db:"status" json:"status" check:"intThan0" empty:"true"`
	//行为标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否为历史
	IsHistory bool `db:"is_history" json:"isHistory" check:"bool"`
	//搜索
	Search string `db:"search" json:"search" check:"search" empty:"true"`
}

// GetLogList 获取指定任务的日志数据
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		var missionData FieldsMission
		err = Router2SystemConfig.MainDB.Get(&missionData, "SELECT id FROM iot_mission WHERE org_id = $1", args.OrgID)
		if err != nil || missionData.ID < 1 {
			err = errors.New("mission not org")
			return
		}
	}
	if args.MissionID > 0 {
		where = where + "mission_id = :mission_id"
		maps["mission_id"] = args.MissionID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(content ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "iot_mission_log"
	if args.IsHistory {
		tableName = "iot_mission_log_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, mission_id, status, mark, content FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsCreateLog 创建新的日志数据参数
type ArgsCreateLog struct {
	//任务ID
	MissionID int64 `db:"mission_id" json:"missionID" check:"id"`
	//状态
	// 0 wait 等待发起 / 1 send 已经发送 / 2 success 已经完成 / 3 failed 已经失败 / 4 cancel 取消
	Status int `db:"status" json:"status" check:"intThan0"`
	//行为标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//日志内容
	Content string `db:"content" json:"content" check:"des" min:"1" max:"1000" empty:"true"`
}

// CreateLog 创建新的日志数据
func CreateLog(args *ArgsCreateLog) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO iot_mission_log (mission_id, status, mark, content) VALUES (:mission_id,:status,:mark,:content)", args)
	return
}
