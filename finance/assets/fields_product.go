package FinanceAssets

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsProduct 组织产品数据类型
type FieldsProduct struct {
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
	//描述
	Des string `db:"des" json:"des"`
	//封面列
	// 第一张作为封面
	CoverFiles pq.Int64Array `db:"cover_files" json:"coverFiles"`
	//描述信息
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//产品编码
	Code string `db:"code" json:"code"`
	//储蓄货币类型
	// 采用CoreCurrency匹配
	Currency int `db:"currency" json:"currency"`
	//产品价值
	Price int64 `db:"price" json:"price"`
	//关联的仓储产品
	// 用于仓储识别本资产时使用
	WarehouseProductIDs pq.Int64Array `db:"warehouse_product_ids" json:"warehouseProductIDs"`
	//关联的商品数据
	// 用于商品ID识别本资产时使用
	MallCommodityIDs pq.Int64Array `db:"mall_commodity_ids" json:"mallCommodityIDs"`
	//总数统计
	Count int64 `db:"count" json:"count"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
