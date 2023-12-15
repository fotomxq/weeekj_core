package OrgWorkTip

import (
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetMsgList 获取通知列表参数
type ArgsGetMsgList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//是否已读
	IsRead bool `db:"is_read" json:"isRead"`
}

// GetMsgList 获取通知列表
func GetMsgList(args *ArgsGetMsgList) (dataList []FieldsWorkTip, dataCount int64, err error) {
	where := "is_read = :is_read"
	maps := map[string]interface{}{
		"is_read": args.IsRead,
	}
	if args.OrgBindID > -1 {
		where = where + " AND org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
	}
	tableName := "org_work_tip"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, org_bind_id, msg, system, bind_id, is_read FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	return
}

// GetLastMsg 获取最近一条通知
func GetLastMsg(orgBindID int64) (data FieldsWorkTip, count int64) {
	cacheMark := getLastCacheMark(orgBindID)
	cacheCountMark := getLastCacheMark(orgBindID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		count, _ = Router2SystemConfig.MainCache.GetInt64(cacheCountMark)
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, msg, system, bind_id, is_read FROM org_work_tip WHERE org_bind_id = $1 AND is_read = false ORDER BY id DESC LIMIT 1", orgBindID)
	if err != nil {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM org_work_tip WHERE org_bind_id = $1 AND is_read = false", orgBindID)
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 3600)
	Router2SystemConfig.MainCache.SetInt64(cacheCountMark, count, 3600)
	return
}

// ReadID 阅读指定ID
func ReadID(id int64, orgBindID int64) (err error) {
	var data FieldsWorkTip
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, org_bind_id FROM org_work_tip WHERE id = $1", id)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_work_tip SET is_read = true WHERE id = :id AND org_bind_id = :org_bind_id", map[string]interface{}{
		"id":          data.ID,
		"org_bind_id": orgBindID,
	})
	if err != nil {
		return
	}
	deleteCache(data.OrgBindID)
	return
}

// argsAppendTip 添加一个数据
type argsAppendTip struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID"`
	//消息内容
	Msg string `db:"msg" json:"msg"`
	//系统
	System string `db:"system" json:"system"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID"`
}

func appendTip(args argsAppendTip) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_work_tip (org_id, org_bind_id, msg, system, bind_id, is_read) VALUES (:org_id,:org_bind_id,:msg,:system,:bind_id,false)", args)
	if err != nil {
		return
	}
	deleteCache(args.OrgBindID)
	return
}

// 缓冲
func getLastCacheMark(orgBindID int64) string {
	return fmt.Sprint("org:work:tip:last:", orgBindID)
}

func getCountCacheMark(orgBindID int64) string {
	return fmt.Sprint("org:work:tip:count:", orgBindID)
}

func deleteCache(orgBindID int64) {
	Router2SystemConfig.MainCache.DeleteMark(getCountCacheMark(orgBindID))
	Router2SystemConfig.MainCache.DeleteMark(getLastCacheMark(orgBindID))
}
