package MarketGivingUserSub

import (
	"errors"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	UserSubscriptionMod "github.com/fotomxq/weeekj_core/v5/user/subscription/mod"
	"github.com/lib/pq"
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
	tableName := "market_giving_user_sub"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, config_id, sub_config_id, sub_buy_count, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetConditionsMoreMap 获取一组配置参数
type ArgsGetConditionsMoreMap struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetConditionsMoreMap 获取一组配置名称组
func GetConditionsMoreMap(args *ArgsGetConditionsMoreMap) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("market_giving_user_sub", args.IDs, args.HaveRemove)
	return
}

// ArgsCheckConfigAndOrg 检查配置和商户是否关联参数
type ArgsCheckConfigAndOrg struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// CheckConfigAndOrg 检查配置和商户是否关联
func CheckConfigAndOrg(args *ArgsCheckConfigAndOrg) (err error) {
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM market_giving_user_sub WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	if err == nil && id < 1 {
		err = errors.New("not exist")
		return
	}
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
	//订阅ID
	SubConfigID int64 `db:"sub_config_id" json:"subConfigID" check:"id"`
	//订阅单位
	SubBuyCount int64 `db:"sub_buy_count" json:"subBuyCount" check:"int64Than0"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateConditions 创建条件
func CreateConditions(args *ArgsCreateConditions) (data FieldsConditions, err error) {
	//检查订阅
	err = UserSubscriptionMod.CheckConfigAndOrg(&UserSubscriptionMod.ArgsCheckConfigAndOrg{
		ID:    args.SubConfigID,
		OrgID: args.OrgID,
	})
	if err != nil {
		err = errors.New("user sub config not exist")
		return
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "market_giving_user_sub", "INSERT INTO market_giving_user_sub (org_id, name, config_id, sub_config_id, sub_buy_count, params) VALUES (:org_id,:name,:config_id,:sub_config_id,:sub_buy_count,:params)", args, &data)
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
	//订阅ID
	SubConfigID int64 `db:"sub_config_id" json:"subConfigID" check:"id"`
	//订阅单位
	SubBuyCount int64 `db:"sub_buy_count" json:"subBuyCount" check:"int64Than0"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConditions 修改条件
func UpdateConditions(args *ArgsUpdateConditions) (err error) {
	//检查订阅
	err = UserSubscriptionMod.CheckConfigAndOrg(&UserSubscriptionMod.ArgsCheckConfigAndOrg{
		ID:    args.SubConfigID,
		OrgID: args.OrgID,
	})
	if err != nil {
		err = errors.New("user sub config not exist")
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE market_giving_user_sub SET update_at = NOW(), name = :name, config_id = :config_id, sub_config_id = :sub_config_id, sub_buy_count = :sub_buy_count, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
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
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "market_giving_user_sub", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
