package ERPProduct

import "time"

// FieldsProductVals 产品扩展数据表
type FieldsProductVals struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//采用模板ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id"`
	//顺序序号
	OrderNum int64 `db:"order_num" json:"orderNum"`
	//插槽值
	SlotID int64 `db:"slot_id" json:"slotID" check:"id"`
	//值(字符串)
	DataValue string `db:"data_value" json:"dataValue"`
	//值(浮点数)
	DataValueNum float64 `db:"data_value_num" json:"dataValueNum"`
	//值(整数)
	DataValueInt int64 `db:"data_value_int" json:"dataValueInt"`
	//参数
	Params string `db:"params" json:"params"`
}
