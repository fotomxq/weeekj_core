package ERPWarehouse

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateWarehouse 创建仓库参数
type ArgsUpdateWarehouse struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//仓库名称
	Name string `db:"name" json:"name" check:"name"`
	//承载重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//存储尺寸
	SizeW int `db:"size_w" json:"sizeW" check:"intThan0" empty:"true"`
	SizeH int `db:"size_h" json:"sizeH" check:"intThan0" empty:"true"`
	SizeZ int `db:"size_z" json:"sizeZ" check:"intThan0" empty:"true"`
	//地址信息
	AddressData CoreSQLAddress.FieldsAddress `db:"address_data" json:"addressData"`
}

// UpdateWarehouse 创建仓库
func UpdateWarehouse(args *ArgsUpdateWarehouse) (err error) {
	//修改数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_warehouse_warehouse SET update_at = NOW(), name = :name, weight = :weight, size_w = :size_w, size_h = :size_h, size_z = :size_z, address_data = :address_data WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//删除缓冲
	deleteWarehouseCache(args.ID)
	//反馈
	return
}
