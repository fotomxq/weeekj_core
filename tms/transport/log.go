package TMSTransport

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//数据产生来源
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID" check:"id" empty:"true"`
	//配送人员
	TransportBindID int64 `db:"transport_bind_id" json:"transportBindID" check:"id" empty:"true"`
	//行为特征
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.TransportID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "transport_id = :transport_id"
		maps["transport_id"] = args.TransportID
	}
	if args.TransportBindID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "transport_bind_id = :transport_bind_id"
		maps["transport_bind_id"] = args.TransportBindID
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
	tableName := "tms_transport_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, bind_id, transport_id, transport_bind_id, mark, des FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// argsAppendLog 创建日志参数
type argsAppendLog struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//数据产生来源
	BindID int64 `db:"bind_id" json:"bindID"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID"`
	//配送人员
	TransportBindID int64 `db:"transport_bind_id" json:"transportBindID"`
	//行为特征
	Mark string `db:"mark" json:"mark"`
	//备注
	Des string `db:"des" json:"des"`
}

// appendLog 创建日志
func appendLog(args *argsAppendLog) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO tms_transport_log (org_id, bind_id, transport_id, transport_bind_id, mark, des) VALUES (:org_id,:bind_id,:transport_id,:transport_bind_id,:mark,:des)", args)
	return
}
