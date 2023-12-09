package ToolsShortURL

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/robfig/cron"
	"time"
)

//短网址服务模块

var (
	//定时器
	runTimer      *cron.Cron
	runExpireLock = false
)

// ArgsGetList 查看列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetList 查看列表
func GetList(args *ArgsGetList) (dataList []FieldsShortURL, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if where == "" {
		where = " true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"tools_short_url",
		"id",
		"SELECT id, create_at, expire_at, key, org_id, user_id, is_public, data, params FROM tools_short_url WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "expire_at"},
	)
	return
}

// ArgsGetByID 查看ID参数
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可选
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetByID 查看ID
func GetByID(args *ArgsGetByID) (data FieldsShortURL, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, org_id, user_id, is_public, data, params FROM tools_short_url WHERE id = $1 AND (is_public = true OR (org_id = $2 AND user_id = $3))", args.ID, args.OrgID, args.UserID)
	return
}

// ArgsGetByKey 查询key参数
type ArgsGetByKey struct {
	//key
	Key string `db:"key" json:"key" check:"mark"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可选
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetByKey 查询key
func GetByKey(args *ArgsGetByKey) (data FieldsShortURL, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, expire_at, key, org_id, user_id, is_public, data, params FROM tools_short_url WHERE key = $1 AND (is_public = true OR (org_id = $2 AND user_id = $3))", args.Key, args.OrgID, args.UserID)
	return
}

// ArgsCreate 建立新的短域名参数
type ArgsCreate struct {
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否公开
	IsPublic bool `db:"is_public" json:"isPublic" check:"bool"`
	//存储的数据集合
	Data string `db:"data" json:"data"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Create 建立新的短域名
func Create(args *ArgsCreate) (data FieldsShortURL, err error) {
	if args.Data == "" {
		err = errors.New("data is empty")
		return
	}
	//重试3次
	tryCount := 0
	for {
		var key string
		key, err = CoreFilter.GetRandStr3(10)
		if err != nil {
			err = errors.New("rand failed")
			return
		}
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tools_short_url", "INSERT INTO tools_short_url (expire_at, key, org_id, user_id, is_public, data, params) VALUES (:expire_at, :key, :org_id, :user_id, :is_public, :data, :params)", map[string]interface{}{
			"expire_at": args.ExpireAt,
			"key":       key,
			"org_id":    args.OrgID,
			"user_id":   args.UserID,
			"is_public": args.IsPublic,
			"data":      args.Data,
			"params":    args.Params,
		}, &data)
		tryCount += 1
		if err != nil {
			if tryCount > 3 {
				return
			} else {
				err = nil
			}
		} else {
			break
		}
	}
	return
}

// ArgsDeleteByID 删除ID参数
type ArgsDeleteByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可选
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// DeleteByID 删除ID
func DeleteByID(args *ArgsDeleteByID) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "tools_short_url", "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:user_id < 1 OR user_id = :user_id)", args)
	return
}
