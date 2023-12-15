package UserRecordCore

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserRecord2Mod "github.com/fotomxq/weeekj_core/v5/user/record2/mod"
)

//日志记录
// 本方法用于用户和设备等渠道端的操作记录内容，用于运维和运营快速了解目标的操作记录，以对提供相关支持。
// 需要注意，本方法不会反馈系统错误、其他关联错误、底层操作记录等信息，对部分操作可能存在遗漏，不能作为唯一的参考依据。

// ArgsGetList 查询列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//指定行为mark
	ContentMark string `json:"contentMark" check:"mark" empty:"true"`
	//是否为历史数据
	IsHistory bool `json:"isHistory" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 查询列表
func GetList(args *ArgsGetList) (dataList []FieldsRecordType, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.UserID > 0 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ContentMark != "" {
		where = where + " AND content_mark = :content_mark"
		maps["content_mark"] = args.ContentMark
	}
	tableName := "user_record"
	if args.IsHistory {
		tableName = tableName + "_history"
	}
	if args.Search != "" {
		where = where + " AND (username ILIKE '%' || :search || '%' OR content_mark ILIKE '%' || :search || '%' OR content ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, user_id, username, content_mark, content FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "username", "content_mark"},
	)
	return
}

// ArgsCreate 插入数据参数
type ArgsCreate struct {
	//组织ID
	// 留空则表明为平台所有
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户信息
	UserID   int64  `db:"user_id" json:"userID"`
	UserName string `db:"username" json:"username"`
	//记录内容
	ContentMark string `db:"content_mark" json:"contentMark"`
	Content     string `db:"content" json:"content"`
}

// Create 插入数据
func Create(args *ArgsCreate) (err error) {
	UserRecord2Mod.AppendData(args.OrgID, 0, args.UserID, "", 0, args.ContentMark, args.Content)
	/** 废弃代码
	//获取用户数据
	if args.UserID > 0 && args.UserName == "" {
		var userData UserCore.FieldsUserType
		userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
			ID:    args.UserID,
			OrgID: -1,
		})
		if err == nil && userData.ID > 0 {
			args.OrgID = userData.OrgID
			args.UserName = userData.Name
		}
	}
	//创建日志
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_record (org_id, user_id, username, content_mark, content) VALUES (:org_id,:user_id,:username,:content_mark,:content)", args)
	if err != nil {
		return
	}
	*/
	return
}
