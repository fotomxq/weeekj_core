package RouterFinance

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	FinanceDeposit "gitee.com/weeekj/weeekj_core/v5/finance/deposit"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
)

//储蓄相关组件

// DepositGetByUser 获取该用户的储蓄数据结构体
// 如果没有会自动给创建
func DepositGetByUser(userData *UserCore.DataUserDataType, fromInfo CoreSQLFrom.FieldsFrom, configMark string) (depositData FinanceDeposit.FieldsDepositType, errCode string, err error) {
	//查询存储
	depositData, err = FinanceDeposit.GetByFrom(&FinanceDeposit.ArgsGetByFrom{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     userData.Info.ID,
			Mark:   "",
			Name:   "",
		},
		FromInfo:   fromInfo,
		ConfigMark: configMark,
	})
	if err == nil {
		return
	}
	//不存在则创建
	depositData, errCode, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
		UpdateHash: "",
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     userData.Info.ID,
			Mark:   "",
			Name:   userData.Info.Name,
		},
		FromInfo:        fromInfo,
		ConfigMark:      configMark,
		AppendSavePrice: 0,
	})
	//反馈最终数据
	return
}
