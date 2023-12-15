package ERPWarehouse

import (
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	"time"
)

// FieldsWarehouse 仓库
type FieldsWarehouse struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//仓库名称
	Name string `db:"name" json:"name"`
	//承载重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//存储尺寸
	SizeW int `db:"size_w" json:"sizeW"`
	SizeH int `db:"size_h" json:"sizeH"`
	SizeZ int `db:"size_z" json:"sizeZ"`
	//地址信息
	AddressData CoreSQLAddress.FieldsAddress `db:"address_data" json:"addressData"`
}
