package BaseEmail

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsEmailServerType 发送渠道
type FieldsEmailServerType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//显示的名称
	Name string `db:"name" json:"name"`
	//host
	Host string `db:"host" json:"host"`
	//端口
	Port string `db:"port" json:"port"`
	//是否启用SSL
	IsSSL bool `db:"is_ssl" json:"isSSL"`
	//发件方邮件地址和用户名
	Email string `db:"email" json:"email"`
	//密码
	Password string `db:"password" json:"password"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
