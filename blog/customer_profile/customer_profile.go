package BlogCustomerProfile

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetList 获取列表参数
type ArgsGetList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetList 获取列表
func GetList(args *ArgsGetList) (dataList []FieldsProfile, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR msg ILIKE '%' || :search || '%' OR address ->> 'name' ILIKE '%' || :search || '%' OR address ->> 'address' ILIKE '%' || :search || '%' OR address ->> 'phone' ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "blog_customer_profile"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, delete_at, org_id, name, address, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "delete_at"},
	)
	return
}

// ArgsGetByID 获取ID参数
type ArgsGetByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetByID 获取ID
func GetByID(args *ArgsGetByID) (data FieldsProfile, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, delete_at, org_id, name, address, msg, params FROM blog_customer_profile WHERE org_id = $1 AND id = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.ID)
	return
}

// ArgsCreate 创建新记录参数
type ArgsCreate struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//客户姓名
	Name string `db:"name" json:"name" check:"name"`
	//客户联系地址组件
	Address CoreSQLAddress.FieldsAddress `db:"address" json:"address" check:"address_data" empty:"true"`
	//留言信息
	Msg string `db:"msg" json:"msg" check:"des" min:"1" max:"1000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// Create 创建新记录
func Create(args *ArgsCreate) (data FieldsProfile, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "blog_customer_profile", "INSERT INTO blog_customer_profile (org_id, name, address, msg, params) VALUES (:org_id, :name, :address, :msg, :params)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsDelete 删除记录参数
type ArgsDelete struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// Delete 删除记录
func Delete(args *ArgsDelete) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "blog_customer_profile", "id = :id AND org_id = :org_id", args)
	return
}
