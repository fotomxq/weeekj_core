package MallCore

import "time"

//FieldsTransport 配送运费模版
type FieldsTransport struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//模版名称
	Name string `db:"name" json:"name"`
	//计费规则
	// 0 无配送费；1 按件计算；2 按重量计算；3 按公里数计算
	// 后续所有计费标准必须统一，否则系统将拒绝创建或修改
	Rules int `db:"rules" json:"rules"`
	//首费标准
	// 0 留空；1 N件开始收费；2 N重量开始收费；3 N公里开始收费
	RulesUnit int `db:"rules_unit" json:"rulesUnit"`
	//首费金额
	RulesPrice int64 `db:"rules_price" json:"rulesPrice"`
	//增费标准
	// 0 无增量；1 每N件增加费用；2 每N重量增加费用；3 每N公里增加费用
	AddType int `db:"add_type" json:"addType"`
	//增费单位
	AddUnit int `db:"add_unit" json:"addUnit"`
	//增费金额
	// 单位增加的费用
	AddPrice int64 `db:"add_price" json:"addPrice"`
	//免邮条件
	// 0 无免费；1 按件免费；2 按重量免费; 3 公里数内免费
	FreeType int `db:"free_type" json:"freeType"`
	//免邮单位
	FreeUnit int `db:"free_unit" json:"freeUnit"`
}
