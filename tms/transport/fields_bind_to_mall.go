package TMSTransport

import (
	"time"
)

//FieldsBindToMall 配送员和商品绑定关系
type FieldsBindToMall struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
	//绑定商品
	BindMallID int64 `db:"bind_mall_id" json:"bindMallID" check:"id"`
}
