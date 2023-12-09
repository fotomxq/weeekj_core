package ERPProduct

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//公司被删除
	CoreNats.SubDataByteNoErr("/service/company", subNatsDeleteCompany)
	//更新产品资料
	CoreNats.SubDataByteNoErr("/erp/product/update", subNatsUpdateProduct)
}

// 删除公司处理
func subNatsDeleteCompany(_ *nats.Msg, action string, id int64, _ string, _ []byte) {
	//找到所有数据
	var dataList []FieldsProductCompany
	err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM erp_product_company WHERE company_id = $1 AND delete_at < to_timestamp(1000000)", id)
	if err != nil || len(dataList) < 1 {
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_product_company", "company_id = :company_id", map[string]interface{}{
		"company_id": id,
	})
	if err != nil {
		CoreLog.Error("erp product sub nats delete company, ", err)
		return
	}
	//遍历删除缓冲
	for _, v := range dataList {
		deleteProductCompanyCache(v.ID)
	}
}

// argsSubNatsUpdateProduct 更新产品资料参数
type argsSubNatsUpdateProduct struct {
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

// subNatsUpdateProduct 更新产品资料
func subNatsUpdateProduct(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	//日志
	appendLog := "erp product sub nats update product, "
	//获取参数
	var args argsSubNatsUpdateProduct
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error(appendLog, "get params, ", err)
		return
	}
	//获取商品
	productData := getProductByID(args.ID)
	if productData.ID < 1 {
		CoreLog.Error(appendLog, "get product data, id: ", args.ID)
		return
	}
	//更新数据
	_, errCode, err := SetProduct(&ArgsSetProduct{
		OrgID:        productData.OrgID,
		CompanyID:    productData.CompanyID,
		SortID:       productData.SortID,
		Tags:         productData.Tags,
		SN:           productData.SN,
		Code:         productData.Code,
		Title:        args.Title,
		TitleDes:     args.TitleDes,
		Des:          args.Des,
		CoverFileIDs: args.CoverFileIDs,
		ExpireHour:   productData.ExpireHour,
		Weight:       args.Weight,
		SizeW:        productData.SizeW,
		SizeH:        productData.SizeH,
		SizeZ:        productData.SizeZ,
		PackType:     productData.PackType,
		PackUnit:     productData.PackUnit,
		TipPrice:     args.Price,
		TipTaxPrice:  args.Price,
		Params:       productData.Params,
		SyncMallCore: false,
	})
	if err != nil {
		CoreLog.Error(appendLog, "update product data, id: ", args.ID, ", ", errCode, ", err: ", err)
		return
	}
}
