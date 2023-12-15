package OrgCert

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsCert 证件记录
type FieldsCert struct {
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
	//审核时间
	AuditAt time.Time `db:"audit_at" json:"auditAt"`
	//审核人
	AuditBindID int64 `db:"audit_bind_id" json:"auditBindID"`
	//是否被审核拒绝
	AuditBanAt time.Time `db:"audit_ban_at" json:"auditBanAt"`
	//审核留言
	AuditDes string `db:"audit_des" json:"auditDes"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID"`
	//名称
	Name string `db:"name" json:"name"`
	//证件序列号
	SN string `db:"sn" json:"sn"`
	//拍照文件ID序列
	FileIDs pq.Int64Array `db:"file_ids" json:"fileIDs"`
	//是否已经缴费
	PayAt time.Time `db:"pay_at" json:"payAt"`
	//支付失败原因
	PayFailed bool `db:"pay_failed" json:"payFailed"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID"`
	//审核费用
	Currency int   `db:"currency" json:"currency"`
	Price    int64 `db:"price" json:"price"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
