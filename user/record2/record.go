package UserRecord2

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetList 查询列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//成员ID
	OrgBindID int64 `json:"orgBindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//系统来源
	System string `json:"system" check:"mark" empty:"true"`
	//影响ID
	ModID int64 `json:"modID" check:"id" empty:"true"`
	//操作内容标识码
	Mark string `json:"mark" check:"mark" empty:"true"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 查询列表
func GetList(args *ArgsGetList) (dataList []FieldsRecord, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.OrgBindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.System != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "system = :system"
		maps["system"] = args.System
	}
	if args.ModID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mod_id = :mod_id"
		maps["mod_id"] = args.ModID
	}
	if args.Mark != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "user_record2"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, org_bind_id, user_id, system, mod_id, mark, des FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// argsAppendData 插入数据参数
type argsAppendData struct {
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//系统来源
	System string `db:"system" json:"system"`
	//影响ID
	ModID int64 `db:"mod_id" json:"modID"`
	//操作内容标识码
	Mark string `db:"mark" json:"mark"`
	//操作内容概述
	Des string `db:"des" json:"des"`
}

// appendData 插入数据
func appendData(args *argsAppendData) (err error) {
	//创建日志
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_record2 (org_id, org_bind_id, user_id, system, mod_id, mark, des) VALUES (:org_id, :org_bind_id, :user_id, :system, :mod_id, :mark, :des)", args)
	if err != nil {
		return
	}
	return
}
