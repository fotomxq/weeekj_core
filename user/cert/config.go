package UserCert

import Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"

// 获取配置列表
type ArgsGetConfigList struct {
}

func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	return
}

// 获取指定配置
type ArgsGetConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//绑定组织
	// 该组织根据资源来源设定
	// 如果是平台资源，则为0
	OrgID int64 `db:"org_id" json:"orgID" empty:"true"`
}

func GetConfig(args *ArgsGetConfig) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_cert_config WHERE id = $1 AND org_id = $2", args)
	return
}

// 创建配置
type ArgsCreateConfig struct {
}

func CreateConfig(args *ArgsCreateConfig) (data FieldsConfig, err error) {
	return
}

// 修改配置
type ArgsUpdateConfig struct {
}

func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	return
}

// 删除配置
type ArgsDeleteConfig struct {
}

func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	return
}
