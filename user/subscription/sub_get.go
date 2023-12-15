package UserSubscription

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetSubList 获取订阅列表参数
type ArgsGetSubList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否到期
	NeedIsExpire bool `json:"needIsExpire" check:"bool" empty:"true"`
	IsExpire     bool `json:"isExpire" check:"bool" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetSubList 获取订阅列表
func GetSubList(args *ArgsGetSubList) (dataList []FieldsSub, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > -1 {
		where = where + " AND config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.NeedIsExpire {
		if args.IsExpire {
			where = where + " AND expire_at < NOW()"
		} else {
			where = where + " AND expire_at >= NOW()"
		}
	}
	if args.Search != "" {
		var configList []FieldsConfig
		err = Router2SystemConfig.MainDB.Select(&configList, "SELECT id FROM user_sub_config WHERE delete_at < to_timestamp(1000000) AND (title ILIKE '%' || $1 || '%' OR des ILIKE '%' || $1 || '%')", args.Search)
		if err == nil && len(configList) > 0 {
			var configIDs pq.Int64Array
			for _, v := range configList {
				configIDs = append(configIDs, v.ID)
			}
			where = where + " AND config_id = ANY(:config_ids)"
			maps["config_ids"] = configIDs
		}
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_sub",
		"id",
		"SELECT id, create_at, update_at, delete_at, expire_at, org_id, config_id, user_id, params FROM user_sub WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "expire_at"},
	)
	return
}

// ArgsGetSub 获取指定人的订阅信息参数
type ArgsGetSub struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// GetSub 获取指定人的订阅信息
func GetSub(args *ArgsGetSub) (data FieldsSub, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, config_id, user_id, params FROM user_sub WHERE ($1 < 1 OR config_id = $1) AND user_id = $2 AND delete_at < to_timestamp(1000000) LIMIT 1", args.ConfigID, args.UserID)
	if err == nil && data.ID < 1 {
		err = errors.New("data not exist")
	}
	return
}

func GetSubNoErr(userID int64, configID int64) (data FieldsSub) {
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, config_id, user_id, params FROM user_sub WHERE ($1 < 1 OR config_id = $1) AND user_id = $2 AND delete_at < to_timestamp(1000000) LIMIT 1", configID, userID)
	return
}

// CheckHaveAnySub 用户是否具有任意会员？
func CheckHaveAnySub(userID int64) (b bool) {
	var id int64
	err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_sub WHERE delete_at < to_timestamp(1000000) AND expire_at >= NOW() AND user_id = $1 ORDER BY id DESC LIMIT 1", userID)
	if err != nil || id < 1 {
		return
	}
	b = true
	return
}

// 获取订阅的hash
func getSubHash(subData *FieldsSub) string {
	return CoreFilter.GetSha1Str(fmt.Sprint(subData.ID, ".", subData.OrgID, ".", subData.UserID, ".", subData.ConfigID, ".", subData.ExpireAt))
}
