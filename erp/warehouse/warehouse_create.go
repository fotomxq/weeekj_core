package ERPWarehouse

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsCreateWarehouse 创建仓库参数
type ArgsCreateWarehouse struct {
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

// CreateWarehouse 创建仓库
func CreateWarehouse(args *ArgsCreateWarehouse) (err error) {
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_warehouse_warehouse (org_id, name, weight, size_w, size_h, size_z, address_data) VALUES (:org_id, :name, :weight, :size_w, :size_h, :size_z, :address_data)", args)
	if err != nil {
		return
	}
	//反馈
	return
}
