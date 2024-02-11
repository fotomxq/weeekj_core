package UserTicket

import (
	"errors"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisUse2 统计使用情况参数
type ArgsGetAnalysisUse2 struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//模式
	// 0 无效; 1 赠送；2 使用
	Mode int `db:"mode" json:"mode" check:"intThan0" empty:"true"`
	//时间段
	TimeBetween CoreSQLTime2.DataCoreTime `json:"timeBetween"`
}

type DataAnalysisUse2 struct {
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//配置名称
	ConfigName string `db:"config_name" json:"configName"`
	//数量
	Count int64 `db:"count" json:"count"`
}

// GetAnalysisUse2 统计使用情况
func GetAnalysisUse2(args *ArgsGetAnalysisUse2) (dataList []DataAnalysisUse, err error) {
	//获取所有配置
	var configList []FieldsConfig
	err = Router2SystemConfig.MainDB.Select(&configList, "SELECT id, title FROM user_ticket_config WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", args.OrgID)
	if err != nil {
		return
	}
	if len(configList) < 1 {
		err = errors.New("no any config")
		return
	}
	//遍历配置，生成统计数据集合
	for _, vConfig := range configList {
		var vData DataAnalysisUse
		_ = Router2SystemConfig.MainDB.Get(&vData, "SELECT SUM(count) as count FROM user_ticket_log WHERE config_id = $1 AND create_at >= $2 AND create_at <= $3 AND mode = $4 LIMIT 1", vConfig.ID, args.TimeBetween.MinTime, args.TimeBetween.MaxTime, args.Mode)
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
