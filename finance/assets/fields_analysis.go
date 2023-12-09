package FinanceAssets

import "time"

//统计用户持有和成员经手的变动记录
//TODO: 等待增加模块衔接

//FieldsAnalysisUser
//  以小时为单位记录
type FieldsAnalysisUser struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//统计周期
	// 完全相同的一个来源体系，1小时仅构建一条数据
	DayTime time.Time `db:"day_time" json:"dayTime"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//变动数量
	Count int64 `db:"count" json:"count"`
}

//FieldsAnalysisBind
//  以小时为单位记录
type FieldsAnalysisBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//实际操作人，组织绑定成员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//变动数量
	Count int64 `db:"count" json:"count"`
}