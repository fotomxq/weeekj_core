package ServiceCompany

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetCompanyList 获取公司列表参数
type ArgsGetCompanyList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用途
	// client 客户; supplier 供应商; partners 合作商; service 服务商
	UseType string `db:"use_type" json:"useType"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetCompanyList 获取公司列表
func GetCompanyList(args *ArgsGetCompanyList) (dataList []FieldsCompany, dataCount int64, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UseType != "" {
		where = where + " AND use_type = :use_type"
		maps["use_type"] = args.UseType
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%' OR sn ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsCompany
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"service_company",
		"id",
		"SELECT id, org_id FROM service_company WHERE "+where,
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
		vData := getCompany(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsSearchCompany 搜索公司信息参数
type ArgsSearchCompany struct {
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用途
	// client 客户; supplier 供应商; partners 合作商; service 服务商
	UseType string `db:"use_type" json:"useType"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataSearchCompany struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//名称
	Name string `db:"name" json:"name"`
}

// SearchCompany 搜索公司信息
func SearchCompany(args *ArgsSearchCompany) (dataList []DataSearchCompany, err error) {
	//获取数据
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(false, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UseType != "" {
		where = where + " AND use_type = :use_type"
		maps["use_type"] = args.UseType
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT id, name FROM service_company WHERE "+where+" LIMIT 10",
		maps,
	)
	if err != nil {
		return
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	//反馈
	return
}

// ArgsGetCompanyID 获取公司信息参数
type ArgsGetCompanyID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetCompanyID 获取公司信息
func GetCompanyID(args *ArgsGetCompanyID) (data FieldsCompany, err error) {
	data = getCompany(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsCompany{}
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetCompanyIDs 查询一组公司
type ArgsGetCompanyIDs struct {
	//一组ID
	IDs pq.Int64Array `db:"ids" json:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetCompanyNames 查询一组公司名字
func GetCompanyNames(args *ArgsGetCompanyIDs) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgNameAndDelete("service_company", args.IDs, args.OrgID, args.HaveRemove)
	return
}

func GetCompanyName(id int64) string {
	if id < 1 {
		return ""
	}
	data := getCompany(id)
	return data.Name
}

func GetCompanyUseType(id int64) string {
	data := getCompany(id)
	return data.UseType
}

// GetCompanyByHash 通过hash获取公司信息
func GetCompanyByHash(hash string, orgID int64) (data FieldsCompany, err error) {
	cacheMark := getCompanyCacheHashMark(orgID, hash)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err != nil || data.ID < 1 {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, hash, org_id, use_type, bind_org_id, bind_user_id, name, sn, des, country, city, address, map_type, longitude, latitude, params FROM service_company WHERE hash = $1 AND org_id = $2 LIMIT 1", hash, orgID)
		if err != nil {
			return
		}
	}
	if CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		data = FieldsCompany{}
		err = errors.New("no data")
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 604800)
	return
}

// GetCompanyBySN 通过公司SN获取公司信息
func GetCompanyBySN(orgID int64, sn string, useType string) (data FieldsCompany) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_company WHERE org_id = $1 AND sn = $2 AND use_type = $3 AND delete_at < to_timestamp(1000000)", orgID, sn, useType)
	if err != nil || data.ID < 1 {
		data = FieldsCompany{}
		return
	}
	data = getCompany(data.ID)
	return
}

// GetCompanyByName 通过公司名称获取公司信息
func GetCompanyByName(orgID int64, name string, useType string) (data FieldsCompany) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_company WHERE org_id = $1 AND name = $2 AND use_type = $3 AND delete_at < to_timestamp(1000000)", orgID, name, useType)
	if err != nil || data.ID < 1 {
		data = FieldsCompany{}
		return
	}
	data = getCompany(data.ID)
	return
}

func getCompany(id int64) (data FieldsCompany) {
	cacheMark := getCompanyCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, hash, org_id, use_type, bind_org_id, bind_user_id, name, sn, des, country, city, address, map_type, longitude, latitude, params FROM service_company WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 604800)
	return
}
