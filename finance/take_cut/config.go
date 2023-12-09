package FinanceTakeCut

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetConfigList 获取配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//组织分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
}

// GetConfigList 获取配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.SortID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if where == "" {
		where = "true"
	}
	tableName := "finance_take_cut_config"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, sort_id, org_id, order_system, cut_price_proportion FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetOrgConfig 获取商户的佣金设置参数
type ArgsGetOrgConfig struct {
	//组织
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetOrgConfig 获取商户的佣金设置
func GetOrgConfig(args *ArgsGetOrgConfig) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, sort_id, org_id, order_system, cut_price_proportion FROM finance_take_cut_config WHERE org_id = $1", args.OrgID)
	if err != nil || data.ID < 1 {
		var orgData OrgCoreCore.FieldsOrg
		orgData, err = OrgCoreCore.GetOrg(&OrgCoreCore.ArgsGetOrg{
			ID: args.OrgID,
		})
		if err != nil {
			return
		}
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, sort_id, org_id, order_system, cut_price_proportion FROM finance_take_cut_config WHERE sort_id = $1", orgData.SortID)
		if err == nil && data.ID < 1 {
			err = errors.New("no data")
			return
		}
	}
	return
}

// ArgsSetConfig 设置配置参数
type ArgsSetConfig struct {
	//组织分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//组织ID
	// 商户和分类必须指定其中一个，会优先采用商户进行，否则采用分类
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//针对订单的系统来源
	// eg: user_sub / org_sub / mall
	OrderSystem string `db:"order_system" json:"orderSystem" check:"mark"`
	//抽成比例
	// 单位10万分之%，例如300000对应3%；3对应0.00003%
	CutPriceProportion int64 `db:"cut_price_proportion" json:"cutPriceProportion"`
}

// SetConfig 设置配置
func SetConfig(args *ArgsSetConfig) (err error) {
	if args.SortID > 0 && args.OrgID > 0 {
		err = errors.New("sort and org same")
		return
	}
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM finance_take_cut_config WHERE sort_id = $1 AND org_id = $2", args.SortID, args.OrgID)
	if err == nil && id > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_take_cut_config SET order_system = :order_system, cut_price_proportion = :cut_price_proportion WHERE id = :id", map[string]interface{}{
			"id":                   id,
			"order_system":         args.OrderSystem,
			"cut_price_proportion": args.CutPriceProportion,
		})
		return
	}
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_take_cut_config (sort_id, org_id, order_system, cut_price_proportion) VALUES (:sort_id,:org_id,:order_system,:cut_price_proportion)", args)
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteConfig 删除配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "finance_take_cut_config", "id", args)
	return
}
