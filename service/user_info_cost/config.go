package ServiceUserInfoCost

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetConfigList 获取配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark" check:"mark" empty:"true"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark" check:"mark" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.RoomBindMark != "" {
		where = where + " AND room_bind_mark = :room_bind_mark"
		maps["room_bind_mark"] = args.RoomBindMark
	}
	if args.SensorMark != "" {
		where = where + " AND sensor_mark = :sensor_mark"
		maps["sensor_mark"] = args.SensorMark
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "service_user_info_cost_config"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, name, room_bind_mark, sensor_mark, count_type, each_unit, each_price, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetConfig 获取配置ID参数
type ArgsGetConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfig 获取配置ID
func GetConfig(args *ArgsGetConfig) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, room_bind_mark, sensor_mark, count_type, each_unit, each_price, params FROM service_user_info_cost_config WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsGetConfigs 获取更多配置参数
type ArgsGetConfigs struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetConfigs 获取更多配置
func GetConfigs(args *ArgsGetConfigs) (dataList []FieldsConfig, err error) {
	err = CoreSQLIDs.GetIDsOrgAndDelete(&dataList, "service_user_info_cost_config", "id, create_at, update_at, delete_at, org_id, name, room_bind_mark, sensor_mark, count_type, each_unit, each_price, params", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// GetConfigsName 获取多个配置名称
func GetConfigsName(args *ArgsGetConfigs) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgNameAndDelete("service_user_info_cost_config", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// ArgsCreateConfig 创建配置参数
type ArgsCreateConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark" check:"mark"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark" check:"mark"`
	//计算方式
	// 0 合并计算，将时间阶段内的所有遥感数据合并进行统计计算
	// 1 平均值计算，将时间段内的数据平均化计算
	CountType int `db:"count_type" json:"countType" check:"intThan0" empty:"true"`
	//每小时能耗值
	EachUnit float64 `db:"each_unit" json:"eachUnit"`
	//每小时费用
	// 每累计产生EachUnit，将增加该金额1次，不足将不增加继续等待累计
	Currency  int   `db:"currency" json:"currency" check:"currency"`
	EachPrice int64 `db:"each_price" json:"eachPrice" check:"price"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateConfig 创建配置
func CreateConfig(args *ArgsCreateConfig) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_user_info_cost_config WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND room_bind_mark = $2 AND sensor_mark = $3", args.OrgID, args.RoomBindMark, args.SensorMark)
	if err == nil && data.ID > 0 {
		err = errors.New("config is exist")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "service_user_info_cost_config", "INSERT INTO service_user_info_cost_config (org_id, name, room_bind_mark, sensor_mark, count_type, each_unit, currency, each_price, params) VALUES (:org_id, :name, :room_bind_mark, :sensor_mark, :count_type, :each_unit, :currency, :each_price, :params)", args, &data)
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name"`
	//房间场景值
	// 设备和房间绑定关系的mark值
	RoomBindMark string `db:"room_bind_mark" json:"roomBindMark" check:"mark"`
	//数据类型标识码
	// 遥感数据及传感器数据值
	SensorMark string `db:"sensor_mark" json:"sensorMark" check:"mark"`
	//计算方式
	// 0 合并计算，将时间阶段内的所有遥感数据合并进行统计计算
	// 1 平均值计算，将时间段内的数据平均化计算
	CountType int `db:"count_type" json:"countType" check:"intThan0" empty:"true"`
	//每小时能耗值
	EachUnit float64 `db:"each_unit" json:"eachUnit"`
	//每小时费用
	// 每累计产生EachUnit，将增加该金额1次，不足将不增加继续等待累计
	Currency  int   `db:"currency" json:"currency" check:"currency"`
	EachPrice int64 `db:"each_price" json:"eachPrice" check:"price"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	var data FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM service_user_info_cost_config WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND room_bind_mark = $2 AND sensor_mark = $3", args.OrgID, args.RoomBindMark, args.SensorMark)
	if err == nil && data.ID > 0 && data.ID != args.ID {
		err = errors.New("config is exist")
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_user_info_cost_config SET update_at = NOW(), name = :name, room_bind_mark = :room_bind_mark, sensor_mark = :sensor_mark, count_type = :count_type, each_unit = :each_unit, currency = :currency, each_price = :each_price, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteConfig 删除配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_user_info_cost_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
