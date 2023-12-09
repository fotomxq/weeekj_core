package OrgDomain

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetDomainList 获取域名列表参数
type ArgsGetDomainList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetDomainList 获取域名列表
func GetDomainList(args *ArgsGetDomainList) (dataList []FieldsDomain, dataCount int64, err error) {
	where := "(:org_id < 1 OR org_id = :org_id) AND (host ILIKE '%' || :search || '%')"
	tableName := "org_domain"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, host, params FROM "+tableName+" WHERE "+where,
		where,
		map[string]interface{}{
			"org_id": args.OrgID,
			"search": args.Search,
		},
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetDomainOrg 通过host获取org参数
type ArgsGetDomainOrg struct {
	//Host
	// 全局唯一
	Host string `db:"host" json:"host"`
}

// GetDomainOrg 通过host获取org
func GetDomainOrg(args *ArgsGetDomainOrg) (orgID int64, params CoreSQLConfig.FieldsConfigsType, err error) {
	var data FieldsDomain
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT org_id, params FROM org_domain WHERE host = $1", args.Host)
	if err == nil && data.OrgID < 1 {
		err = errors.New("no org")
		return
	}
	orgID = data.OrgID
	params = data.Params
	return
}

// ArgsCreateDomain 创建新的域名参数
type ArgsCreateDomain struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//Host
	// 全局唯一
	Host string `db:"host" json:"host"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateDomain 创建新的域名
func CreateDomain(args *ArgsCreateDomain) (data FieldsDomain, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_domain WHERE host = $1", args.Host)
	if err == nil || data.ID > 0 {
		err = errors.New("host is replace")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_domain", "INSERT INTO org_domain (org_id, host, params) VALUES (:org_id,:host,:params)", args, &data)
	return
}

// ArgsUpdateDomain 修改域名参数
type ArgsUpdateDomain struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//Host
	// 全局唯一
	Host string `db:"host" json:"host"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateDomain 修改域名
func UpdateDomain(args *ArgsUpdateDomain) (err error) {
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_domain WHERE host = $1", args.Host)
	if err == nil && id != args.ID {
		err = errors.New("host is replace")
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_domain SET host = :host, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteDomain 删除域名参数
type ArgsDeleteDomain struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteDomain 删除域名
func DeleteDomain(args *ArgsDeleteDomain) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "org_domain", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
