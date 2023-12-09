package ERPProduct

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetProductCompanyList 获取产品供应商列表参数
type ArgsGetProductCompanyList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//供应商公司
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
}

// GetProductCompanyList 获取产品供应商列表
func GetProductCompanyList(args *ArgsGetProductCompanyList) (dataList []FieldsProductCompany, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ProductID > -1 {
		where = where + " AND product_id = :product_id"
		maps["product_id"] = args.ProductID
	}
	if args.CompanyID > -1 {
		where = where + " AND company_id = :company_id"
		maps["company_id"] = args.CompanyID
	}
	tableName := "erp_product_company"
	var rawList []FieldsProductCompany
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "weight", "tip_price", "tip_tax_price"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getProductCompanyByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetProductCompany 获取供货商配置数据
func GetProductCompany(orgID int64, productID int64, companyID int64) (data FieldsProductCompany) {
	var b bool
	data, b = CheckProductCompany(orgID, productID, companyID)
	if !b {
		data = FieldsProductCompany{}
		return
	}
	return
}

// CheckProductCompany 检查供应商是否供应商品，并反馈数据
func CheckProductCompany(orgID int64, productID int64, companyID int64) (data FieldsProductCompany, b bool) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM erp_product_company WHERE org_id = $1 AND product_id = $2 AND company_id = $3 AND delete_at < to_timestamp(1000000)", orgID, productID, companyID)
	if err != nil {
		return
	}
	data = getProductCompanyByID(data.ID)
	b = true
	return
}

// GetProductCompanyLast 获取产品的历史信息
func GetProductCompanyLast(orgID int64, productID int64, companyID int64, beforeAt time.Time) (data FieldsProductCompany) {
	//如果实际数据更新时间较早，则直接获取当前数据即可
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM erp_product_company WHERE org_id = $1 AND product_id = $2 AND company_id = $3 AND delete_at < to_timestamp(1000000) AND update_at <= $4", orgID, productID, companyID, beforeAt)
	if err == nil && data.ID > 0 {
		data = getProductCompanyByID(data.ID)
		return
	}
	return
}

// GetProductCompanyRand 随机抽取一个供货商
// randMode 随机模式，0 只要反馈即可; 1 成本价最低
func GetProductCompanyRand(productID int64, randMode int) (data FieldsProductCompany) {
	dataList := getProductCompanyListByProductID(productID)
	switch randMode {
	case 0:
		if len(dataList) > 0 {
			data = dataList[0]
		}
	case 1:
		for _, v := range dataList {
			if data.ID < 0 {
				v = data
				continue
			}
			if data.TaxCostPrice > v.TaxCostPrice {
				data = v
				continue
			}
		}
	}
	if data.ID < 1 {
		if len(dataList) > 0 {
			data = dataList[0]
		}
	}
	return
}

// getProductCompanyListByProductID 通过产品获取关联的供应商
func getProductCompanyListByProductID(productID int64) (dataList []FieldsProductCompany) {
	var rawList []FieldsProductCompany
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM erp_product_company WHERE product_id = $1 AND delete_at < to_timestamp(1000000)", productID)
	if err != nil || len(rawList) < 1 {
		return
	}
	for _, v := range rawList {
		vData := getProductCompanyByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// getProductCompanyByID 获取产品关联供应商
func getProductCompanyByID(id int64) (data FieldsProductCompany) {
	cacheMark := getProductCompanyCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, product_id, company_id, currency, cost_price, tax, tax_cost_price, rebate_price, tip_price, tip_tax_price, params FROM erp_product_company WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheProductCompanyTime)
	return
}
