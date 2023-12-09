package UserChat

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
)

// ArgsGetChatList 获取成员列表参数
type ArgsGetChatList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//聊天室ID
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//获取该数据的用户ID
	// -1 跳过
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否包含已经离开的人
	HaveLeave bool `json:"haveLeave" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetChatList 获取成员列表
func GetChatList(args *ArgsGetChatList) (dataList []FieldsChat, dataCount int64, err error) {
	//检查用户
	if args.UserID > 0 {
		if !checkChatUser(args.GroupID, args.UserID) {
			err = errors.New("no permission")
			return
		}
	}
	where := "group_id = :group_id"
	maps := map[string]interface{}{
		"group_id": args.GroupID,
	}
	if !args.HaveLeave {
		where = where + " AND leave_at < to_timestamp(1000000)"
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "user_chat_chat"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, leave_at, user_id, name, group_id, un_read_count FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "leave_at"},
	)
	return
}

// ArgsInviteUser 邀请其他人加入聊天室参数
type ArgsInviteUser struct {
	//聊天室ID
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//邀请人
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//被邀请人
	InviteUserID int64 `db:"invite_user_id" json:"inviteUserID" check:"id"`
}

// InviteUser 邀请其他人加入聊天室
func InviteUser(args *ArgsInviteUser) (err error) {
	//检查是否在聊天室？
	var chatID int64
	err = Router2SystemConfig.MainDB.Get(&chatID, "SELECT id FROM user_chat_chat WHERE group_id = $1 AND user_id = $2 AND leave_at < to_timestamp(1000000)", args.GroupID, args.InviteUserID)
	if err == nil && chatID > 0 {
		err = nil
		return
	}
	//检查房间是否存在且符合基本条件
	var groupID int64
	err = Router2SystemConfig.MainDB.Get(&groupID, "SELECT id FROM user_chat_group WHERE id = $1 AND (only_create_invite = false OR (only_create_invite = true AND user_id = $2))", args.GroupID, args.UserID)
	if err != nil || groupID < 1 {
		err = errors.New("group not exist")
		return
	}
	//检查是否已经在列表中？如果存在则修改离开时间
	if chatID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_chat_group SET leave_at = to_timestamp(0) AND id = :id", map[string]interface{}{
			"id": args.GroupID,
		})
	} else {
		var userData UserCore.FieldsUserType
		userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
			ID:    args.InviteUserID,
			OrgID: -1,
		})
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_chat_chat (user_id, name, group_id) VALUES (:user_id,:name,:group_id)", map[string]interface{}{
			"user_id":  args.InviteUserID,
			"name":     userData.Name,
			"group_id": args.GroupID,
		})
	}
	if err == nil {
		updateLastAtByChat(args.GroupID, args.UserID)
		pushUpdateByGroup(args.GroupID, 0)
	}
	//反馈
	return
}

// ArgsOutChat 退出聊天参数
type ArgsOutChat struct {
	//聊天室ID
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//房间的创建人
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//被剔除的人
	InviteUserID int64 `db:"invite_user_id" json:"inviteUserID" check:"id"`
}

// OutChat 退出聊天
func OutChat(args *ArgsOutChat) (err error) {
	//检查是否为剔除模式
	if args.UserID > 0 {
		//检查房间是否存在且符合基本条件
		var groupID int64
		err = Router2SystemConfig.MainDB.Get(&groupID, "SELECT id FROM user_chat_group WHERE id = $1 AND user_id = $2", args.GroupID, args.UserID)
		if err != nil || groupID < 1 {
			err = errors.New("group not exist")
			return
		}
	}
	//检查是否在聊天室？
	var chatID int64
	err = Router2SystemConfig.MainDB.Get(&chatID, "SELECT id FROM user_chat_chat WHERE group_id = $1 AND user_id = $2 AND leave_at < to_timestamp(1000000)", args.GroupID, args.InviteUserID)
	if err == nil && chatID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_chat_group SET leave_at = to_timestamp(0) AND id = :id", map[string]interface{}{
			"id": chatID,
		})
	}
	//检查聊天室是否只剩下1个人，则自动关闭聊天室
	var count int64
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM user_chat_chat WHERE group_id = $1 AND leave_at < to_timestamp(1000000)", args.GroupID)
	if count < 2 {
		err = DeleteGroup(&ArgsDeleteGroup{
			ID:     args.GroupID,
			OrgID:  -1,
			UserID: -1,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("delete chat group failed, ", err))
			return
		}
	}
	if err == nil {
		pushUpdateByGroup(args.GroupID, 1)
	}
	//反馈
	return
}

// ArgsUpdateChatName 修改成员姓名参数
type ArgsUpdateChatName struct {
	//聊天室ID
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//成员用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//新的姓名
	Name string `db:"name" json:"name" check:"name"`
}

// UpdateChatName 修改成员姓名
func UpdateChatName(args *ArgsUpdateChatName) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_chat_chat SET name = :name WHERE group_id = :group_id AND user_id = :user_id", args)
	if err == nil {
		pushUpdateByGroup(args.GroupID, 3)
	}
	return
}

// ArgsUpdateChatRead 标记已读参数
type ArgsUpdateChatRead struct {
	//聊天室ID
	GroupID int64 `db:"group_id" json:"groupID" check:"id"`
	//成员用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// UpdateChatRead 标记已读
func UpdateChatRead(args *ArgsUpdateChatRead) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_chat_chat SET un_read_count = 0 WHERE group_id = :group_id AND user_id = :user_id", args)
	return
}

// 检查用户是否在房间内
func checkChatUser(groupID int64, userID int64) (b bool) {
	//检查用户
	var id int64
	err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_chat_chat WHERE user_id = $1 AND group_id = $2 AND leave_at < to_timestamp(1000000)", userID, groupID)
	if err != nil || id < 1 {
		return
	}
	return true
}

// 更新最后一次时间
func updateLastAtByChat(groupID, userID int64) {
	_, err := CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_chat_chat SET last_at = NOW(), un_read_count = un_read_count + 1 WHERE groupID = :group_id AND userID = :user_id", map[string]interface{}{
		"groupID": groupID,
		"user_id": userID,
	})
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_chat_group SET last_at = NOW() WHERE id = :id", map[string]interface{}{
		"id": groupID,
	})
	return
}
