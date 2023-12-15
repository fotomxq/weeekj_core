package BlogStuRead

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetLogList 获取日志列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//学习内容
	ContentID int64 `db:"content_id" json:"contentID" check:"id" empty:"true"`
}

// GetLogList 获取日志列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ContentID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "content_id = :content_id"
		maps["content_id"] = args.ContentID
	}
	if where == "" {
		where = "true"
	}
	tableName := "blog_stu_read_log"
	var rawList []FieldsLog
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "end_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getLog(v.ID)
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

func getLog(id int64) (data FieldsLog) {
	cacheMark := getLogCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, end_at, run_time, org_id, user_id, content_id FROM blog_stu_read_log WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 604800)
	return
}
