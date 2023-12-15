package RouterOrgFinance

import (
	"errors"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinanceDeposit "github.com/fotomxq/weeekj_core/v5/finance/deposit"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
)

// GetDepositDataAndDefaultMark 获取组织默认收款配置项
func GetDepositDataAndDefaultMark(orgID int64) (depositData FinanceDeposit.FieldsDepositType, defaultDepositMark string, err error) {
	//获取储蓄的默认配置项
	defaultDepositMark, err = OrgCore.Config.GetConfigVal(&ClassConfig.ArgsGetConfig{
		BindID:    orgID,
		Mark:      "finance_deposit_default_mark",
		VisitType: "admin",
	})
	if err != nil {
		err = errors.New("get org config default deposit mark, " + err.Error())
		return
	}
	//写入反馈数据集
	depositData, err = getDepositData(orgID, defaultDepositMark)
	if err != nil {
		err = errors.New("not have any deposit data")
		return
	}
	//反馈成功
	return
}

// getDepositData 获取组织的储蓄数据结构体
func getDepositData(orgID int64, saveMark string) (depositData FinanceDeposit.FieldsDepositType, err error) {
	//查询存储
	depositData, err = FinanceDeposit.GetByFrom(&FinanceDeposit.ArgsGetByFrom{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     orgID,
			Mark:   "",
			Name:   "",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     orgID,
			Mark:   "",
			Name:   "",
		},
		ConfigMark: saveMark,
	})
	if err == nil {
		return
	}
	//不存在则创建
	var orgData OrgCore.FieldsOrg
	orgData, err = OrgCore.GetOrg(&OrgCore.ArgsGetOrg{
		ID: orgID,
	})
	if err != nil {
		err = errors.New("org not exist, " + err.Error())
		return
	}
	depositData, _, err = FinanceDeposit.SetByFrom(&FinanceDeposit.ArgsSetByFrom{
		UpdateHash: "",
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     orgData.ID,
			Mark:   "",
			Name:   orgData.Name,
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "org",
			ID:     orgData.ID,
			Mark:   "",
			Name:   orgData.Name,
		},
		ConfigMark:      saveMark,
		AppendSavePrice: 0,
	})
	if err != nil {
		err = errors.New("set deposit failed, " + err.Error())
	}
	//反馈最终数据
	return
}
