package FinanceReturnedMoney

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

type ArgsDeleteCompany struct {
	//ID
	CompanyID int64 `db:"company_id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func DeleteCompany(args *ArgsDeleteCompany) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "finance_returned_money_company", "company_id = :company_id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteCompanyCache(args.CompanyID)
	return
}
