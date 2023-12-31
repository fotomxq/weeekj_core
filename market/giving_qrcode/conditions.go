package MarketGivingQrcode

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetConditionsList 获取条件列表参数
type ArgsGetConditionsList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//关联的奖励
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConditionsList 获取条件列表
func GetConditionsList(args *ArgsGetConditionsList) (dataList []FieldsConditions, dataCount int64, err error) {
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
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "market_giving_qrcode"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, config_id, have_phone, have_order, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsCreateConditions 创建条件参数
type ArgsCreateConditions struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"title" min:"1" max:"300"`
	//赠礼配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//是否需要绑定手机号
	HavePhone bool `db:"have_phone" json:"havePhone" check:"bool"`
	//是否需要发生交易
	HaveOrder bool `db:"have_order" json:"haveOrder" check:"bool"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateConditions 创建条件
func CreateConditions(args *ArgsCreateConditions) (data FieldsConditions, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "market_giving_qrcode", "INSERT INTO market_giving_qrcode (org_id, name, config_id, have_phone, have_order, params) VALUES (:org_id,:name,:config_id,:have_phone,:have_order,:params)", args, &data)
	return
}

// ArgsUpdateConditions 修改条件参数
type ArgsUpdateConditions struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"title" min:"1" max:"300"`
	//赠礼配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//是否需要绑定手机号
	HavePhone bool `db:"have_phone" json:"havePhone" check:"bool"`
	//是否需要发生交易
	HaveOrder bool `db:"have_order" json:"haveOrder" check:"bool"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConditions 修改条件
func UpdateConditions(args *ArgsUpdateConditions) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE market_giving_qrcode SET update_at = NOW(), name = :name, config_id = :config_id, have_phone = :have_phone, have_order = :have_order, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteConditions 删除条件参数
type ArgsDeleteConditions struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteConditions 删除条件
func DeleteConditions(args *ArgsDeleteConditions) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "market_giving_qrcode", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
