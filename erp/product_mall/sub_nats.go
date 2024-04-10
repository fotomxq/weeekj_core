package ERPProductMall

import (
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//ERP产品被删除
	CoreNats.SubDataByteNoErr("erp_product_delete", "/erp/product/delete", subNatsProductDelete)
}

// ERP产品被删除
func subNatsProductDelete(_ *nats.Msg, _ string, erpProductID int64, _ string, _ []byte) {
	var page int64 = 1
	for {
		dataList, _, _ := GetProductMallList(&ArgsGetProductMallList{
			Pages: CoreSQL2.ArgsPages{
				Page: page,
				Max:  100,
				Sort: "id",
				Desc: false,
			},
			OrgID:      -1,
			ProductID:  erpProductID,
			CategoryID: -1,
			IsRemove:   false,
			Search:     "",
		})
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			_ = DeleteProductMall(&ArgsDeleteProductMall{
				ID:    v.ID,
				OrgID: v.OrgID,
			})
		}
		page += 1
	}
}
