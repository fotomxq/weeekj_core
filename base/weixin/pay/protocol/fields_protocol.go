package BaseWeixinPayProtocol

import "time"

//FieldsProtocol 续约请求记录表
type FieldsProtocol struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 也是商户ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID"`
	//签约模版ID
	TemplateID int64 `db:"template_id" json:"templateID"`
	//签约模块
	// 0 会员模块; 1 平台组织会员模块
	ConfigSystem int `db:"config_system" json:"configSystem"`
	//签约模块ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//下一检查时间
	// 将在下一次到期之前检查会员到期情况，如果即将在24之后到期，将触发扣费请求
	NextAt time.Time `db:"next_at" json:"nextAt"`
}
