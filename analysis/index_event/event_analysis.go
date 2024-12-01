package AnalysisIndexEvent

import "fmt"

type DataGetEventLevelCount struct {
	//预警等级
	// 根据项目需求划定等级
	// 例如：0 低风险; 1 中风险; 2 高风险
	Level int `db:"level" json:"level"`
	//数量
	Count int64 `db:"count" json:"count"`
}

// GetEventLevelCount 获取风险等级统计
// 该方法会无视维度，获取所有风险等级的数量统计
func GetEventLevelCount() (dataList []DataGetEventLevelCount) {
	_ = eventDB.GetClient().DB.GetPostgresql().Select(&dataList, "SELECT level,count(*) as count FROM "+eventDB.GetClient().TableName+" GROUP BY level")
	return
}

// ArgsGetEventLevelCountByExtend 获取指定维度的风险事件数量关系参数
type ArgsGetEventLevelCountByExtend struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" empty:"true" index:"true" field_list:"true"`
	//扩展维度1
	// 可建议特别的维度关系，例如特定供应商的数据、特定地区的数据等
	Extend1 string `db:"extend1" json:"extend1" index:"true" field_list:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true" field_list:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true" field_list:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true" field_list:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true" field_list:"true"`
}

// DataGetEventLevelCountByExtend 获取指定维度的风险事件数量关系数据
type DataGetEventLevelCountByExtend struct {
	//预警等级
	// 根据项目需求划定等级
	// 例如：0 低风险; 1 中风险; 2 高风险
	Level int `db:"level" json:"level"`
	//数量
	Count int64 `db:"count" json:"count"`
}

// GetEventLevelCountByExtend 获取指定维度的风险事件数量关系
// 该方法需指定具体维度，以获取指定维度的风险等级数量统计
func GetEventLevelCountByExtend(args *ArgsGetEventLevelCountByExtend) (dataList []DataGetEventLevelCountByExtend) {
	_ = eventDB.GetClient().DB.GetPostgresql().Select(&dataList, fmt.Sprintf("SELECT level, count(*) as count FROM %s WHERE code = $1 AND extend1 = $2 AND extend2 = $3 AND extend3 = $4 AND extend4 = $5 AND extend5 = $6 GROUP BY level", eventDB.GetClient().TableName), args.Code, args.Extend1, args.Extend2, args.Extend3, args.Extend4, args.Extend5)
	return
}

type DataArgsDataGetEventLevelCountRanking struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//预警等级
	// 根据项目需求划定等级
	// 例如：0 低风险; 1 中风险; 2 高风险
	Level int `db:"level" json:"level"`
	//数量
	Count int64 `db:"count" json:"count"`
}

// GetEventLevelCountRanking 获取指标的数量排名
func GetEventLevelCountRanking(args *ArgsGetEventLevelCountByExtend) (dataList []DataArgsDataGetEventLevelCountRanking) {
	_ = eventDB.GetClient().DB.GetPostgresql().Select(&dataList, fmt.Sprintf("SELECT code, level, count(*) as count FROM %s GROUP BY code, level ORDER BY count DESC", eventDB.GetClient().TableName))
	return
}

// GetEventExtendDistinctList 获取指定维度的所有可选值
func GetEventExtendDistinctList(extendNum int) (dataList []string, err error) {
	//获取数据
	dataList, err = eventDB.GetList().GetDistinctList(fmt.Sprintf("extend%d", extendNum))
	if err != nil {
		return
	}
	//反馈
	return
}
