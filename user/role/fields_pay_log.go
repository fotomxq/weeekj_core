package UserRole

import "time"

//FieldsPayLog 付款记录
// 给角色人员付款储蓄的交易记录
type FieldsPayLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//角色ID
	RoleID int64 `db:"role_id" json:"roleID"`
	//储蓄配置
	DepositMark string `db:"deposit_mark" json:"depositMark"`
	//货币
	// eg: 86
	Currency int `db:"currency" json:"currency" check:"currency"`
	//支付金额
	Price int64 `db:"price" json:"price"`
	//系统来源
	SystemFrom string `db:"system_from" json:"systemFrom"`
	//系统ID
	FromID int64 `db:"from_id" json:"fromID"`
	//支付ID
	// 可能不存在，当平台向用户发起支付时，会直接给调账账户资金，而不是产生交易记录
	PayID int64 `db:"pay_id" json:"payID"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"300" empty:"true"`
}
