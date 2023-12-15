package UserSubscription

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/lib/pq"
)

// ArgsGetConfigList 获取配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "user_sub_config"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, mark, time_type, time_n, currency, price, price_old, title, des, cover_file_id, des_files, user_groups, exemption_discount, exemption_price, exemption_min_price, limits, exemption_time, style_id, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetConfigByID 获取指定配置ID参数
type ArgsGetConfigByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfigByID 获取指定配置ID
func GetConfigByID(args *ArgsGetConfigByID) (data FieldsConfig, err error) {
	data, err = getConfigByID(args.ID)
	if err != nil {
		err = errors.New(fmt.Sprint("no data, id: ", args.ID, ", org id: ", args.OrgID))
		return
	}
	if CoreFilter.EqID2(args.OrgID, data.OrgID) {
		return
	}
	err = errors.New(fmt.Sprint("no data, id: ", args.ID, ", org id: ", args.OrgID, ", real org id: ", data.OrgID))
	return
}

// GetConfigOnlyOne 获取全局唯一一个会员配置
func GetConfigOnlyOne() (data FieldsConfig) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_sub_config WHERE delete_at < to_timestamp(1000000) ORDER BY id ASC LIMIT 1")
	if err != nil || data.ID < 1 {
		return
	}
	data, err = getConfigByID(data.ID)
	if err != nil {
		return
	}
	return
}

// 获取配置
func getConfigByID(id int64) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, mark, time_type, time_n, currency, price, price_old, title, des, cover_file_id, des_files, user_groups, exemption_discount, exemption_price, exemption_min_price, limits, exemption_time, style_id, params FROM user_sub_config WHERE id = $1", id)
	return
}

// ArgsGetConfigMore 获取一组配置参数
type ArgsGetConfigMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetConfigMore 获取一组配置
func GetConfigMore(args *ArgsGetConfigMore) (dataList []FieldsConfig, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "user_sub_config", "id, create_at, update_at, delete_at, org_id, mark, time_type, time_n, currency, price, price_old, title, des, cover_file_id, des_files, user_groups, exemption_discount, exemption_price, exemption_min_price, limits, exemption_time, style_id, params", args.IDs, args.HaveRemove)
	return
}

// GetConfigMoreMap 获取一组配置名称组
func GetConfigMoreMap(args *ArgsGetConfigMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsTitleAndDelete("user_sub_config", args.IDs, args.HaveRemove)
	return
}

// GetConfigName 获取配置名称
func GetConfigName(id int64) (name string) {
	if id < 1 {
		return ""
	}
	_ = Router2SystemConfig.MainDB.Get(&name, "SELECT title FROM user_sub_config WHERE id = $1", id)
	return
}

// 检查用户组是否属于商户？
func checkUserGroup(orgID int64, groupIDs []int64) (err error) {
	if len(groupIDs) < 1 {
		return
	}
	err = UserCore.CheckGroupIsOrg(orgID, groupIDs)
	return
}
