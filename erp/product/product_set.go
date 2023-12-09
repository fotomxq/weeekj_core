package ERPProduct

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	MallCoreMod "gitee.com/weeekj/weeekj_core/v5/mall/core/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceCompany "gitee.com/weeekj/weeekj_core/v5/service/company"
	"github.com/lib/pq"
)

// ArgsSetProduct 设置产品参数
type ArgsSetProduct struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//选择供应商
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供货商名称
	CompanyName string `db:"company_name" json:"companyName" check:"title" min:"1" max:"300" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//SN
	SN string `db:"sn" json:"sn"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
	//拼音助记码
	PinYin string `db:"pin_yin" json:"pinYin" check:"des" min:"1" max:"300" empty:"true"`
	//英文名称
	EnName string `db:"en_name" json:"enName" check:"des" min:"1" max:"300" empty:"true"`
	//生产厂商名称
	ManufacturerName string `db:"manufacturer_name" json:"manufacturerName" check:"des" min:"1" max:"300" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//副标题
	TitleDes string `db:"title_des" json:"titleDes" check:"des" min:"1" max:"300" empty:"true"`
	//商品描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面ID
	CoverFileIDs pq.Int64Array `db:"cover_file_ids" json:"coverFileIDs"`
	//保质期小时
	// 如果输入小于1，则永远不过期
	ExpireHour int `db:"expire_hour" json:"expireHour" check:"intThan0" empty:"true"`
	//货物重量
	// 单位g
	Weight int `db:"weight" json:"weight" check:"intThan0" empty:"true"`
	//存储尺寸
	SizeW int `db:"size_w" json:"sizeW" check:"intThan0" empty:"true"`
	SizeH int `db:"size_h" json:"sizeH" check:"intThan0" empty:"true"`
	SizeZ int `db:"size_z" json:"sizeZ" check:"intThan0" empty:"true"`
	//规格类型
	// 0 盒装; 1 袋装; 3 散装; 4 瓶装
	PackType int `db:"pack_type" json:"packType"`
	//包装单位名称
	PackUnitName string `db:"pack_unit_name" json:"packUnitName" check:"des" min:"1" max:"100" empty:"true"`
	//包装内部含有产品数量
	PackUnit int `db:"pack_unit" json:"packUnit" check:"intThan0" empty:"true"`
	//建议零售价（不含税）
	TipPrice int64 `db:"tip_price" json:"tipPrice" check:"price" empty:"true"`
	//建议零售价（含税）
	// 该建议价格用于可直接用于最终零售价填入
	TipTaxPrice int64 `db:"tip_tax_price" json:"tipTaxPrice" check:"price" empty:"true"`
	//是否允许打折
	IsDiscount bool `db:"is_discount" json:"isDiscount" check:"bool" empty:"true"`
	//货币
	Currency int `db:"currency" json:"currency" check:"currency"`
	//单价成本（不含税）
	CostPrice int64 `db:"cost_price" json:"costPrice" check:"price" empty:"true"`
	//税率
	// 实际税率=tax/100000
	Tax int64 `db:"tax" json:"tax"`
	//单价成本（含税）
	TaxCostPrice int64 `db:"tax_cost_price" json:"taxCostPrice" check:"price" empty:"true"`
	//返利设计
	RebatePrice FieldsProductRebateList `db:"rebate_price" json:"rebatePrice" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//是否同步商品
	SyncMallCore bool `json:"syncMallCore" check:"bool"`
}

// SetProduct 设置产品
func SetProduct(args *ArgsSetProduct) (data FieldsProduct, errCode string, err error) {
	//修正参数
	if args.Tags == nil || len(args.Tags) < 1 {
		args.Tags = pq.Int64Array{}
	}
	if args.CoverFileIDs == nil || len(args.CoverFileIDs) < 1 {
		args.CoverFileIDs = pq.Int64Array{}
	}
	if args.Params == nil || len(args.Params) < 1 {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//如果条形码不存在，则自动构建
	if args.Code == "" {
		args.Code = CoreFilter.GetRandStr4(50)
		data = getProductByCode(args.OrgID, args.Code)
		if data.ID > 0 {
			errCode = "err_code_replace"
			err = errors.New("product code is replace")
			return
		}
	}
	//检查编码是否重复
	data = getProductByCode(args.OrgID, args.Code)
	//SN不能重复
	if args.SN != "" {
		findSNData := GetProductBySN(args.OrgID, args.SN)
		if findSNData.ID > 0 && findSNData.ID != data.ID {
			errCode = "err_sn_replace"
			err = errors.New("product code is replace")
			return
		}
	}
	//检查公司是否存在
	var companyData ServiceCompany.FieldsCompany
	var productCompanyData FieldsProductCompany
	var haveProductCompany bool
	if args.CompanyID > 0 {
		companyData, err = ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
			ID:    args.CompanyID,
			OrgID: args.OrgID,
		})
		if err != nil {
			errCode = "err_no_data"
			err = errors.New("company not exist")
			return
		}
		//检查供货商是否为该商品的供货商
		productCompanyData, haveProductCompany = CheckProductCompany(args.OrgID, data.ID, args.CompanyID)
	}
	//CoreLog.Warn("companyData", companyData)
	//CoreLog.Warn("productCompanyData", productCompanyData)
	//编辑或创建产品
	if data.ID > 0 && CoreFilter.EqID2(args.OrgID, data.OrgID) {
		//修正供货商信息
		if args.CompanyID > 0 {
			args.CompanyID = productCompanyData.CompanyID
		}
		//更新数据
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_product SET update_at = NOW(), delete_at = to_timestamp(0), company_id = :company_id, company_name = :company_name, sort_id = :sort_id, tags = :tags, sn = :sn, code = :code, pin_yin = :pin_yin, en_name = :en_name, manufacturer_name = :manufacturer_name, title = :title, title_des = :title_des, des = :des, cover_file_ids = :cover_file_ids, expire_hour = :expire_hour, weight = :weight, size_w = :size_w, size_h = :size_h, size_z = :size_z, pack_type = :pack_type, pack_unit_name = :pack_unit_name, pack_unit = :pack_unit, tip_price = :tip_price, tip_tax_price = :tip_tax_price, is_discount = :is_discount, currency = :currency, cost_price = :cost_price, tax = :tax, tax_cost_price = :tax_cost_price, rebate_price = :rebate_price, params = :params WHERE id = :id", map[string]interface{}{
			"id":                data.ID,
			"company_id":        args.CompanyID,
			"company_name":      args.CompanyName,
			"sort_id":           args.SortID,
			"tags":              args.Tags,
			"sn":                args.SN,
			"code":              args.Code,
			"pin_yin":           args.PinYin,
			"en_name":           args.EnName,
			"manufacturer_name": args.ManufacturerName,
			"title":             args.Title,
			"title_des":         args.TitleDes,
			"des":               args.Des,
			"cover_file_ids":    args.CoverFileIDs,
			"expire_hour":       args.ExpireHour,
			"weight":            args.Weight,
			"size_w":            args.SizeW,
			"size_h":            args.SizeH,
			"size_z":            args.SizeZ,
			"pack_type":         args.PackType,
			"pack_unit_name":    args.PackUnitName,
			"pack_unit":         args.PackUnit,
			"tip_price":         args.TipPrice,
			"tip_tax_price":     args.TipTaxPrice,
			"is_discount":       args.IsDiscount,
			"currency":          args.Currency,
			"cost_price":        args.CostPrice,
			"tax":               args.Tax,
			"tax_cost_price":    args.TaxCostPrice,
			"rebate_price":      args.RebatePrice,
			"params":            args.Params,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
		//删除缓冲
		deleteProductCache(data.ID)
	} else {
		//创建产品
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "erp_product", "INSERT INTO erp_product (org_id, company_id, company_name, sort_id, tags, sn, code, pin_yin, en_name, manufacturer_name, title, title_des, des, cover_file_ids, expire_hour, weight, size_w, size_h, size_z, pack_type, pack_unit_name, pack_unit, tip_price, tip_tax_price, is_discount, currency, cost_price, tax, tax_cost_price, rebate_price, params) VALUES (:org_id, :company_id, :company_name, :sort_id, :tags, :sn, :code, :pin_yin, :en_name, :manufacturer_name, :title, :title_des, :des, :cover_file_ids, :expire_hour, :weight, :size_w, :size_h, :size_z, :pack_type, :pack_unit_name, :pack_unit, :tip_price, :tip_tax_price, :is_discount, :currency, :cost_price, :tax, :tax_cost_price, :rebate_price, :params)", map[string]interface{}{
			"org_id":            args.OrgID,
			"company_id":        args.CompanyID,
			"company_name":      args.CompanyName,
			"sort_id":           args.SortID,
			"tags":              args.Tags,
			"code":              args.Code,
			"pin_yin":           args.PinYin,
			"en_name":           args.EnName,
			"manufacturer_name": args.ManufacturerName,
			"sn":                args.SN,
			"title":             args.Title,
			"title_des":         args.TitleDes,
			"des":               args.Des,
			"cover_file_ids":    args.CoverFileIDs,
			"expire_hour":       args.ExpireHour,
			"weight":            args.Weight,
			"size_w":            args.SizeW,
			"size_h":            args.SizeH,
			"size_z":            args.SizeZ,
			"pack_type":         args.PackType,
			"pack_unit_name":    args.PackUnitName,
			"pack_unit":         args.PackUnit,
			"tip_price":         args.TipPrice,
			"tip_tax_price":     args.TipTaxPrice,
			"is_discount":       args.IsDiscount,
			"currency":          args.Currency,
			"cost_price":        args.CostPrice,
			"tax":               args.Tax,
			"tax_cost_price":    args.TaxCostPrice,
			"rebate_price":      args.RebatePrice,
			"params":            args.Params,
		}, &data)
		if err != nil {
			errCode = "err_insert"
			return
		}
		//如果不存在供货商则自动创建
		if !haveProductCompany && args.CompanyID > 0 {
			var tax int64 = 0
			if args.TipTaxPrice != args.TipPrice && args.TipPrice > 0 {
				tax = int64(((float64(args.TipPrice) - float64(args.TipTaxPrice)) / float64(args.TipPrice)) * 100000)
			}
			err = CreateProductCompany(&ArgsCreateProductCompany{
				OrgID:        data.OrgID,
				ProductID:    data.ID,
				CompanyID:    companyData.ID,
				Currency:     86,
				CostPrice:    args.TipPrice,
				Tax:          tax,
				TaxCostPrice: args.TipTaxPrice,
				RebatePrice:  nil,
				TipPrice:     args.TipPrice,
				TipTaxPrice:  args.TipTaxPrice,
				Params:       nil,
			})
			if err != nil {
				errCode = "err_insert"
				return
			}
		}
	}
	//重新获取数据
	data = getProductByID(data.ID)
	//如果启动同步
	if args.SyncMallCore {
		mallList := MallCoreMod.GetProductListByWarehouseProductID(data.ID)
		for _, v := range mallList {
			MallCoreMod.UpdateProduct(MallCoreMod.ArgsUpdateProduct{
				ID:           v.ID,
				Title:        data.Title,
				TitleDes:     data.TitleDes,
				Des:          data.Des,
				CoverFileIDs: data.CoverFileIDs,
				DesFiles:     v.DesFiles,
				Weight:       data.Weight,
				Price:        data.TipTaxPrice,
				PriceNoTax:   data.TipPrice,
			})
		}
	}
	//重新获取产品
	data = getProductByID(data.ID)
	//反馈
	return
}

// SetProduct2 设置产品2
func SetProduct2(args *ArgsSetProduct) (data FieldsProduct, errCode string, err error) {
	//修正参数
	if args.Tags == nil || len(args.Tags) < 1 {
		args.Tags = pq.Int64Array{}
	}
	if args.CoverFileIDs == nil || len(args.CoverFileIDs) < 1 {
		args.CoverFileIDs = pq.Int64Array{}
	}
	if args.Params == nil || len(args.Params) < 1 {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//如果条形码不存在，则自动构建
	if args.Code == "" {
		args.Code = CoreFilter.GetRandStr4(50)
		data = getProductByCode(args.OrgID, args.Code)
		if data.ID > 0 {
			errCode = "err_code_replace"
			err = errors.New("product code is replace")
			return
		}
	}
	//检查编码是否重复
	data = getProductByCode(args.OrgID, args.Code)
	//SN不能重复
	if args.SN != "" {
		findSNData := GetProductBySN(args.OrgID, args.SN)
		if findSNData.ID > 0 && findSNData.ID != data.ID {
			errCode = "err_sn_replace"
			err = errors.New("product code is replace")
			return
		}
	}
	//检查公司是否存在
	var companyData ServiceCompany.FieldsCompany
	if args.CompanyID > 0 {
		companyData, err = ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
			ID:    args.CompanyID,
			OrgID: args.OrgID,
		})
		if err != nil {
			errCode = "err_no_data"
			err = errors.New("company not exist")
			return
		}
	}
	//编辑或创建产品
	if data.ID > 0 && CoreFilter.EqID2(args.OrgID, data.OrgID) {
		//更新数据
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_product SET update_at = NOW(), delete_at = to_timestamp(0), company_id = :company_id, company_name = :company_name, sort_id = :sort_id, tags = :tags, sn = :sn, code = :code, pin_yin = :pin_yin, en_name = :en_name, manufacturer_name = :manufacturer_name, title = :title, title_des = :title_des, des = :des, cover_file_ids = :cover_file_ids, expire_hour = :expire_hour, weight = :weight, size_w = :size_w, size_h = :size_h, size_z = :size_z, pack_type = :pack_type, pack_unit_name = :pack_unit_name, pack_unit = :pack_unit, tip_price = :tip_price, tip_tax_price = :tip_tax_price, is_discount = :is_discount, currency = :currency, cost_price = :cost_price, tax = :tax, tax_cost_price = :tax_cost_price, rebate_price = :rebate_price, params = :params WHERE id = :id", map[string]interface{}{
			"id":                data.ID,
			"company_id":        companyData.ID,
			"company_name":      args.CompanyName,
			"sort_id":           args.SortID,
			"tags":              args.Tags,
			"sn":                args.SN,
			"code":              args.Code,
			"pin_yin":           args.PinYin,
			"en_name":           args.EnName,
			"manufacturer_name": args.ManufacturerName,
			"title":             args.Title,
			"title_des":         args.TitleDes,
			"des":               args.Des,
			"cover_file_ids":    args.CoverFileIDs,
			"expire_hour":       args.ExpireHour,
			"weight":            args.Weight,
			"size_w":            args.SizeW,
			"size_h":            args.SizeH,
			"size_z":            args.SizeZ,
			"pack_type":         args.PackType,
			"pack_unit_name":    args.PackUnitName,
			"pack_unit":         args.PackUnit,
			"tip_price":         args.TipPrice,
			"tip_tax_price":     args.TipTaxPrice,
			"is_discount":       args.IsDiscount,
			"currency":          args.Currency,
			"cost_price":        args.CostPrice,
			"tax":               args.Tax,
			"tax_cost_price":    args.TaxCostPrice,
			"rebate_price":      args.RebatePrice,
			"params":            args.Params,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
		//删除缓冲
		deleteProductCache(data.ID)
	} else {
		//创建产品
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "erp_product", "INSERT INTO erp_product (org_id, company_id, company_name, sort_id, tags, sn, code, pin_yin, en_name, manufacturer_name, title, title_des, des, cover_file_ids, expire_hour, weight, size_w, size_h, size_z, pack_type, pack_unit_name, pack_unit, tip_price, tip_tax_price, is_discount, currency, cost_price, tax, tax_cost_price, rebate_price, params) VALUES (:org_id, :company_id, :company_name, :sort_id, :tags, :sn, :code, :pin_yin, :en_name, :manufacturer_name, :title, :title_des, :des, :cover_file_ids, :expire_hour, :weight, :size_w, :size_h, :size_z, :pack_type, :pack_unit_name, :pack_unit, :tip_price, :tip_tax_price, :is_discount, :currency, :cost_price, :tax, :tax_cost_price, :rebate_price, :params)", map[string]interface{}{
			"org_id":            args.OrgID,
			"company_id":        args.CompanyID,
			"company_name":      args.CompanyName,
			"sort_id":           args.SortID,
			"tags":              args.Tags,
			"code":              args.Code,
			"pin_yin":           args.PinYin,
			"en_name":           args.EnName,
			"manufacturer_name": args.ManufacturerName,
			"sn":                args.SN,
			"title":             args.Title,
			"title_des":         args.TitleDes,
			"des":               args.Des,
			"cover_file_ids":    args.CoverFileIDs,
			"expire_hour":       args.ExpireHour,
			"weight":            args.Weight,
			"size_w":            args.SizeW,
			"size_h":            args.SizeH,
			"size_z":            args.SizeZ,
			"pack_type":         args.PackType,
			"pack_unit_name":    args.PackUnitName,
			"pack_unit":         args.PackUnit,
			"tip_price":         args.TipPrice,
			"tip_tax_price":     args.TipTaxPrice,
			"is_discount":       args.IsDiscount,
			"currency":          args.Currency,
			"cost_price":        args.CostPrice,
			"tax":               args.Tax,
			"tax_cost_price":    args.TaxCostPrice,
			"rebate_price":      args.RebatePrice,
			"params":            args.Params,
		}, &data)
		if err != nil {
			errCode = "err_insert"
			return
		}
	}
	//反馈
	return
}
