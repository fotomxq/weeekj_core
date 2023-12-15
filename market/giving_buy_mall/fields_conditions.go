package MarketGivingBuyMall

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsConditions 赠送条件配置
type FieldsConditions struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//赠礼配置ID
	ConfigID int64 `db:"config_id" json:"configID"`
	//商品ID
	MallProductID int64 `db:"mall_product_id" json:"mallProductID"`
	//商品分类
	SortID int64 `db:"sort_id" json:"sortID"`
	//商品标签
	Tag int64 `db:"tag" json:"tag"`
	//订单的最小金额
	MinPrice int64 `db:"min_price" json:"minPrice"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
