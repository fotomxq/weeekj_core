package ERPProduct

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsProductCompanyHistory 产品和供应商关联模块历史
type FieldsProductCompanyHistory struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
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
	//单价税额
	TaxCostPrice int64 `db:"tax_cost_price" json:"taxCostPrice" check:"price" empty:"true"`
	//返利设计
	RebatePrice FieldsProductRebateList `db:"rebate_price" json:"rebatePrice" empty:"true"`
	//建议零售价（不含税）
	// 供货商建议
	TipPrice int64 `db:"tip_price" json:"tipPrice" check:"price" empty:"true"`
	//建议零售价（含税）
	// 加入税收后的指导价
	TipTaxPrice int64 `db:"tip_tax_price" json:"tipTaxPrice" check:"price" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
