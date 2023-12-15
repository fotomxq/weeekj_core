package OrgSubscription

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetSubList 获取订阅列表参数
type ArgsGetSubList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//开通配置
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//是否到期
	NeedIsExpire bool `db:"need_is_expire" json:"needIsExpire" check:"bool"`
	IsExpire     bool `db:"is_expire" json:"isExpire" check:"bool"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetSubList 获取订阅列表
func GetSubList(args *ArgsGetSubList) (dataList []FieldsSub, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > 0 {
		where = where + " AND config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.NeedIsExpire {
		if args.IsExpire {
			where = where + " AND expire_at < NOW()"
		} else {
			where = where + " AND expire_at >= NOW()"
		}
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"org_sub",
		"id",
		"SELECT id, create_at, update_at, delete_at, expire_at, org_id, config_id, params FROM org_sub WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "expire_at"},
	)
	return
}

// ArgsCheckSub 检查组织的订阅状态参数
type ArgsCheckSub struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//开通配置
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
}

// CheckSub 检查组织的订阅状态
func CheckSub(args *ArgsCheckSub) (expireAt time.Time, b bool) {
	var data FieldsSub
	if err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, expire_at FROM org_sub WHERE org_id = $1 AND config_id = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.ConfigID); err != nil {
		return
	}
	if data.ID < 1 {
		return
	}
	expireAt = data.ExpireAt
	b = expireAt.Unix() >= CoreFilter.GetNowTime().Unix()
	return
}

// 获取订阅
func getSub(id int64) (data FieldsSub, err error) {
	if err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, org_id, config_id, params FROM org_sub WHERE id = $1 AND delete_at < to_timestamp(1000000)", id); err != nil {
		return
	}
	return
}

// 获取订阅的hash
func getSubHash(subData *FieldsSub) string {
	return CoreFilter.GetSha1Str(fmt.Sprint(subData.ID, ".", subData.OrgID, ".", subData.ConfigID, ".", subData.ExpireAt))
}
