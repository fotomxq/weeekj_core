package FinanceDeposit

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisPrice 储蓄账户金额参数
type ArgsGetAnalysisPrice struct {
	//存储来源
	// 如不同的加盟商，储蓄均有所差异，也可以留空不指定，则为平台总的储蓄资金池
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//储蓄配置标识码
	// 可用于同一类货币下，多个用途，如赠送的储值额度、或用户自行充值的额度
	// user 用户自己储值 ; deposit 押金 ; free 免费赠送额度 ; ... 特定系统下的充值模块
	ConfigMark string `db:"config_mark" json:"configMark"`
}

// GetAnalysisPrice 储蓄账户金额
func GetAnalysisPrice(args *ArgsGetAnalysisPrice) (count int64, err error) {
	where := "delete_at < TO_TIMESTAMP(1000000) AND config_mark = :config_mark"
	maps := map[string]interface{}{
		"config_mark": args.ConfigMark,
	}
	var newWhere string
	newWhere, maps, err = args.FromInfo.GetList("from_info", "from_info", maps)
	if err != nil {
		return
	} else {
		if newWhere != "" {
			where = where + " AND " + newWhere
		}
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "save_price", "price", where, maps)
	return
}
