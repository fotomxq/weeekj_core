package ERPProduct

import (
	BaseService "github.com/fotomxq/weeekj_core/v5/base/service"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//公司被删除
	CoreNats.SubDataByteNoErr("service_company", "/service/company", subNatsDeleteCompany)
	//更新产品资料
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "ERP产品更新通知",
		Description:  "",
		EventSubType: "all",
		Code:         "erp_product_update",
		EventType:    "nats",
		EventURL:     "/erp/product/update",
		//TODO:待补充
		EventParams: "",
	})
	CoreNats.SubDataByteNoErr("erp_product_update", "/erp/product/update", subNatsUpdateProduct)
	//注册服务
	_ = BaseService.SetService(&BaseService.ArgsSetService{
		ExpireAt:     CoreFilter.GetNowTimeCarbon().AddDay().Time,
		Name:         "ERP产品删除通知",
		Description:  "ERP产品删除通知",
		EventSubType: "push",
		Code:         "erp_product_delete",
		EventType:    "nats",
		EventURL:     "/erp/product/delete",
		EventParams:  "<<id>>:int64:产品ID",
	})
}

// 删除公司处理
func subNatsDeleteCompany(_ *nats.Msg, _ string, companyID int64, _ string, _ []byte) {
	//删除所有公司关联的产品信息
	subNatsDeleteCompanyStep1(companyID)
	//删除公司关联的品牌设置
	subNatsDeleteCompanyStep2(companyID)
}

func subNatsDeleteCompanyStep1(companyID int64) {
	logAppend := "erp product sub nats delete company, step 1, "
	//找到所有数据
	var dataList []FieldsProductCompany
	err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM erp_product_company WHERE company_id = $1 AND delete_at < to_timestamp(1000000)", companyID)
	if err != nil || len(dataList) < 1 {
		return
	}
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_product_company", "company_id = :company_id", map[string]interface{}{
		"company_id": companyID,
	})
	if err != nil {
		CoreLog.Error(logAppend, err)
		return
	}
	//遍历删除缓冲
	for _, v := range dataList {
		deleteProductCompanyCache(v.ID)
	}
}

func subNatsDeleteCompanyStep2(companyID int64) {
	logAppend := "erp product sub nats delete company, step 2, "
	var page int64 = 1
	for {
		dataList, _, err := GetBrandBindList(&ArgsGetBrandBindList{
			Pages: CoreSQL2.ArgsPages{
				Page: page,
				Max:  100,
				Sort: "id",
				Desc: false,
			},
			OrgID:     -1,
			BrandID:   -1,
			CompanyID: companyID,
			ProductID: -1,
			IsRemove:  false,
		})
		if err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			err = DeleteBrandBind(&ArgsDeleteBrandBind{
				OrgID:     v.OrgID,
				BrandID:   v.BrandID,
				CompanyID: v.CompanyID,
				ProductID: v.ProductID,
			})
			if err != nil {
				CoreLog.Error(logAppend, "delete brand bind, id: ", v.ID, ", ", err)
			}
		}
		page += 1
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
		OrgID:            productData.OrgID,
		CompanyID:        productData.CompanyID,
		CompanyName:      productData.CompanyName,
		SortID:           productData.SortID,
		Tags:             productData.Tags,
		SN:               productData.SN,
		Code:             productData.Code,
		PinYin:           productData.PinYin,
		EnName:           productData.EnName,
		ModelTypeID:      productData.ModelTypeID,
		ManufacturerName: productData.ManufacturerName,
		Title:            args.Title,
		TitleDes:         args.TitleDes,
		Des:              args.Des,
		CoverFileIDs:     args.CoverFileIDs,
		ExpireHour:       productData.ExpireHour,
		Weight:           args.Weight,
		SizeW:            productData.SizeW,
		SizeH:            productData.SizeH,
		SizeZ:            productData.SizeZ,
		PackType:         productData.PackType,
		PackUnitName:     productData.PackUnitName,
		PackUnit:         productData.PackUnit,
		TipPrice:         args.Price,
		TipTaxPrice:      args.Price,
		IsDiscount:       productData.IsDiscount,
		Currency:         productData.Currency,
		CostPrice:        productData.CostPrice,
		Tax:              productData.Tax,
		TaxCostPrice:     productData.TaxCostPrice,
		RebatePrice:      productData.RebatePrice,
		Params:           productData.Params,
		SyncMallCore:     false,
	})
	if err != nil {
		CoreLog.Error(appendLog, "update product data, id: ", args.ID, ", ", errCode, ", err: ", err)
		return
	}
}
