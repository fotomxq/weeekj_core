package ERPProductMod

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	"github.com/lib/pq"
)

type ArgsUpdateProduct struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes" check:"des" min:"1" max:"300" empty:"true"`
	//商品描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面ID
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs"`
	//货物重量
	// 单位g
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//建议售价
	Price int64 `json:"price"`
}

func UpdateProduct(args ArgsUpdateProduct) {
	CoreNats.PushDataNoErr("erp_product_update", "/erp/product/update", "", 0, "", args)
}
