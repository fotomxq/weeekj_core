package UserMessage

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//发送用户ID
	SendUserID int64 `db:"send_user_id" json:"sendUserID" check:"id" empty:"true"`
	//预计发送时间
	WaitSendAt time.Time `db:"wait_send_at" json:"waitSendAt" check:"isoTime" empty:"true"`
	//发送状态
	// 发送完成后，发送人无法删除，接收人可以标记已读或删除，具体其他字段完成该操作
	// 0 草稿; 1 等待审核; 2 发送成功
	Status int `db:"status" json:"status"`
	//接收人ID
	ReceiveUserID int64 `db:"receive_user_id" json:"receiveUserID" check:"id" empty:"true"`
	//接收人阅读时间
	ReceiveReadAt time.Time `db:"receive_read_at" json:"receiveReadAt" check:"isoTime" empty:"true"`
	//接收人删除状态
	ReceiveDeleteAt time.Time `db:"receive_delete_at" json:"receiveDeleteAt" check:"isoTime" empty:"true"`
	//是否被删除
	NeedIsRemove bool `json:"needIsRemove" check:"bool" empty:"true"`
	IsRemove     bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsMessage, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.NeedIsRemove {
		where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	}
	if args.SendUserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "send_user_id = :send_user_id"
		maps["send_user_id"] = args.SendUserID
	}
	if args.Status > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "status = :status"
		maps["status"] = args.Status
	}
	if CoreSQL.CheckTimeHaveData(args.WaitSendAt) {
		if where != "" {
			where = where + " AND "
		}
		where = where + "wait_send_at > :wait_send_at"
		maps["wait_send_at"] = args.WaitSendAt
	}
	if args.ReceiveUserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		if args.ReceiveReadAt.Unix() > 0 {
			where = where + "receive_read_at > :receive_read_at"
			maps["receive_read_at"] = args.ReceiveReadAt
		}
		if args.ReceiveDeleteAt.Unix() > 0 {
			where = where + "receive_delete_at > :receive_delete_at"
			maps["receive_delete_at"] = args.ReceiveDeleteAt
		}
		where = where + "receive_user_id = :receive_user_id"
		maps["receive_user_id"] = args.ReceiveUserID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(title ILIKE '%' || :search || '%' OR content ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	var rawList []FieldsMessage
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"user_message",
		"id",
		"SELECT id FROM user_message WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "wait_send_at", "receive_read_at", "receive_delete_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		vData.Content = ""
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetByID 获取ID参数
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//发送人ID
	SendUserID int64 `db:"send_user_id" json:"sendUserID" check:"id" empty:"true"`
	//接收人ID
	ReceiveUserID int64 `db:"receive_user_id" json:"receiveUserID" check:"id" empty:"true"`
}

// GetByID 获取ID
func GetByID(args *ArgsGetByID) (data FieldsMessage, err error) {
	data = getByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.SendUserID, data.SendUserID) && !CoreFilter.EqID2(args.ReceiveUserID, data.ReceiveUserID) {
		data = FieldsMessage{}
		err = errors.New("no data")
		return
	}
	//如果是接收人
	if CoreFilter.EqID2(args.ReceiveUserID, data.ReceiveUserID) {
		_ = UpdateReceiveRead(&ArgsUpdateReceiveRead{
			ID:            data.ID,
			ReceiveUserID: -1,
		})
		deleteMessageCache(data.ID)
	}
	//反馈
	return
}

// GetReceiveCountByUserID 获取用户收到的消息数量
func GetReceiveCountByUserID(userID int64) (count int64) {
	cacheMark := getReceiveMessageCountCacheMark(userID)
	var err error
	count, err = Router2SystemConfig.MainCache.GetInt64(cacheMark)
	if err == nil && count > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT (id) FROM user_message WHERE receive_user_id = $1 AND wait_send_at >= to_timestamp(1000000) AND status = 2", userID)
	Router2SystemConfig.MainCache.SetInt64(cacheMark, count, 1800)
	return
}

// 获取消息
func getByID(id int64) (data FieldsMessage) {
	cacheMark := getMessageCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, wait_send_at, status, send_user_id, receive_user_id, receive_read_at, receive_delete_at, title, content, files, params FROM user_message WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 21600)
	return
}
