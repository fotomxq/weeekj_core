package OrgMap

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetMapList 获取地图列表参数
type ArgsGetMapList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//上级ID
	// 用于叠加展示
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//是否审核
	NeedIsAudit bool `json:"needIsAudit" check:"bool"`
	IsAudit     bool `json:"isAudit" check:"bool"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetMapList 获取地图列表
func GetMapList(args *ArgsGetMapList) (dataList []FieldsMap, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.NeedIsAudit {
		if args.IsAudit {
			where = where + " AND audit_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at < to_timestamp(1000000)"
		}
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_map"
	var rawList []FieldsMap
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "audit_at", "ad_count"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getMapByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetMapByID 根据ID获取数据参数
type ArgsGetMapByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetMapByID 根据ID获取数据
func GetMapByID(args *ArgsGetMapByID) (data FieldsMap, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_map WHERE id = $1 AND delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000)", args.ID)
	if err != nil {
		return
	}
	data = getMapByID(data.ID)
	return
}

func GetMapNoAuditByID(args *ArgsGetMapByID) (data FieldsMap, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_map WHERE id = $1", args.ID)
	if err != nil {
		return
	}
	data = getMapByID(data.ID)
	return
}

func GetMapNameByID(args *ArgsGetMapByID) (data string, err error) {
	if args.ID < 1 {
		err = errors.New("no data")
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT name FROM org_map WHERE id = $1 AND delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000)", args.ID)
	if err != nil {
		return
	}
	return
}

// ArgsGetMapByOrg 获取商户信息参数
type ArgsGetMapByOrg struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//是否审核
	IsAudit bool `json:"isAudit"`
}

// GetMapByOrg 获取商户信息
func GetMapByOrg(args *ArgsGetMapByOrg) (data FieldsMap, err error) {
	if args.IsAudit {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_map WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000) LIMIT 1", args.OrgID)
	} else {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_map WHERE org_id = $1 AND delete_at < to_timestamp(1000000) LIMIT 1", args.OrgID)
	}
	if err != nil {
		return
	}
	data = getMapByID(data.ID)
	return
}

// GetMapCountByOrgOrUser 获取用户或商户的广告数量
func GetMapCountByOrgOrUser(orgID int64, userID int64) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM org_map WHERE (org_id = $1 OR user_id = $2) AND delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000)", orgID, userID)
	return
}

// ArgsGetMapByUser 获取用户信息参数
type ArgsGetMapByUser struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//是否审核
	IsAudit bool `json:"isAudit"`
}

// GetMapByUser 获取商户信息
func GetMapByUser(args *ArgsGetMapByUser) (data FieldsMap, err error) {
	if args.IsAudit {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_map WHERE user_id = $1 AND delete_at < to_timestamp(1000000) AND audit_at > to_timestamp(1000000) LIMIT 1", args.UserID)
	} else {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_map WHERE user_id = $1 AND delete_at < to_timestamp(1000000) LIMIT 1", args.UserID)
	}
	if err != nil {
		return
	}
	data = getMapByID(data.ID)
	return
}

// ArgsGetMapChildCount 检查上级的数量参数
type ArgsGetMapChildCount struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
}

// GetMapChildCount 检查上级的数量
func GetMapChildCount(args *ArgsGetMapChildCount) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM org_map WHERE parent_id = $1 AND delete_at < to_timestamp(1000000)", args.ID)
	return
}

// 获取地图数据包
func getMapByID(id int64) (data FieldsMap) {
	cacheMark := getMapCacheMark(id)
	_ = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data)
	if data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, audit_at, org_id, user_id, parent_id, cover_file_id, cover_file_ids, name, des, country, province, city, address, map_type, longitude, latitude, ad_count, ad_count_limit, view_time_limit, params FROM org_map WHERE id = $1 LIMIT 1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}
