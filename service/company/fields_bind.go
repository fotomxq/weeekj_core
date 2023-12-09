package ServiceCompany

import (
	"github.com/lib/pq"
	"time"
)

// FieldsBind 公司和用户绑定关系
type FieldsBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID"`
	//可以预设手机号，手续用户绑定后自动绑定对应用户
	//绑定手机号的国家代码
	NationCode string `db:"nation_code" json:"nationCode"`
	//手机号码，绑定后的手机
	Phone string `db:"phone" json:"phone"`
	//赋予能力
	// 系统约定的几个特定能力，平台无法编辑该范围，只能授权
	Managers pq.StringArray `db:"managers" json:"managers"`
}
