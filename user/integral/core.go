package UserIntegral

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetList 获取积分列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//积分范围
	Min int64 `json:"min" check:"int64Than0" empty:"true"`
	Max int64 `json:"max" check:"int64Than0" empty:"true"`
}

// GetList 获取积分列表
func GetList(args *ArgsGetList) (dataList []FieldsIntegral, dataCount int64, err error) {
	var where string
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
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
	if args.Min > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "count >= :min"
		maps["min"] = args.Min
	}
	if args.Max > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "count <= :max"
		maps["max"] = args.Max
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_integral",
		"id",
		"SELECT id, create_at, update_at, org_id, user_id, count FROM user_integral WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "count"},
	)
	return
}

// ArgsGetLogList 查看积分变动记录参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//请求变动范围
	// -999999~999999
	Min int64 `json:"min"`
	Max int64 `json:"max"`
	//搜索
	// 备注
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 查看积分变动记录
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
	var where string
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
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
	if args.Min > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "add_count >= :min"
		maps["min"] = args.Min
	}
	if args.Max > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "add_count <= :max"
		maps["max"] = args.Max
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_integral_log",
		"id",
		"SELECT id, create_at, org_id, user_id, add_count, des FROM user_integral_log WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "add_count"},
	)
	return
}

// ArgsGetUser 查看某个用户的积分参数
type ArgsGetUser struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetUser 查看某个用户的积分
func GetUser(args *ArgsGetUser) (data FieldsIntegral, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, org_id, user_id, count FROM user_integral WHERE org_id = $1 AND user_id = $2", args.OrgID, args.UserID)
	return
}

// GetUserCount 查看某个用户的积分，只要积分部分
func GetUserCount(orgID int64, userID int64) (count int64) {
	var data FieldsIntegral
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT count FROM user_integral WHERE org_id = $1 AND user_id = $2", orgID, userID)
	if err == nil {
		count = data.Count
	}
	return
}

// ArgsAddCount 变动积分参数
type ArgsAddCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//变动
	AddCount int64 `db:"add_count" json:"addCount"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
}

// AddCount 变动积分
func AddCount(args *ArgsAddCount) (err error) {
	var data FieldsIntegral
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, count FROM user_integral WHERE org_id = $1 AND user_id = $2", args.OrgID, args.UserID)
	if err == nil {
		data.Count += args.AddCount
		if data.Count < 1 {
			data.Count = 0
		}
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE user_integral SET update_at = NOW(), count = :count WHERE id = :id", map[string]interface{}{
			"id":    data.ID,
			"count": data.Count,
		})
	} else {
		if args.AddCount < 1 {
			err = errors.New("count less 1")
			return
		}
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_integral (org_id, user_id, count) VALUES (:org_id, :user_id, :count)", map[string]interface{}{
			"org_id":  args.OrgID,
			"user_id": args.UserID,
			"count":   args.AddCount,
		})
	}
	if err != nil {
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO user_integral_log (org_id, user_id, add_count, des) VALUES (:org_id,:user_id,:add_count,:des)", args)
	if err != nil {
		return
	}
	if args.OrgID > 0 {
		OrgUserMod.PushUpdateUserData(args.OrgID, args.UserID)
	}
	return
}

// ArgsClearUser 清空某用户的积分参数
type ArgsClearUser struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// ClearUser 清空某用户的积分
// 彻底清理用户的相关记录
// 注意，日志数据会保留，以确保可追溯性
func ClearUser(args *ArgsClearUser) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_integral", "user_id = :user_id", args)
	return
}

// ArgsClearOrg 清空某组织的积分参数
type ArgsClearOrg struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// ClearOrg 清空某组织的积分
func ClearOrg(args *ArgsClearOrg) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_integral", "org_id = :org_id", args)
	if err == nil {
		_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "user_integral_log", "org_id = :org_id", args)
	}
	return
}
