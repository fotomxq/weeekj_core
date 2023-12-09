package OrgCert2

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetWarningList 获取异常列表参数
type ArgsGetWarningList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否反馈
	NeedIsFinish bool `db:"need_is_finish" json:"needIsFinish" check:"bool"`
	IsFinish     bool `db:"is_finish" json:"isFinish" check:"bool"`
	//证件标识码
	ConfigMarks pq.StringArray `db:"config_marks" json:"configMarks" check:"marks" empty:"true"`
	//证件配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetWarningList 获取异常列表参数
func GetWarningList(args *ArgsGetWarningList) (dataList []FieldsWarning, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	where = CoreSQL.GetNeedChange(where, "finish_at", args.NeedIsFinish, args.IsFinish)
	if len(args.ConfigMarks) > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_mark = ANY(:config_marks)"
		maps["config_marks"] = args.ConfigMarks
	}
	if args.ConfigID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.TimeBetween.MinTime != "" && args.TimeBetween.MaxTime != "" {
		var timeBetween CoreSQLTime.FieldsCoreTime
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
		var newWhere string
		newWhere, maps = CoreSQLTime.GetBetweenByTime("create_at", timeBetween, maps)
		if newWhere != "" {
			if where != "" {
				where = where + " AND "
			}
			where = where + newWhere
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(msg ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "org_cert_warning2"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, finish_at, org_id, cert_id, config_id, config_mark, msg FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}
