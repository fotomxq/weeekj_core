package TMSTransport

import "time"

type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//数据产生来源
	BindID int64 `db:"bind_id" json:"bindID"`
	//配送单ID
	TransportID int64 `db:"transport_id" json:"transportID"`
	//配送人员
	TransportBindID int64 `db:"transport_bind_id" json:"transportBindID"`
	//行为特征
	Mark string `db:"mark" json:"mark"`
	//备注
	Des string `db:"des" json:"des"`
}
