package ERPWarehouse

import "time"

// FieldsLocation 货位信息
type FieldsLocation struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域ID
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//货位编号
	Code string `db:"code" json:"code" check:"des" min:"1" max:"300"`
}
