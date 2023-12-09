package ERPProduct

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceCompany "gitee.com/weeekj/weeekj_core/v5/service/company"
)

// ArgsCreateProductCompany 创建产品和供应商绑定关系参数
type ArgsCreateProductCompany struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//供应商公司
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//单价成本（不含税）
	CostPrice int64 `db:"cost_price" json:"costPrice" check:"price" empty:"true"`
	//税率
	// 实际税率=tax/100000
	Tax int64 `db:"tax" json:"tax"`
	//单价成本（含税）
	TaxCostPrice int64 `db:"tax_cost_price" json:"taxCostPrice" check:"price" empty:"true"`
	//返利设计
	RebatePrice FieldsProductRebateList `db:"rebate_price" json:"rebatePrice"`
	//建议零售价（不含税）
	// 供货商建议
	TipPrice int64 `db:"tip_price" json:"tipPrice" check:"price" empty:"true"`
	//建议零售价（含税）
	// 加入税收后的指导价
	TipTaxPrice int64 `db:"tip_tax_price" json:"tipTaxPrice" check:"price" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateProductCompany 创建产品和供应商绑定关系
func CreateProductCompany(args *ArgsCreateProductCompany) (err error) {
	//修正参数
	if len(args.RebatePrice) < 1 {
		args.RebatePrice = FieldsProductRebateList{}
	}
	if len(args.Params) < 1 {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//检查组织和商品关系
	productData := getProductByID(args.ProductID)
	if productData.ID < 1 || CoreSQL.CheckTimeHaveData(productData.DeleteAt) || !CoreFilter.EqID2(args.OrgID, productData.OrgID) {
		err = errors.New("product not exist")
		return
	}
	//检查供应商
	var companyData ServiceCompany.FieldsCompany
	companyData, err = ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
		ID:    args.CompanyID,
		OrgID: args.OrgID,
	})
	if err != nil || companyData.ID < 1 || CoreSQL.CheckTimeHaveData(companyData.DeleteAt) {
		err = errors.New("company not exist")
		return
	}
	//检查必须是唯一的关系，不能重复创建关系
	var findID int64
	err = Router2SystemConfig.MainDB.DB.Get(&findID, "SELECT id FROM erp_product_company WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND product_id = $2 AND company_id = $3 LIMIT 1", args.OrgID, args.ProductID, args.CompanyID)
	if err == nil && findID > 0 {
		err = errors.New("have bind")
		return
	}
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_product_company (org_id, product_id, company_id, currency, cost_price, tax, tax_cost_price, rebate_price, tip_price, tip_tax_price, params) VALUES (:org_id, :product_id, :company_id, :currency, :cost_price, :tax, :tax_cost_price, :rebate_price, :tip_price, :tip_tax_price, :params)", args)
	if err != nil {
		return
	}
	//反馈
	return
}

// 增加历史归档数据
func createProductCompanyHistory(args *ArgsCreateProductCompany) (err error) {
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_product_company_history (org_id, product_id, company_id, currency, cost_price, tax, tax_cost_price, rebate_price, tip_price, tip_tax_price, params) VALUES (:org_id, :product_id, :company_id, :currency, :cost_price, :tax, :tax_cost_price, :rebate_price, :tip_price, :tip_tax_price, :params)", args)
	if err != nil {
		return
	}
	//反馈
	return
}
