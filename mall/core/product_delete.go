package MallCore

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsDeleteProduct 删除商品参数
type ArgsDeleteProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteProduct 删除商品
func DeleteProduct(args *ArgsDeleteProduct) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "mall_core", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteProductCache(args.ID)
	//反馈
	return
}

// ArgsReturnProduct 恢复商品参数
type ArgsReturnProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// ReturnProduct 恢复商品
func ReturnProduct(args *ArgsReturnProduct) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE mall_core SET update_at = NOW(), delete_at = to_timestamp(0), publish_at = to_timestamp(0) WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteProductCache(args.ID)
	//反馈
	return
}
