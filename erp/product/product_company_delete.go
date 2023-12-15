package ERPProduct

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteProductCompany 删除产品和供应商关系
type ArgsDeleteProductCompany struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func DeleteProductCompany(args *ArgsDeleteProductCompany) (err error) {
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_product_company", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteProductCompanyCache(args.ID)
	//反馈
	return
}

// deleteProductCompanyByProductID 删除产品关联的所有公司
func deleteProductCompanyByProductID(productID int64) (err error) {
	dataList := getProductCompanyListByProductID(productID)
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_product_company", "product_id = :product_id", map[string]interface{}{
		"product_id": productID,
	})
	if err != nil {
		return
	}
	for _, v := range dataList {
		deleteProductCompanyCache(v.ID)
	}
	return
}
