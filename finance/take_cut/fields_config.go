package FinanceTakeCut

import "time"

//FieldsConfig 抽成约定配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//组织ID
	// 商户和分类必须指定其中一个，会优先采用商户进行，否则采用分类
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//针对订单的系统来源
	// eg: user_sub / org_sub / mall
	OrderSystem string `db:"order_system" json:"orderSystem" check:"mark"`
	//抽成比例
	// 单位10万分之%，例如300000对应3%；3对应0.00003%
	CutPriceProportion int64 `db:"cut_price_proportion" json:"cutPriceProportion"`
}
