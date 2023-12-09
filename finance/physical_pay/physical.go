package FinancePhysicalPay

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetPhysicalList 获取实物列表参数
type ArgsGetPhysicalList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetPhysicalList 获取实物列表
func GetPhysicalList(args *ArgsGetPhysicalList) (dataList []FieldsPhysical, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "finance_physical_pay_physical"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, bind_from, need_count, limit_count, take_count, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetPhysicalID 获取指定实物参数
type ArgsGetPhysicalID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetPhysicalID 获取指定实物
func GetPhysicalID(args *ArgsGetPhysicalID) (data FieldsPhysical, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, bind_from, need_count, limit_count, take_count, params FROM finance_physical_pay_physical WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetPhysicalMore 获取多个实物配置参数
type ArgsGetPhysicalMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetPhysicalMoreNames 获取多个实物配置
func GetPhysicalMoreNames(args *ArgsGetPhysicalMore) (dataList map[int64]string, err error) {
	dataList, err = CoreSQLIDs.GetIDsOrgNameAndDelete("finance_physical_pay_physical", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// ArgsGetPhysicalByFrom 获取指定来源的数据参数
type ArgsGetPhysicalByFrom struct {
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//可用的来源标的物
	BindFrom CoreSQLFrom.FieldsFrom `db:"bind_from" json:"bindFrom"`
}

// GetPhysicalByFrom 获取指定来源的数据
func GetPhysicalByFrom(args *ArgsGetPhysicalByFrom) (data FieldsPhysical, err error) {
	var bindFrom string
	bindFrom, err = args.BindFrom.GetRawNoName()
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, bind_from, need_count, limit_count, take_count, params FROM finance_physical_pay_physical WHERE org_id = $1 AND bind_from @> $2 AND delete_at < to_timestamp(1000000)", args.OrgID, bindFrom)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// CheckPhysicalByFrom 检查标的物是否定义过
func CheckPhysicalByFrom(args *ArgsGetPhysicalByFrom) (err error) {
	var bindFrom string
	bindFrom, err = args.BindFrom.GetRawNoName()
	if err != nil {
		return
	}
	var data FieldsPhysical
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM finance_physical_pay_physical WHERE org_id = $1 AND bind_from @> $2 AND delete_at < to_timestamp(1000000)", args.OrgID, bindFrom)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreatePhysical 创建新的实物参数
type ArgsCreatePhysical struct {
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//可用的来源标的物
	BindFrom CoreSQLFrom.FieldsFrom `db:"bind_from" json:"bindFrom"`
	//置换一件商品需对应几个标的物
	NeedCount int64 `db:"need_count" json:"needCount" check:"int64Than0"`
	//标的物市场总投放量限制
	LimitCount int64 `db:"limit_count" json:"limitCount" check:"int64Than0"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreatePhysical 创建新的实物
func CreatePhysical(args *ArgsCreatePhysical) (data FieldsPhysical, err error) {
	var bindFrom string
	bindFrom, err = args.BindFrom.GetRawNoName()
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM finance_physical_pay_physical WHERE org_id = $1 AND bind_from @> $2 AND delete_at < to_timestamp(1000000)", args.OrgID, bindFrom)
	if err == nil && data.ID > 0 {
		err = errors.New("have data")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_physical_pay_physical", "INSERT INTO finance_physical_pay_physical (org_id, name, bind_from, need_count, limit_count, take_count, params) VALUES (:org_id,:name,:bind_from,:need_count,:limit_count,0,:params)", args, &data)
	return
}

// ArgsUpdatePhysical 修改实物参数
type ArgsUpdatePhysical struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//置换一件商品需对应几个标的物
	NeedCount int64 `db:"need_count" json:"needCount" check:"int64Than0"`
	//标的物市场总投放量限制
	LimitCount int64 `db:"limit_count" json:"limitCount" check:"int64Than0"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdatePhysical 修改实物
func UpdatePhysical(args *ArgsUpdatePhysical) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_physical_pay_physical SET update_at = NOW(), name = :name, need_count = :need_count, limit_count = :limit_count, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeletePhysical 删除实物参数
type ArgsDeletePhysical struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 用于检查
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeletePhysical 删除实物
func DeletePhysical(args *ArgsDeletePhysical) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_physical_pay_physical", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
