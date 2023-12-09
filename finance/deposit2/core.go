package FinanceDeposit2

import (
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	"sync"
)

//第二代储蓄模块
/**
1. 每种储蓄模块将独立表记录，有特定的行为特征处理方案
2. 不支持多配置方案处理，所有特殊项目交给配置项记录处理
*/

var (
	orgDepositLock  sync.Mutex
	orgSavingLock   sync.Mutex
	userDepositLock sync.Mutex
	userFreeLock    sync.Mutex
	userSavingLock  sync.Mutex
	//OpenAnalysis 是否启动analysis
	OpenAnalysis = false
)

func Init() {
	if OpenAnalysis {
		AnalysisAny2.SetConfigBeforeNoErr("finance_deposit_org_add_price", 1, 180)
		AnalysisAny2.SetConfigBeforeNoErr("finance_deposit_total_org_deposit_price", 1, 180)
		AnalysisAny2.SetConfigBeforeNoErr("finance_deposit_total_org_saving_price", 1, 180)
		AnalysisAny2.SetConfigBeforeNoErr("finance_deposit_total_user_deposit_price", 1, 180)
		AnalysisAny2.SetConfigBeforeNoErr("finance_deposit_total_user_free_price", 1, 180)
		AnalysisAny2.SetConfigBeforeNoErr("finance_deposit_total_user_saving_price", 1, 180)
	}
}
