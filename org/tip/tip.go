package OrgTip

import (
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"

	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"github.com/lib/pq"
)

//行政专用的提醒模块
// 外部行政类模块可以任意使用，将数据写入该列队
// 列队将不断巡逻，如果发现需要提醒的数据，将自动发送给目标来源系统

// ArgsCreate 创建推送计划参数
type ArgsCreate struct {
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
	//附加文件
	Files pq.Int64Array `db:"files" json:"files"`
	//是否需要短信
	NeedSMS bool `db:"need_sms" json:"needSMS"`
	//短信配置
	SMSConfigID int64 `db:"sms_config_id" json:"smsConfigID"`
	//短信模版参数
	SMSParams CoreSQLConfig.FieldsConfigsType `db:"sms_params" json:"smsParams"`
	//扩展数据
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Create 创建推送计划
func Create(args *ArgsCreate) (err error) {
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO org_tip (org_id, create_info, from_info, tip_at, title, content, files, need_sms, sms_config_id, sms_params, params, allow_send) VALUES (:org_id, :create_info, :from_info, :tip_at, :title, :content, :files, :need_sms, :sms_config_id, :sms_params, :params, false)", args)
	return
}

// ArgsDeleteByCreateFrom 删除某个来源的所有推送参数
type ArgsDeleteByCreateFrom struct {
	//创建来源
	CreateInfo CoreSQLFrom.FieldsFrom
}

// DeleteByCreateFrom 删除某个来源的所有推送
func DeleteByCreateFrom(args *ArgsDeleteByCreateFrom) (err error) {
	var createData string
	createData, err = args.CreateInfo.GetRaw()
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_tip SET delete_at = NOW() WHERE create_info @> :create_info", map[string]interface{}{
		"create_info": createData,
	})
	return
}
