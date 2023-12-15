package BaseSMS

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetConfigList 获取列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (system ILIKE '%' || :search || '%' OR name ILIKE '%' || :search || '%' OR template_id ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "core_sms_config"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, system, name, app_id, app_key, default_expire, time_spacing, template_id, template_sign, template_params, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetConfigByID 获取某个数据参数
type ArgsGetConfigByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfigByID 获取某个数据
func GetConfigByID(args *ArgsGetConfigByID) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, system, name, app_id, app_key, default_expire, time_spacing, template_id, template_sign, template_params, params FROM core_sms_config WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	return
}

// ArgsCreateConfig 创建新配置参数
type ArgsCreateConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//使用哪一家？
	// tencent / aliyun
	System string `db:"system" json:"system" check:"mark"`
	//来源系统的显示名称
	Name string `db:"name" json:"name" check:"name"`
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
	//默认参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateConfig 创建新配置
func CreateConfig(args *ArgsCreateConfig) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_sms_config WHERE system = $1 AND template_id = $2 AND delete_at < to_timestamp(1000000)", args.System, args.TemplateID)
	if err == nil {
		err = errors.New("system and template id is exist")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_sms_config", "INSERT INTO core_sms_config(org_id, system, app_id, app_key, name, default_expire, time_spacing, template_id, template_sign, template_params, params) VALUES(:org_id, :system, :app_id, :app_key, :name, :default_expire, :time_spacing, :template_id, :template_sign, :template_params, :params)", args, &data)
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//使用哪一家？
	// tencent / aliyun
	System string `db:"system" json:"system" check:"mark"`
	//来源系统的显示名称
	Name string `db:"name" json:"name" check:"name"`
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
	//默认参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	var data FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM core_sms_config WHERE system = $1 AND template_id = $2 AND delete_at < to_timestamp(1000000)", args.System, args.TemplateID)
	if err == nil && data.ID != args.ID {
		err = errors.New("system and template id is exist")
		return
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE core_sms_config SET update_at = NOW(), system = :system, name = :name, app_id = :app_id, app_key = :app_key, default_expire = :default_expire, time_spacing = :time_spacing, template_id = :template_id, template_sign = :template_sign, template_params = :template_params, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteConfig 删除配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "core_sms_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
