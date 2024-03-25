package ERPProduct

import (
	"errors"
	"fmt"
	BaseBPM "github.com/fotomxq/weeekj_core/v5/base/bpm"
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

// ArgsSetProductVals 设置产品数据参数
type ArgsSetProductVals struct {
	//组织ID
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
	//获取产品品牌绑定关系（优先级最高）
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
		CompanyID: args.CompanyID,
		ProductID: args.ProductID,
		IsRemove:  false,
	})
	if len(brandBindList) > 0 {
		brandBindData = brandBindList[0]
	}
	//检查品牌和模板的绑定关系
	var templateBindData FieldsTemplateBind
	if brandBindData.ID > 0 {
		templateBindList, _, _ := GetTemplateBindList(&ArgsGetTemplateBindList{
			Pages: CoreSQL2.ArgsPages{
				Page: 1,
				Max:  1,
				Sort: "id",
				Desc: false,
			},
			OrgID:      args.OrgID,
			TemplateID: -1,
			CategoryID: -1,
			BrandID:    templateBindData.BrandID,
			IsRemove:   false,
		})
		if len(templateBindList) > 0 {
			templateBindData = templateBindList[0]
		}
	}
	//如果不存在品牌绑定关系，则检查产品分类绑定关系
	if templateBindData.ID < 1 {
		templateBindData = getTemplateBindRecursionByCategoryID(args.OrgID, productData.SortID)
	}
	//检查是否具备模板
	if templateBindData.ID < 1 {
		errCode = "err_erp_product_no_template"
		err = errors.New("product not bind template")
		return
	}
	//通过绑定关系获取模板数据包
	templateData := GetTemplate(templateBindData.TemplateID, args.OrgID)
	if templateData.ID < 1 {
		errCode = "err_erp_product_no_exist_template"
		err = errors.New("product bind template not exist")
		return
	}
	//如果模板拿到关联主题
	bmpThemeData, _ := BaseBPM.GetThemeByID(&BaseBPM.ArgsGetThemeByID{
		ID: templateData.BPMThemeID,
	})
	if bmpThemeData.ID < 1 {
		errCode = "err_erp_product_no_exist_bpm_theme"
		err = errors.New("product bind template not exist bpm theme")
		return
	}
	//通过BPM主题，查询关联的插槽
	bpmSlotList, _, _ := BaseBPM.GetSlotList(&BaseBPM.ArgsGetSlotList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  999,
			Sort: "id",
			Desc: false,
		},
		ThemeCategoryID: -1,
		ThemeID:         bmpThemeData.ID,
		IsRemove:        false,
		Search:          "",
	})
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
		for _, v := range rawList {
			for _, v2 := range newProductVals {
				if v.SlotID == v2.SlotID {
					v.OrderNum = v2.OrderNum
					v.DataValue = v2.DataValue
					v.DataValueNum = v2.DataValueNum
					v.DataValueInt = v2.DataValueInt
					v.Params = v2.Params
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
				"template_id":    templateData.ID,
				"order_num":      v.OrderNum,
				"slot_id":        v.SlotID,
				"data_value":     v.DataValue,
				"data_value_num": v.DataValueNum,
				"data_value_int": v.DataValueInt,
				"params":         v.Params,
			}).ExecAndCheckID()
			errCode = "err_insert"
			err = errors.New(fmt.Sprint("insert product vals error: ", err))
			return
		}
	}
	//清空缓存
	deleteProductValsCache(args.OrgID, args.ProductID)
	//反馈
	return
}

type ArgsClearProductVals struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
}

// ClearProductVals 清空产品数据
func ClearProductVals(args *ArgsClearProductVals) (err error) {
	err = productValsDB.Delete().NeedSoft(true).AddWhereOrgID(args.OrgID).SetWhereStr(" AND product_id = :product_id", map[string]any{
		"product_id": args.ProductID,
	}).ExecNamed(nil)
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
	err = templateDB.Select().SetFieldsList([]string{"id", "create_at", "update_at", "delete_at", "org_id", "product_id", "template_id", "order_num", "slot_id", "data_value", "data_value_num", "data_value_int", "params"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(CoreSQL2.ArgsPages{
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
