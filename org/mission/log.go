package OrgMission

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 查询列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//绑定人信息
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//任务ID
	MissionID int64 `db:"mission_id" json:"missionID" check:"id" empty:"true"`
	//指定行为mark
	ContentMark string `json:"contentMark" check:"mark" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 查询列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.MissionID > -1 {
		where = where + " AND mission_id = :mission_id"
		maps["mission_id"] = args.MissionID
	}
	if args.ContentMark != "" {
		where = where + " AND content_mark = :content_mark"
		maps["content_mark"] = args.ContentMark
	}
	if args.Search != "" {
		where = where + " AND (content ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_mission_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, bind_id, mission_id, content_mark, content FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "content_mark"},
	)
	return
}

// ArgsCreateLog 插入数据参数
type ArgsCreateLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//任务ID
	MissionID int64 `db:"mission_id" json:"missionID" check:"id"`
	//操作内容标识码
	// 可用于其他语言处理
	ContentMark string `db:"content_mark" json:"contentMark" check:"mark"`
	//操作内容概述
	Content string `db:"content" json:"content" check:"des" min:"1" max:"1000"`
}

// CreateLog 插入数据
func CreateLog(args *ArgsCreateLog) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_mission_log (org_id, bind_id, mission_id, content_mark, content) VALUES (:org_id, :bind_id, :mission_id, :content_mark, :content)", args)
	return
}
