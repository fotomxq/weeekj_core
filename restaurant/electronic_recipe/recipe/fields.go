package RestaurantElectronicRecipe

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsRecipe struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//菜品名称
	Name string `db:"name" json:"name"`
	//组织ID
	RawOrgID int64 `db:"raw_org_id" json:"rawOrgID"`
	//分公司ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//门店ID
	StoreID int64 `db:"store_id" json:"storeID"`
	// 建议售价
	Price int64 `db:"price" json:"price"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
