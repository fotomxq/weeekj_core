package MarketGivingCore

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
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
		where = where + "(name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"market_giving_core_config",
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, market_config_id, limit_time_type, limit_count, user_integral, user_subs, user_tickets, deposit_config_mark, price, count, params FROM market_giving_core_config WHERE "+where,
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
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, market_config_id, limit_time_type, limit_count, user_integral, user_subs, user_tickets, deposit_config_mark, price, count, params FROM market_giving_core_config WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsGetConfigMore 获取一组配置参数
type ArgsGetConfigMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetConfigMore 获取一组配置
func GetConfigMore(args *ArgsGetConfigMore) (dataList []FieldsConfig, err error) {
	err = CoreSQLIDs.GetIDsOrgAndDelete(&dataList, "market_giving_core_config", "id, create_at, update_at, delete_at, org_id, name, market_config_id, limit_time_type, limit_count, user_integral, user_subs, user_tickets, deposit_config_mark, price, count, params", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// GetConfigMoreMap 获取一组配置名称组
func GetConfigMoreMap(args *ArgsGetConfigMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgNameAndDelete("market_giving_core_config", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// ArgsCreateConfig 创建新的配置参数
type ArgsCreateConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"title" min:"1" max:"1000" empty:"true"`
	//推荐后奖励配置
	MarketConfigID int64 `db:"market_config_id" json:"marketConfigID" check:"id" empty:"true"`
	//领取周期类型
	// 0 不限制; 1 一次性; 2 每天限制; 3 每周限制; 4 每月限制; 5 每季度限制; 6 每年限制
	LimitTimeType int `db:"limit_time_type" json:"limitTimeType" check:"intThan0" empty:"true"`
	//领取次数
	LimitCount int `db:"limit_count" json:"limitCount" check:"intThan0" empty:"true"`
	//奖励积分
	UserIntegral int64 `db:"user_integral" json:"userIntegral" check:"int64Than0" empty:"true"`
	//奖励用户订阅
	UserSubs FieldsConfigUserSubs `db:"user_subs" json:"userSubs"`
	//奖励票据
	UserTickets FieldsConfigUserTickets `db:"user_tickets" json:"userTickets"`
	//奖励金储蓄标识码
	DepositConfigMark string `db:"deposit_config_mark" json:"depositConfigMark" check:"mark" empty:"true"`
	//奖励金额
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//奖励次数
	Count int64 `db:"count" json:"count" check:"int64Than0" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateConfig 创建新的配置
func CreateConfig(args *ArgsCreateConfig) (data FieldsConfig, err error) {
	if err = checkUserSub(args.OrgID, args.UserSubs); err != nil {
		return
	}
	if err = checkUserTicket(args.OrgID, args.UserTickets); err != nil {
		return
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "market_giving_core_config", "INSERT INTO market_giving_core_config (org_id, name, market_config_id, limit_time_type, limit_count, user_integral, user_subs, user_tickets, deposit_config_mark, price, count, params) VALUES (:org_id,:name,:market_config_id,:limit_time_type,:limit_count,:user_integral,:user_subs,:user_tickets,:deposit_config_mark,:price,:count,:params)", args, &data)
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"title" min:"1" max:"1000" empty:"true"`
	//推荐后奖励配置
	MarketConfigID int64 `db:"market_config_id" json:"marketConfigID" check:"id" empty:"true"`
	//领取周期类型
	// 0 不限制; 1 一次性; 2 每天限制; 3 每周限制; 4 每月限制; 5 每季度限制; 6 每年限制
	LimitTimeType int `db:"limit_time_type" json:"limitTimeType" check:"intThan0" empty:"true"`
	//领取次数
	LimitCount int `db:"limit_count" json:"limitCount" check:"intThan0" empty:"true"`
	//奖励积分
	UserIntegral int64 `db:"user_integral" json:"userIntegral" check:"int64Than0" empty:"true"`
	//奖励用户订阅
	UserSubs FieldsConfigUserSubs `db:"user_subs" json:"userSubs"`
	//奖励票据
	UserTickets FieldsConfigUserTickets `db:"user_tickets" json:"userTickets"`
	//奖励金储蓄标识码
	DepositConfigMark string `db:"deposit_config_mark" json:"depositConfigMark" check:"mark" empty:"true"`
	//奖励金额
	Price int64 `db:"price" json:"price" check:"price" empty:"true"`
	//奖励次数
	Count int64 `db:"count" json:"count" check:"int64Than0" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	if err = checkUserSub(args.OrgID, args.UserSubs); err != nil {
		return
	}
	if err = checkUserTicket(args.OrgID, args.UserTickets); err != nil {
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE market_giving_core_config SET update_at = NOW(), name = :name, market_config_id = :market_config_id, limit_time_type = :limit_time_type, limit_count = :limit_count, user_integral = :user_integral, user_subs = :user_subs, user_tickets = :user_tickets, deposit_config_mark = :deposit_config_mark, price = :price, count = :count, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "market_giving_core_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

func checkUserSub(orgID int64, args FieldsConfigUserSubs) (err error) {
	for _, v := range args {
		if v.ConfigID < 1 {
			continue
			//err = errors.New("sub config id less 1")
			//return
		}
		if v.Count < 1 || v.CountTime < 1 {
			err = errors.New("count or count time less 1")
			return
		}
		/**
		var vConfig UserSubscription.FieldsConfig
		vConfig, err = UserSubscription.GetConfigByID(&UserSubscription.ArgsGetConfigByID{
			ID:    v.ConfigID,
			OrgID: orgID,
		})
		if err != nil {
			return
		}
		if vConfig.ID < 1 {
			err = errors.New("sub config not exist")
			return
		}
		*/
	}
	return
}

func checkUserTicket(orgID int64, args FieldsConfigUserTickets) (err error) {
	for _, v := range args {
		if v.ConfigID < 1 {
			continue
			//err = errors.New("ticket config id less 1")
			//return
		}
		if v.Count < 1 {
			err = errors.New("count or count time less 1")
			return
		}
		/**
		var vConfig UserTicket.FieldsConfig
		vConfig, err = UserTicket.GetConfigByID(&UserTicket.ArgsGetConfigByID{
			ID:    v.ConfigID,
			OrgID: orgID,
		})
		if err != nil {
			return
		}
		if vConfig.ID < 1 {
			err = errors.New("ticket config not exist")
			return
		}
		*/
	}
	return
}
