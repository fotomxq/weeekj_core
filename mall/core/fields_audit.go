package MallCore

import (
	"github.com/lib/pq"
	"time"
)

//FieldsAudit 审核
type FieldsAudit struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//组织ID
	// 留空则表明为平台的用户留下的内容
	OrgID int64 `db:"org_id" json:"orgID"`
	//商品ID
	ProductID int64 `db:"product_id" json:"productID"`
	//审核通过时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//拒绝原因
	BanDes string `db:"ban_des" json:"banDes"`
	//拒绝附加文件
	BanDesFiles pq.Int64Array `db:"ban_des_files" json:"banDesFiles"`
}