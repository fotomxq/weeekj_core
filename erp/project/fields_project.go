package ERPProject

import "time"

type FieldsProject struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"50"`
	//描述
	Desc string `db:"desc" json:"desc" check:"des" min:"1" max:"300" empty:"true"`
	//关联预算
	BudgetID int64 `db:"budget_id" json:"budgetID" check:"id"`
}
