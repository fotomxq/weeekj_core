package BaseWeixinPayProtocol

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetTemplateList 获取列表参数
type ArgsGetTemplateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTemplateList 获取列表
func GetTemplateList(args *ArgsGetTemplateList) (dataList []FieldsTemplate, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
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
	if where != "" {
		where = "true"
	}
	tableName := "core_weixin_pay_protocol_template"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, code, name FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsCreateTemplate 创建模版参数
type ArgsCreateTemplate struct {
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模版在微信的编号
	// 该编号在组织下唯一
	Code string `db:"code" json:"code" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
}

// CreateTemplate 创建模版
func CreateTemplate(args *ArgsCreateTemplate) (data FieldsTemplate, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_weixin_pay_protocol_template", "INSERT INTO core_weixin_pay_protocol_template (org_id, code, name) VALUES (:org_id,:code,:name)", args, &data)
	return
}

// ArgsUpdateTemplate 更新模版参数
type ArgsUpdateTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模版在微信的编号
	// 该编号在组织下唯一
	Code string `db:"code" json:"code" check:"mark"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
}

// UpdateTemplate 更新模版
func UpdateTemplate(args *ArgsUpdateTemplate) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_weixin_pay_protocol_template SET code = :code, name = :name WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteTemplate 删除模版参数
type ArgsDeleteTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 设备所属的组织，也可能为0
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteTemplate 删除模版
func DeleteTemplate(args *ArgsDeleteTemplate) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "core_weixin_pay_protocol_template", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
