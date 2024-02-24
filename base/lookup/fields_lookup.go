package BaseLookup

import "time"

type FieldsLookup struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//是否为系统预设
	IsSys bool `db:"is_sys" json:"isSys" check:"bool"`
	//领域ID
	DomainID int64 `db:"domain_id" json:"domainID" check:"id" empty:"true"`
	//管理单元
	UnitID int64 `db:"unit_id" json:"unitID" check:"id" empty:"true"`
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"100"`
	//名称
	Name string `db:"name" json:"name" check:"des" min:"1" max:"300"`
}
