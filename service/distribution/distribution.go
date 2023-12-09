package ServiceDistribution

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetDistributionList 获取分销商列表参数
type ArgsGetDistributionList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetDistributionList 获取分销商列表
func GetDistributionList(args *ArgsGetDistributionList) (dataList []FieldsDistribution, dataCount int64, err error) {
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
	tableName := "service_distribution_distribution"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, user_id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsCreateDistribution 添加新的分销商参数
type ArgsCreateDistribution struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分销商姓名
	Name string `db:"name" json:"name" check:"name"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// CreateDistribution 添加新的分销商
func CreateDistribution(args *ArgsCreateDistribution) (data FieldsDistribution, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_distribution_distribution", "INSERT INTO service_distribution_distribution (org_id, name, user_id) VALUES (:org_id,:name,:user_id)", args, &data)
	return
}

// ArgsUpdateDistribution 修改分销商参数
type ArgsUpdateDistribution struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分销商姓名
	Name string `db:"name" json:"name" check:"name"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// UpdateDistribution 修改分销商
func UpdateDistribution(args *ArgsUpdateDistribution) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_distribution_distribution SET update_at = NOW(), name = :name, user_id = :user_id WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteDistribution 删除分销商参数
type ArgsDeleteDistribution struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteDistribution 删除分销商
func DeleteDistribution(args *ArgsDeleteDistribution) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_distribution_distribution", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
