package ERPProduct

import (
	"errors"
	"fmt"
	BaseBPM "github.com/fotomxq/weeekj_core/v5/base/bpm"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// DataProductVal 标准化参数和反馈结构
type DataProductVal struct {
	//顺序序号
	OrderNum int64 `db:"order_num" json:"orderNum"`
	//插槽值
	SlotID int64 `db:"slot_id" json:"slotID" check:"id"`
	//值(字符串)
	DataValue string `db:"data_value" json:"dataValue"`
	//值(浮点数)
	DataValueNum float64 `db:"data_value_num" json:"dataValueNum"`
	//值(整数)
	DataValueInt int64 `db:"data_value_int" json:"dataValueInt"`
	//参数
	Params string `db:"params" json:"params"`
}

// ArgsGetProductVals 获取产品预设模板值参数
type ArgsGetProductVals struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
}

// GetProductVals 获取产品预设模板值
// 根据产品关联的分类、品牌，获取产品预设模板数据包
// 如已经存在值，则反馈具体值；否则反馈默认数据包
func GetProductVals(args *ArgsGetProductVals) (data []DataProductVal, err error) {
	var rawList []FieldsProductVals
	rawList, err = getProductValsRaw(args)
	if err != nil {
		return
	}
	for _, v := range rawList {
		data = append(data, DataProductVal{
			OrderNum:     v.OrderNum,
			SlotID:       v.SlotID,
			DataValue:    v.DataValue,
			DataValueNum: v.DataValueNum,
			DataValueInt: v.DataValueInt,
			Params:       v.Params,
		})
	}
	return
}

// ArgsGetProductValsAndDefault 获取产品预设模板值和默认值参数
type ArgsGetProductValsAndDefault struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品关联的公司ID
	// 可选，可留空
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//是否可以突破模板数据
	// 支持超出模板定义范畴的自定义数据包
	HaveMoreData bool `json:"haveMoreData"`
}

// GetProductValsAndDefault 获取产品预设模板值和默认值
// 如果存在数据，则反馈数据；否则反馈预设值
func GetProductValsAndDefault(args *ArgsGetProductValsAndDefault) (dataList []DataProductVal, errCode string, err error) {
	//获取已经存在的数据包
	rawList, _ := GetProductVals(&ArgsGetProductVals{
		OrgID:     args.OrgID,
		ProductID: args.ProductID,
	})
	//获取产品绑定的模板关系
	var templateBindData FieldsTemplateBind
	templateBindData, errCode, err = GetProductValsTemplateID(&ArgsGetProductValsTemplateID{
		OrgID:     args.OrgID,
		ProductID: args.ProductID,
		CompanyID: args.CompanyID,
	})
	if err != nil {
		return
	}
	//获取当前插槽数据包
	var bpmSlotList []BaseBPM.FieldsSlot
	bpmSlotList, errCode, err = GetTemplateBPMThemeSlotData(args.OrgID, templateBindData.TemplateID)
	if err != nil {
		return
	}
	//交叉匹配数据
	if args.HaveMoreData {
		step := len(rawList)
		for _, v := range rawList {
			dataList = append(dataList, v)
		}
		for _, v := range bpmSlotList {
			isFind := false
			for _, v2 := range rawList {
				if v.ID == v2.SlotID {
					isFind = true
					break
				}
			}
			if isFind {
				continue
			}
			dataList = append(dataList, DataProductVal{
				OrderNum:     int64(step),
				SlotID:       v.ID,
				DataValue:    v.DefaultValue,
				DataValueNum: CoreFilter.GetFloat64ByStringNoErr(v.DefaultValue),
				DataValueInt: CoreFilter.GetInt64ByStringNoErr(v.DefaultValue),
				Params:       v.Params,
			})
			step += 1
		}
	} else {
		for _, v := range bpmSlotList {
			isFind := false
			for _, v2 := range rawList {
				if v.ID == v2.SlotID {
					dataList = append(dataList, v2)
					isFind = true
					break
				}
			}
			if isFind {
				continue
			}
			dataList = append(dataList, DataProductVal{
				OrderNum:     int64(len(dataList)),
				SlotID:       v.ID,
				DataValue:    v.DefaultValue,
				DataValueNum: CoreFilter.GetFloat64ByStringNoErr(v.DefaultValue),
				DataValueInt: CoreFilter.GetInt64ByStringNoErr(v.DefaultValue),
				Params:       v.Params,
			})
		}
	}
	//反馈
	return
}

// ArgsGetProductValsTemplateID 获取产品模板ID参数
type ArgsGetProductValsTemplateID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品关联的公司ID
	// 可选，可留空
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
}

// GetProductValsTemplateID 获取产品模板ID
func GetProductValsTemplateID(args *ArgsGetProductValsTemplateID) (templateBindData FieldsTemplateBind, errCode string, err error) {
	//找到产品对应的模板体系
	productData := getProductByID(args.ProductID)
	if productData.ID < 1 {
		errCode = "err_erp_product_no_data"
		err = errors.New(fmt.Sprint("product no data id: ", args.ProductID))
		return
	}
	//检查公司是否和产品存在绑定关系
	if args.CompanyID > 0 {
		_, b := CheckProductCompany(args.OrgID, args.ProductID, args.CompanyID)
		if !b {
			errCode = "err_erp_product_company"
			err = errors.New(fmt.Sprint("product company no bind by product id: ", args.ProductID, ", company id: ", args.CompanyID))
			return
		}
	}
	//获取产品与型号绑定关系（优先级最高）
	if productData.ModelTypeID > 0 {
		modelTypeData := GetModelType(productData.ModelTypeID, args.OrgID)
		if modelTypeData.ID > 0 {
			templateBindList, _, _ := GetTemplateBindList(&ArgsGetTemplateBindList{
				Pages: CoreSQL2.ArgsPages{
					Page: 1,
					Max:  1,
					Sort: "id",
					Desc: false,
				},
				OrgID:       args.OrgID,
				TemplateID:  -1,
				CategoryID:  -1,
				BrandID:     -1,
				ModelTypeID: modelTypeData.ID,
				IsRemove:    false,
			})
			if len(templateBindList) > 0 {
				templateBindData = templateBindList[0]
			}
		}
	}
	//获取产品品牌绑定关系（优先级中等）
	if templateBindData.ID < 1 {
		var brandBindData FieldsBrandBind
		brandBindList, _, _ := GetBrandBindList(&ArgsGetBrandBindList{
			Pages: CoreSQL2.ArgsPages{
				Page: 1,
				Max:  1,
				Sort: "id",
				Desc: false,
			},
			OrgID:     args.OrgID,
			BrandID:   -1,
			CompanyID: -1,
			ProductID: args.ProductID,
			IsRemove:  false,
		})
		if len(brandBindList) > 0 {
			brandBindData = brandBindList[0]
		} else {
			if args.CompanyID > 0 {
				brandBindList, _, _ = GetBrandBindList(&ArgsGetBrandBindList{
					Pages: CoreSQL2.ArgsPages{
						Page: 1,
						Max:  1,
						Sort: "id",
						Desc: false,
					},
					OrgID:     args.OrgID,
					BrandID:   -1,
					CompanyID: args.CompanyID,
					ProductID: -1,
					IsRemove:  false,
				})
				if len(brandBindList) > 0 {
					brandBindData = brandBindList[0]
				}
			}
		}
		//检查品牌和模板的绑定关系
		if brandBindData.ID > 0 {
			templateBindList, _, _ := GetTemplateBindList(&ArgsGetTemplateBindList{
				Pages: CoreSQL2.ArgsPages{
					Page: 1,
					Max:  1,
					Sort: "id",
					Desc: false,
				},
				OrgID:       args.OrgID,
				TemplateID:  -1,
				CategoryID:  -1,
				BrandID:     brandBindData.BrandID,
				ModelTypeID: -1,
				IsRemove:    false,
			})
			if len(templateBindList) > 0 {
				templateBindData = templateBindList[0]
			}
		}
	}
	//获取产品分类绑定关系（优先级最低）
	if templateBindData.ID < 1 {
		templateBindData = getTemplateBindRecursionByCategoryID(args.OrgID, productData.SortID)
	}
	//检查是否具备模板
	if templateBindData.ID < 1 {
		errCode = "err_erp_product_no_template"
		err = errors.New("product not bind template")
		return
	}
	//反馈
	return
}

// ArgsGetValsByBrandOrCategoryID 通过分类或品牌获取数据包参数
type ArgsGetValsByBrandOrCategoryID struct {
	//组织ID
	// 请勿给-1，否则无法清理缓冲
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id" empty:"true"`
	//规格型号
	// 三选一，分类ID、品牌ID、规格型号ID
	ModelTypeID int64 `db:"model_type_id" json:"modelTypeID" check:"id" empty:"true"`
}

// GetValsByBrandOrCategoryID 通过分类或品牌获取数据包
func GetValsByBrandOrCategoryID(args *ArgsGetValsByBrandOrCategoryID) (templateID int64, themeID int64, bpmSlotList []BaseBPM.FieldsSlot, errCode string, err error) {
	var templateBindData FieldsTemplateBind
	templateBindList, _, _ := GetTemplateBindList(&ArgsGetTemplateBindList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  1,
			Sort: "id",
			Desc: false,
		},
		OrgID:       args.OrgID,
		TemplateID:  -1,
		CategoryID:  -1,
		BrandID:     args.BrandID,
		ModelTypeID: args.ModelTypeID,
		IsRemove:    false,
	})
	if len(templateBindList) > 0 {
		templateBindData = templateBindList[0]
	}
	if templateBindData.ID < 1 {
		templateBindData = getTemplateBindRecursionByCategoryID(args.OrgID, args.CategoryID)
	}
	if templateBindData.ID < 1 {
		errCode = "err_erp_product_no_template"
		err = errors.New("product not bind template")
		return
	}
	templateData := GetTemplate(templateBindData.TemplateID, templateBindData.OrgID)
	if templateData.ID < 1 {
		errCode = "err_erp_product_no_exist_template"
		err = errors.New("product bind template not exist")
		return
	}
	templateID = templateBindData.TemplateID
	themeID = templateData.BPMThemeID
	//如果模板拿到关联主题
	bpmSlotList, errCode, err = GetTemplateBPMThemeSlotData(args.OrgID, templateID)
	if err != nil {
		return
	}
	return
}

// ArgsSetProductVals 设置产品数据参数
type ArgsSetProductVals struct {
	//组织ID
	// 请勿给-1，否则无法清理缓冲
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品关联的公司ID
	// 可选，可留空
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//数据结构
	Vals []DataProductVal `db:"vals" json:"vals"`
	//是否可以突破模板数据
	// 支持超出模板定义范畴的自定义数据包
	HaveMoreData bool `json:"haveMoreData"`
}

// SetProductVals 设置产品数据
func SetProductVals(args *ArgsSetProductVals) (errCode string, err error) {
	//检查产品绑定的模板
	var templateBindData FieldsTemplateBind
	templateBindData, errCode, err = GetProductValsTemplateID(&ArgsGetProductValsTemplateID{
		OrgID:     args.OrgID,
		ProductID: args.ProductID,
		CompanyID: args.CompanyID,
	})
	if err != nil {
		return
	}
	//获取当前插槽数据包
	var bpmSlotList []BaseBPM.FieldsSlot
	bpmSlotList, errCode, err = GetTemplateBPMThemeSlotData(args.OrgID, templateBindData.TemplateID)
	if err != nil {
		return
	}
	//重构数据包
	var newProductVals []DataProductVal
	if !args.HaveMoreData {
		//禁止自定义数据包，将移除提交的无效数据
		for _, v := range bpmSlotList {
			for _, v2 := range args.Vals {
				if v.ID == v2.SlotID {
					newProductVals = append(newProductVals, v2)
					break
				}
			}
		}
	} else {
		for _, v := range args.Vals {
			newProductVals = append(newProductVals, v)
		}
	}
	//获取产品当前数据
	var rawList []FieldsProductVals
	rawList, err = getProductValsRaw(&ArgsGetProductVals{
		OrgID:     args.OrgID,
		ProductID: args.ProductID,
	})
	//如果存在数据，则匹配修改
	if err == nil {
		for k, v := range rawList {
			for _, v2 := range newProductVals {
				if v.SlotID == v2.SlotID {
					rawList[k].OrderNum = v2.OrderNum
					rawList[k].DataValue = v2.DataValue
					rawList[k].DataValueNum = v2.DataValueNum
					rawList[k].DataValueInt = v2.DataValueInt
					rawList[k].Params = v2.Params
					break
				}
			}
		}
	}
	//优先遍历rawList，确保修改
	for _, v := range rawList {
		err = productValsDB.Update().SetFields([]string{"order_num", "data_value", "data_value_num", "data_value_int", "params"}).AddWhereID(v.ID).NamedExec(map[string]any{
			"order_num":      v.OrderNum,
			"data_value":     v.DataValue,
			"data_value_num": v.DataValueNum,
			"data_value_int": v.DataValueInt,
			"params":         v.Params,
		})
		if err != nil {
			errCode = "err_update"
			err = errors.New(fmt.Sprint("update product vals error: ", err))
			return
		}
	}
	for _, v := range newProductVals {
		isFind := false
		for _, v2 := range rawList {
			if v.SlotID == v2.SlotID {
				isFind = true
				break
			}
		}
		if !isFind {
			err = productValsDB.Insert().SetFields([]string{"org_id", "product_id", "template_id", "order_num", "slot_id", "data_value", "data_value_num", "data_value_int", "params"}).Add(map[string]any{
				"org_id":         args.OrgID,
				"product_id":     args.ProductID,
				"template_id":    templateBindData.TemplateID,
				"order_num":      v.OrderNum,
				"slot_id":        v.SlotID,
				"data_value":     v.DataValue,
				"data_value_num": v.DataValueNum,
				"data_value_int": v.DataValueInt,
				"params":         v.Params,
			}).ExecAndCheckID()
			if err != nil {
				errCode = "err_insert"
				err = errors.New(fmt.Sprint("insert product vals error: ", err))
				return
			}
		}
	}
	//清理丢失的数据
	for _, v := range rawList {
		isFind := false
		for _, v2 := range newProductVals {
			if v.SlotID == v2.SlotID {
				isFind = true
				break
			}
		}
		if !isFind {
			err = productValsDB.Delete().NeedSoft(true).AddWhereID(v.ID).ExecNamed(nil)
			if err != nil {
				errCode = "err_delete"
				err = errors.New(fmt.Sprint("delete product vals error: ", err))
				return
			}
		}
	}
	//清空缓存
	deleteProductValsCache(args.OrgID, args.ProductID)
	//反馈
	return
}

type ArgsClearProductVals struct {
	//组织ID
	// 请勿给-1，否则无法清理缓冲
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
}

// ClearProductVals 清空产品数据
func ClearProductVals(args *ArgsClearProductVals) (err error) {
	err = productValsDB.Delete().NeedSoft(true).AddWhereOrgID(args.OrgID).SetWhereAnd("product_id", args.ProductID).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteProductValsCache(args.OrgID, args.ProductID)
	return
}

// getProductValsRaw 获取产品扩展数据的原始数据
func getProductValsRaw(args *ArgsGetProductVals) (dataList []FieldsProductVals, err error) {
	cacheMark := getProductValsCacheMark(args.OrgID, args.ProductID)
	if err = Router2SystemConfig.MainCache.GetStruct(cacheMark, &dataList); err == nil && len(dataList) > 0 {
		return
	}
	err = productValsDB.Select().SetFieldsList([]string{"id", "create_at", "update_at", "delete_at", "org_id", "product_id", "template_id", "order_num", "slot_id", "data_value", "data_value_num", "data_value_int", "params"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(CoreSQL2.ArgsPages{
		Page: 1,
		Max:  100,
		Sort: "id",
		Desc: false,
	}).SetDeleteQuery("delete_at", false).SetIDQuery("org_id", args.OrgID).SetIDQuery("product_id", args.ProductID).SelectList("").Result(&dataList)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, dataList, cacheProductValsTime)
	return
}
