package OrgCert2

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//默认过期时间长度
	DefaultExpire string `db:"default_expire" json:"defaultExpire"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//绑定来源
	// user 用户 / org 商户 / org_bind 商户成员 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark"`
	//证件名称
	Name string `db:"name" json:"name"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//审核模式
	// none 无需审核; wait 人工审核; auto 自动审核(依赖其他模块，根据扩展参数具体识别方案);
	AuditType string `db:"audit_type" json:"auditType"`
	//审核费用
	// 如果为0则无效
	Currency int   `db:"currency" json:"currency"`
	Price    int64 `db:"price" json:"price"`
	//序列号长度
	// 0则不会限制，但数据表最多存储100位，超出需使用扩展参数存储
	SNLen int `db:"sn_len" json:"snLen"`
	//通知类型
	// none 无通知; audit 审核通过后通知; expire 过期前通知; all 全部通知;
	TipType string `db:"tip_type" json:"tipType"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
