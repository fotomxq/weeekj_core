package AnalysisIndexEvent

// ArgsInsertEvent 插入新的预警事件参数
type ArgsInsertEvent struct {
	//指标编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//年月日
	// 可任意持续，如年，或仅年月
	// 不建议构建小时及以下级别的指标
	// 同一个维度和时间范围，仅会存在一个数据，否则将覆盖
	YearMD string `db:"year_md" json:"yearMD" index:"true" field_list:"true"`
	//预警等级
	// 根据项目需求划定等级
	// 例如：0 低风险; 1 中风险; 2 高风险
	Level int `db:"level" json:"level" index:"true" field_list:"true"`
	//来源指标值的系统和ID
	// 避免重复触发预警
	FromSystem string `db:"from_system" json:"fromSystem" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	FromID     int64  `db:"from_id" json:"fromID" check:"id" index:"true" field_list:"true"`
	//触发类型
	// 根据项目需求划定类型，可以留空
	FromType string `db:"from_type" json:"fromType" check:"des" min:"1" max:"50" index:"true" field_list:"true"`
	//扩展维度1
	// 可建立特别的维度关系，例如特定供应商的数据、特定地区的数据等
	// 该维度主要用于筛选数据
	Extend1 string `db:"extend1" json:"extend1" index:"true" field_list:"true"`
	//扩展维度2
	Extend2 string `db:"extend2" json:"extend2" index:"true" field_list:"true"`
	//扩展维度3
	Extend3 string `db:"extend3" json:"extend3" index:"true" field_list:"true"`
	//扩展维度4
	Extend4 string `db:"extend4" json:"extend4" index:"true" field_list:"true"`
	//扩展维度5
	Extend5 string `db:"extend5" json:"extend5" index:"true" field_list:"true"`
	//指标预警阈值，触发预警时的值
	Threshold int64 `db:"threshold" json:"threshold" index:"true"`
	//触发值
	Val float64 `db:"val" json:"val" index:"true"`
}

// InsertEvent 插入新的预警事件
func InsertEvent(args *ArgsInsertEvent) (err error) {
	//检查指标是否已触发过风险
	_, err = getEventBySystem(args.FromSystem, args.FromID, args.FromType)
	if err != nil {
		err = nil
		return
	}
	//写入数据
	_, err = eventDB.GetInsert().InsertRow(args)
	if err != nil {
		return
	}
	//反馈
	return
}
