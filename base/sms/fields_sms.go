package BaseSMS

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsSMS 短信结构体
type FieldsSMS struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//过期时间
	// 过期后不清理数据，但存在保留的最大时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//发送时间
	SendAt time.Time `db:"send_at" json:"sendAt"`
	//失败原因
	// 如果为本地原因则显示错误代码，否则显示API提供方反馈信息
	FailedMsg string `db:"failed_msg" json:"failedMsg"`
	//是否已经验证
	// 仅用于验证码处理
	IsCheck bool `db:"is_check" json:"isCheck"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//会话
	Token int64 `db:"token" json:"token"`
	//国家代码
	NationCode string `db:"nation_code" json:"nationCode"`
	//目标手机号
	// 目标手机号是唯一的标识码
	Phone string `db:"phone" json:"phone"`
	//附带参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//创建来源和创建来源ID
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
}
