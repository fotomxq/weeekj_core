package ERPProduct

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateProductCompany 更新产品和供货商关系参数
type ArgsUpdateProductCompany struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//单价成本（不含税）
	CostPrice int64 `db:"cost_price" json:"costPrice" check:"price" empty:"true"`
	//税率
	// 实际税率=tax/100000
	Tax int64 `db:"tax" json:"tax"`
	//单价税额
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

// UpdateProductCompany 更新产品和供货商关系
func UpdateProductCompany(args *ArgsUpdateProductCompany) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_product_company SET update_at = NOW(), currency = :currency, cost_price = :cost_price, tax = :tax, tax_cost_price = :tax_cost_price, rebate_price = :rebate_price, tip_price = :tip_price, tip_tax_price = :tip_tax_price, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteProductCompanyCache(args.ID)
	//增加历史归档数据
	data := getProductCompanyByID(args.ID)
	err = createProductCompanyHistory(&ArgsCreateProductCompany{
		OrgID:        data.OrgID,
		ProductID:    data.ProductID,
		CompanyID:    data.CompanyID,
		Currency:     data.Currency,
		CostPrice:    data.CostPrice,
		Tax:          data.Tax,
		TaxCostPrice: data.TaxCostPrice,
		RebatePrice:  data.RebatePrice,
		TipPrice:     data.TipPrice,
		TipTaxPrice:  data.TipTaxPrice,
		Params:       data.Params,
	})
	if err != nil {
		CoreLog.Error("erp product update product company, create product company history, product id: ", data.ProductID, ", company id: ", data.CompanyID, ", err: ", err)
		err = nil
	}
	//反馈
	return
}
