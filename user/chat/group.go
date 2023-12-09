package UserChat

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetGroupList 获取群组列表参数
type ArgsGetGroupList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// -1 跳过
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetGroupList 获取群组列表
func GetGroupList(args *ArgsGetGroupList) (dataList []FieldsGroup, dataCount int64, err error) {
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
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	tableName := "user_chat_group"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, last_at, org_id, name, user_id, only_create_invite, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "last_at"},
	)
	return
}

// ArgsCreateGroup 创建新聊天室参数
type ArgsCreateGroup struct {
	//绑定组织
	// 商户可以查看构建的相关聊天室
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//聊天室创建人
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//只有创建人能邀请其他人？
	OnlyCreateInvite bool `db:"only_create_invite" json:"onlyCreateInvite"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateGroup 创建新聊天室
func CreateGroup(args *ArgsCreateGroup) (data FieldsGroup, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_chat_group", "INSERT INTO user_chat_group (last_at, org_id, name, user_id, only_create_invite, params) VALUES(NOW(),:org_id, :name, :user_id, :only_create_invite, :params)", map[string]interface{}{
		"org_id":             args.OrgID,
		"name":               args.Name,
		"user_id":            args.UserID,
		"only_create_invite": args.OnlyCreateInvite,
		"params":             args.Params,
	}, &data)
	if err != nil {
		return
	}
	err = InviteUser(&ArgsInviteUser{
		GroupID:      data.ID,
		UserID:       args.UserID,
		InviteUserID: args.UserID,
	})
	if err != nil {
		return
	}
	pushUpdateByUser(args.UserID, data.ID, 0)
	return
}

// ArgsDeleteGroup 删除聊天室参数
type ArgsDeleteGroup struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// -1 跳过
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// DeleteGroup 删除聊天室
func DeleteGroup(args *ArgsDeleteGroup) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_chat_group", "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", args)
	if err == nil {
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_chat_chat", "group_id = :group_id", map[string]interface{}{
			"group_id": args.ID,
		})
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_chat_message", "group_id = :group_id", map[string]interface{}{
			"group_id": args.ID,
		})
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_chat_message_money", "group_id = :group_id", map[string]interface{}{
			"group_id": args.ID,
		})
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_chat_message_ticket", "group_id = :group_id", map[string]interface{}{
			"group_id": args.ID,
		})
		if args.UserID > 0 {
			pushUpdateByUser(args.UserID, args.ID, 2)
		}
	}
	return
}
