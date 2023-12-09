package FinancePayMod

import (
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
)

// FixOrgID 修正组织ID
func FixOrgID(orgID int64) (newOrgID int64) {
	//如果存在组织ID
	if orgID > 0 {
		//检查全局是否打开了强制平台支付处理
		financePayOtherInOne, _ := BaseConfig.GetDataBool("FinancePayOtherInOne")
		if financePayOtherInOne {
			orgID = 0
		} else {
			//检查商户是否开通的独立支付体系
			openFinanceIndependent := OrgCoreCore.CheckOrgPermissionFunc(orgID, "finance_independent")
			openAllFunc := OrgCoreCore.CheckOrgPermissionFunc(orgID, "all")
			if !openFinanceIndependent && !openAllFunc {
				orgID = 0
			}
		}
	} else {
		orgID = 0
	}
	newOrgID = orgID
	return
}
