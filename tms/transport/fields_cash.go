package TMSTransport

import "time"

//FieldsCash 配送员收付款金额
type FieldsCash struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//配送员ID
	BindID int64 `db:"bind_id" json:"bindID"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID"`
	//收付款类型
	// 0 收款(配送员收到款项) 1 付款(配送员付出款项)
	PayType int `db:"pay_type" json:"payType"`
	//金额
	Price int64 `db:"price" json:"price"`
}
