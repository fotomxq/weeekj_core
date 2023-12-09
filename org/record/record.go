package OrgRecord

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserRecord2Mod "gitee.com/weeekj/weeekj_core/v5/user/record2/mod"
	"github.com/lib/pq"
)

//日志记录
// 本方法用于用户和设备等渠道端的操作记录内容，用于运维和运营快速了解目标的操作记录，以对提供相关支持。
// 需要注意，本方法不会反馈系统错误、其他关联错误、底层操作记录等信息，对部分操作可能存在遗漏，不能作为唯一的参考依据。

// ArgsGetList 查询列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//模块标识码
	FromMark string `db:"from_mark" json:"fromMark" check:"mark" empty:"true"`
	//修改内容ID
	FromID int64 `db:"from_id" json:"fromID" check:"id" empty:"true"`
	//绑定人信息
	BindID int64 `json:"bindID" check:"id" empty:"true"`
	//指定行为mark
	ContentMarks pq.StringArray `json:"contentMarks" check:"marks" empty:"true"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 查询列表
func GetList(args *ArgsGetList) (dataList []FieldsRecordType, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.FromMark != "" {
		where = where + " AND from_mark = :from_mark"
		maps["from_mark"] = args.FromMark
	}
	if args.FromID > 0 {
		where = where + " AND from_id = :from_id"
		maps["from_id"] = args.FromID
	}
	if args.BindID > 0 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if len(args.ContentMarks) > 0 {
		where = where + " AND content_mark = ANY(:content_marks)"
		maps["content_marks"] = args.ContentMarks
	}
	if args.Search != "" {
		where = where + " AND (content ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_record"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, from_mark, bind_id, content_mark, content, change_data FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "content_mark"},
	)
	return
}

// ArgsCreate 插入数据参数
type ArgsCreate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//模块标识码
	FromMark string `db:"from_mark" json:"fromMark"`
	//修改内容ID
	FromID int64 `db:"from_id" json:"fromID"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//操作内容标识码
	// 可用于其他语言处理
	ContentMark string `db:"content_mark" json:"contentMark"`
	//操作内容概述
	Content string `db:"content" json:"content"`
	//变动内容列
	ChangeData FieldsRecordChangeList `db:"change_data" json:"changeData"`
}

// Create 插入数据
func Create(args *ArgsCreate) (err error) {
	UserRecord2Mod.AppendData(args.OrgID, args.BindID, 0, args.FromMark, args.FromID, args.ContentMark, args.Content)
	/** 废弃代码
	if args.ChangeData == nil {
		args.ChangeData = FieldsRecordChangeList{}
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_record (org_id, from_mark, from_id, bind_id, content_mark, content, change_data) VALUES (:org_id, :from_mark, :from_id, :bind_id, :content_mark, :content, :change_data)", args)
	*/
	return
}
