package ERPProduct

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteProduct 删除产品参数
type ArgsDeleteProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteProduct 删除产品
func DeleteProduct(args *ArgsDeleteProduct) (err error) {
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_product", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteProductCache(args.ID)
	//删除关联的供应商关系
	_ = deleteProductCompanyByProductID(args.ID)
	//反馈
	return
}
