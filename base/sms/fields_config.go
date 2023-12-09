package BaseSMS

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsConfig 短信模版和配置信息结构
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
	//使用哪一家？
	// tencent / aliyun
	System string `db:"system" json:"system"`
	//来源系统的显示名称
	Name string `db:"name" json:"name"`
	//应用ID
	AppID string `db:"app_id" json:"appID"`
	//应用密钥
	AppKey string `db:"app_key" json:"appKey"`
	//默认过期时间
	DefaultExpire string `db:"default_expire" json:"defaultExpire"`
	//获取间隔时间 秒
	TimeSpacing int64 `db:"time_spacing" json:"timeSpacing"`
	//模版ID
	TemplateID string `db:"template_id" json:"templateID"`
	//签名名称
	TemplateSign string `db:"template_sign" json:"templateSign"`
	//扩展参数
	TemplateParams CoreSQLConfig.FieldsConfigsType `db:"template_params" json:"templateParams"`
	//扩展参数
	// mark: 'check' 验证码短信
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
