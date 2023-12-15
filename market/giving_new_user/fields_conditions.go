package MarketGivingNewUser

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsConditions 赠送条件配置
type FieldsConditions struct {
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
	//名称
	Name string `db:"name" json:"name"`
	//赠礼配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//是否需要绑定手机号
	HavePhone bool `db:"have_phone" json:"havePhone"`
	//什么时候注册之后的
	AfterSign time.Time `db:"after_sign" json:"afterSign"`
	//什么时候注册之前的
	BeforeSign time.Time `db:"before_sign" json:"beforeSign"`
	//是否需要发生交易
	HaveOrder bool `db:"have_order" json:"haveOrder"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
