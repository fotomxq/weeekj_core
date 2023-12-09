package ERPPermanentAssets

// FieldsSortAnalysis 固定资产清查统计
type FieldsSortAnalysis struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//年月
	YearMonth int `db:"year_month" json:"yearMonth"`
	//更新时间
	UpdateAt int64 `db:"update_at" json:"updateAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID"`
	//购进价格
	AllBuyPrice int64 `db:"all_buy_price" json:"allBuyPrice"`
	//期初数量
	BeginCount int64 `db:"begin_count" json:"beginCount"`
	//期末数量
	EndCount int64 `db:"end_count" json:"endCount"`
	//期初余额
	BeginBalance int64 `db:"begin_balance" json:"beginBalance"`
	//期末余额
	EndBalance int64 `db:"end_balance" json:"endBalance"`
}
