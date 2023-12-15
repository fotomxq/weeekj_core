package ServiceCompany

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetCompanyAuditList 获取公司审核列表参数
type ArgsGetCompanyAuditList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定组织ID
	BindOrgID int64 `db:"bind_org_id" json:"bindOrgID"`
	//绑定用户ID
	// 主绑定关系具备所有能力，类似组织的拥有人
	BindUserID int64 `db:"bind_user_id" json:"bindUserID"`
	//用途
	// client 客户; supplier 供应商; partners 合作商; service 服务商
	UseType string `db:"use_type" json:"useType"`
	//是否审核
	NeedIsAudit bool `json:"needIsAudit" check:"bool"`
	IsAudit     bool `json:"isAudit" check:"bool"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetCompanyAuditList 获取公司审核列表
func GetCompanyAuditList(args *ArgsGetCompanyAuditList) (dataList []FieldsCompanyAudit, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindOrgID > -1 {
		where = where + " AND bind_org_id = :bind_org_id"
		maps["bind_org_id"] = args.BindOrgID
	}
	if args.BindUserID > -1 {
		where = where + " AND bind_user_id = :bind_user_id"
		maps["bind_user_id"] = args.BindUserID
	}
	if args.UseType != "" {
		where = where + " AND use_type = :use_type"
		maps["use_type"] = args.UseType
	}
	if args.NeedIsAudit {
		where = CoreSQL.GetDeleteSQLField(args.IsAudit, where, "audit_at")
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%' OR sn ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsCompanyAudit
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"service_company_audit",
		"id",
		"SELECT id FROM service_company_audit WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "name"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	//遍历重组数据
	for _, v := range rawList {
		vData := getCompanyAudit(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}
func getCompanyAudit(id int64) (data FieldsCompanyAudit) {
	cacheMark := getCompanyAuditCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, audit_at, hash, org_id, use_type, bind_org_id, bind_user_id, name, sn, des, country, city, address, map_type, longitude, latitude, params FROM service_company_audit WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 604800)
	return
}
