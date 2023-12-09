package UserSubscription

import (
	"errors"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisUse 统计使用情况参数
type ArgsGetAnalysisUse struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

type DataAnalysisUse struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//配置名称
	ConfigName string `db:"config_name" json:"configName"`
	//数量
	Count int64 `db:"count" json:"count"`
}

// GetAnalysisUse 统计使用情况
func GetAnalysisUse(args *ArgsGetAnalysisUse) (dataList []DataAnalysisUse, err error) {
	//获取所有配置
	var configList []FieldsConfig
	if err = Router2SystemConfig.MainDB.Select(&configList, "SELECT id, title FROM user_sub_config WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", args.OrgID); err != nil {
		return
	}
	if len(configList) < 1 {
		err = errors.New("no any config")
		return
	}
	//获取时间段
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	//遍历配置，生成统计数据集合
	for _, vConfig := range configList {
		var vData DataAnalysisUse
		_ = Router2SystemConfig.MainDB.Get(&vData, "SELECT COUNT(id) as count FROM user_sub_log WHERE config_id = $1 AND create_at >= $2 AND create_at <= $3", vConfig.ID, timeBetween.MinTime, timeBetween.MaxTime)
		vData.ConfigID = vConfig.ID
		vData.ConfigName = vConfig.Title
		dataList = append(dataList, vData)
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetAnalysisArea 获取在续订阅分布情况参数
type ArgsGetAnalysisArea struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

type DataAnalysisArea struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//配置名称
	ConfigName string `db:"config_name" json:"configName"`
	//数量
	Count int64 `db:"count" json:"count"`
}

// GetAnalysisArea 获取在续订阅分布情况
func GetAnalysisArea(args *ArgsGetAnalysisArea) (dataList []DataAnalysisUse, err error) {
	//获取所有配置
	var configList []FieldsConfig
	if err = Router2SystemConfig.MainDB.Select(&configList, "SELECT id, title FROM user_sub_config WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", args.OrgID); err != nil {
		return
	}
	if len(configList) < 1 {
		err = errors.New("no any config")
		return
	}
	if err != nil {
		return
	}
	//遍历配置，生成统计数据集合
	for _, vConfig := range configList {
		var vData DataAnalysisUse
		_ = Router2SystemConfig.MainDB.Get(&vData, "SELECT COUNT(id) as count FROM user_sub WHERE config_id = $1 AND delete_at < to_timestamp(1000000) AND expire_at >= NOW()", vConfig.ID)
		vData.ConfigID = vConfig.ID
		vData.ConfigName = vConfig.Title
		dataList = append(dataList, vData)
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}
