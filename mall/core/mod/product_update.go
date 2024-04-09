package MallCoreMod

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/lib/pq"
)

// ArgsUpdateProduct 修改商品参数
type ArgsUpdateProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes" check:"title" min:"1" max:"300" empty:"true"`
	//商品描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//封面ID
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs" check:"ids" empty:"true"`
	//描述图组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//货物重量
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//建议售价
	Price int64 `json:"price"`
	//不含税价格
	PriceNoTax int64 `db:"price_no_tax" json:"priceNoTax" check:"price" empty:"true"`
}

// UpdateProduct 修改商品
func UpdateProduct(args ArgsUpdateProduct) {
	CoreNats.PushDataNoErr("mall_core_product_update", "/mall/core/product_update", "", 0, "", args)
}
