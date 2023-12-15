package OrgTip

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/lib/pq"
	"time"
)

// FieldsTipType 提醒数据列队
type FieldsTipType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//推送目标系统
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//计划提醒时间
	TipAt time.Time `db:"tip_at" json:"tipAt"`
	//提醒标题
	Title string `db:"title" json:"title"`
	//提醒内容
	Content string `db:"content" json:"content"`
	//附件
	Files pq.Int64Array `db:"files" json:"files"`
	//是否需要短信
	NeedSMS bool `db:"need_sms" json:"needSMS"`
	//短信配置
	SMSConfigID int64 `db:"sms_config_id" json:"smsConfigID"`
	//短信模版参数
	SMSParams CoreSQLConfig.FieldsConfigsType `db:"sms_params" json:"smsParams"`
	//扩展数据
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//是否已经发送
	AllowSend bool `db:"allow_send" json:"allowSend"`
}
