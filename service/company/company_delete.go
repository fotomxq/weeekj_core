package ServiceCompany

import (
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteCompany 删除公司参数
type ArgsDeleteCompany struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteCompany 删除公司
func DeleteCompany(args *ArgsDeleteCompany) (err error) {
	//获取公司信息
	data := getCompany(args.ID)
	//删除公司信息
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "service_company", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//请求删除绑定信息
	deleteBindByCompanyID(data.OrgID, data.ID)
	//清理缓冲
	deleteCompanyCache(data.ID)
	//推送通知
	CoreNats.PushDataNoErr("/service/company", "delete", data.ID, "", nil)
	//反馈
	return
}
