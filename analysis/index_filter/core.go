package AnalysisIndexFilter

import BaseSQLTools "github.com/fotomxq/weeekj_core/v5/base/sql_tools"

//规则预设筛选
/**
用途：
1. 用于特定算法模型规则，对指定数据范围进行筛选
2. 可标记筛选后的特定维度+来源，用于后续分析/业务使用

使用方法：
1. 使用Clear清理掉数据
2. 使用Append添加数据
3. 使用GetAll获取所有数据，注意如果数据量预计特别庞大，建议使用GetList获取数据列表遍历
*/

var (
	filterDB BaseSQLTools.Quick
)

func Init() (err error) {
	//规则预设筛选
	if err = filterDB.Init("analysis_index_filter", &FieldsFilter{}); err != nil {
		return
	}
	return
}
