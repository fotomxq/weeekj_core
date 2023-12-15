package AnalysisAny

import (
	"errors"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"

	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
)

// ArgsInitConfig 初始化配置设置参数
type ArgsInitConfig struct {
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//归档天数
	// 必须指定，小于1则强制按照3天计算
	FileDay int `db:"file_day" json:"fileDay" check:"intThan0"`
	//是否需要推送组织MQTT
	MqttOrg  bool `db:"mqtt_org" json:"mqttOrg"`
	MqttUser bool `db:"mqtt_user" json:"mqttUser"`
	MqttBind bool `db:"mqtt_bind" json:"mqttBind"`
}

// InitConfig 初始化配置设置
// Deprecated
func InitConfig(args *ArgsInitConfig) (err error) {
	//检查mark是否已经存在
	var configID int64
	err = Router2SystemConfig.MainDB.Get(&configID, "SELECT id FROM analysis_any_config WHERE mark = $1 AND delete_at < to_timestamp(1000000) LIMIT 1", args.Mark)
	if err == nil && configID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE analysis_any_config SET update_at = NOW(), file_day = :file_day, mqtt_org = :mqtt_org, mqtt_user = :mqtt_user, mqtt_bind = :mqtt_bind WHERE id = :id", map[string]interface{}{
			"id":        configID,
			"file_day":  args.FileDay,
			"mqtt_org":  args.MqttOrg,
			"mqtt_user": args.MqttUser,
			"mqtt_bind": args.MqttBind,
		})
		return
	}
	//不存在创建新的配置
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO analysis_any_config (mark, file_day, mqtt_org, mqtt_user, mqtt_bind) VALUES (:mark,:file_day,:mqtt_org, :mqtt_user, :mqtt_bind)", args)
	return
}

// getConfigMark 获取指定的配置
func getConfigMark(mark string) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, mark, file_day, mqtt_org, mqtt_user, mqtt_bind FROM analysis_any_config WHERE mark = $1 AND delete_at < to_timestamp(1000000) LIMIT 1", mark)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}
