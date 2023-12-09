package OrgActivity

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsConfig 活动配置表
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//开始时间
	StartAt time.Time `db:"start_at" json:"startAt"`
	//结束时间
	EndAt time.Time `db:"end_at" json:"endAt"`
	//活动名称
	Name string `db:"name" json:"name"`
	//活动封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//商品描述
	Des string `db:"des" json:"des"`
	//描述图组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//关联的平台会员ID列
	// 允许商户使用这些内容
	UserSubIDs pq.Int64Array `db:"user_sub_ids" json:"userSubIDs"`
	//关联的平台票据ID列
	// 允许商户使用这些内容
	UserTicketIDs pq.Int64Array `db:"user_ticket_ids" json:"userTicketIDs"`
	//关联的储蓄设计
	// 具体依赖于特定代码实现
	FinanceDepositIDs FieldsConfigFinances `db:"finance_deposit_ids" json:"financeDepositIDs"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// FieldsConfigFinances 储蓄相关设计
type FieldsConfigFinances []FieldsConfigFinance

type FieldsConfigFinance struct {
	//储蓄配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//是否可用于消费
	UseBuy bool `db:"use_buy" json:"useBuy"`
	//是否可用于赠与
	UseGiving bool `db:"use_giving" json:"useGiving"`
}
