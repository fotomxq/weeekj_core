package AnalysisAny

import "time"

//FieldsConfig 统计配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//数据标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//推送MQTT时间
	LastMQTT time.Time `db:"last_mqtt" json:"lastMQTT"`
	//上次推送的hash
	LastHash string `db:"last_hash" json:"lastHash"`
	//归档天数
	// 必须指定，小于1则强制按照3天计算
	FileDay int `db:"file_day" json:"fileDay" check:"intThan0"`
	//是否需要推送组织MQTT
	MqttOrg  bool `db:"mqtt_org" json:"mqttOrg"`
	MqttUser bool `db:"mqtt_user" json:"mqttUser"`
	MqttBind bool `db:"mqtt_bind" json:"mqttBind"`
}
