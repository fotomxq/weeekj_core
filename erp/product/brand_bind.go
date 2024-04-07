package ERPProduct

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetBrandBindList 获取绑定关系列表参数
type ArgsGetBrandBindList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//产品ID
	// 可选，如果给与值，则认为本数据为直接绑定到对应产品
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
}

// GetBrandBindList 获取绑定关系列表
func GetBrandBindList(args *ArgsGetBrandBindList) (dataList []FieldsBrandBind, dataCount int64, err error) {
	dataCount, err = brandBindDB.Select().SetFieldsList([]string{"id", "org_id", "brand_id", "company_id", "product_id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SelectList("((delete_at < to_timestamp(1000000) AND $1 = false) OR (delete_at >= to_timestamp(1000000) AND $1 = true)) AND (org_id = $2 OR $2 < 0) AND (brand_id = $3 OR $3 < 0) AND (company_id = $4 OR $4 < 0) AND (product_id = $5 OR $5 < 0)", args.IsRemove, args.OrgID, args.BrandID, args.CompanyID, args.ProductID).ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := GetBrandBindData(&ArgsGetBrandBindData{
			OrgID:     v.OrgID,
			BrandID:   v.BrandID,
			CompanyID: v.CompanyID,
			ProductID: v.ProductID,
		})
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetBrandBindData 获取品牌绑定关系参数
type ArgsGetBrandBindData struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//产品ID
	// 可选，如果给与值，则认为本数据为直接绑定到对应产品
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
}

// GetBrandBindData 获取品牌绑定关系
func GetBrandBindData(args *ArgsGetBrandBindData) (data FieldsBrandBind) {
	//获取缓冲
	cacheMark := getBrandBindCacheMark(args.OrgID, args.BrandID, args.CompanyID, args.ProductID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	//获取数据
	err := brandBindDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "org_id", "brand_id", "company_id", "product_id"}).AppendWhere("(org_id = $1 OR $1 < 0) AND brand_id = $2 AND company_id = $3 AND (product_id = $4 OR $4 < 0)", args.OrgID, args.BrandID, args.CompanyID, args.ProductID).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	//保存缓冲
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheBrandBindTime)
	//反馈
	return
}

// ArgsCheckBrandBind 检查品牌绑定关系参数
type ArgsCheckBrandBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//产品ID
	// 可选，如果给与值，则认为本数据为直接绑定到对应产品
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
}

// CheckBrandBind 检查品牌绑定关系
func CheckBrandBind(args *ArgsCheckBrandBind) (b bool) {
	//获取数据
	data := GetBrandBindData(&ArgsGetBrandBindData{
		OrgID:     args.OrgID,
		BrandID:   args.BrandID,
		CompanyID: args.CompanyID,
		ProductID: args.ProductID,
	})
	if data.ID < 1 {
		return
	}
	b = true
	//反馈
	return
}

// GetProductBrandID 获取产品的所属品牌
func GetProductBrandID(orgID int64, productID int64) (brandData FieldsBrand, brandDataBind FieldsBrandBind) {
	//获取产品直接关联的品牌
	dataList, _, _ := GetBrandBindList(&ArgsGetBrandBindList{
		Pages: CoreSQL2.ArgsPages{
			Page: 1,
			Max:  1,
			Sort: "id",
			Desc: false,
		},
		OrgID:     orgID,
		BrandID:   -1,
		CompanyID: -1,
		ProductID: productID,
		IsRemove:  false,
	})
	if len(dataList) < 1 {
		productData := getProductByID(productID)
		if productData.CompanyID < 1 {
			return
		}
		dataList, _, _ = GetBrandBindList(&ArgsGetBrandBindList{
			Pages: CoreSQL2.ArgsPages{
				Page: 1,
				Max:  1,
				Sort: "id",
				Desc: false,
			},
			OrgID:     orgID,
			BrandID:   -1,
			CompanyID: productData.CompanyID,
			ProductID: -1,
			IsRemove:  false,
		})
	}
	brandDataBind = GetBrandBindData(&ArgsGetBrandBindData{
		OrgID:     dataList[0].OrgID,
		BrandID:   dataList[0].BrandID,
		CompanyID: dataList[0].CompanyID,
		ProductID: dataList[0].ProductID,
	})
	brandData = getBrand(brandDataBind.BrandID)
	return
}

// ArgsCreateBrandBind 添加新品牌绑定关系参数
type ArgsCreateBrandBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//产品ID
	// 可选，如果给与值，则认为本数据为直接绑定到对应产品
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
}

// CreateBrandBind 添加新品牌绑定关系
func CreateBrandBind(args *ArgsCreateBrandBind) (id int64, err error) {
	//检查数据
	data := GetBrandBindData(&ArgsGetBrandBindData{
		OrgID:     args.OrgID,
		BrandID:   args.BrandID,
		CompanyID: args.CompanyID,
		ProductID: args.ProductID,
	})
	if data.ID > 0 {
		if CoreFilter.CheckHaveTime(data.DeleteAt) {
			id = data.ID
			err = brandBindDB.Update().NeedSoft(false).NeedUpdateTime().AddWhereID(data.ID).SetFields([]string{"delete_at"}).NamedExec(map[string]any{
				"delete_at": time.Time{},
			})
			return
		} else {
			//禁止重复绑定
			err = errors.New("have replace")
			return
		}
	}
	//检查产品是否存在绑定关系
	_, oldBindData := GetProductBrandID(args.OrgID, args.ProductID)
	if oldBindData.ID > 0 {
		if args.ProductID > 0 {
			if oldBindData.CompanyID < 1 {
				//产品已经绑定了品牌关系
				err = errors.New(fmt.Sprint("product have bind brand, arg product id: ", args.ProductID))
				return
			}
		} else {
			if args.CompanyID > 0 {
				if args.CompanyID == oldBindData.CompanyID {
					//产品所属公司已经绑定了品牌关系
					err = errors.New("company have bind brand")
					return
				}
			}
		}
	}
	//创建数据
	id, err = brandBindDB.Insert().SetFields([]string{"org_id", "brand_id", "company_id", "product_id"}).Add(map[string]any{
		"org_id":     args.OrgID,
		"brand_id":   args.BrandID,
		"company_id": args.CompanyID,
		"product_id": args.ProductID,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsDeleteBrandBind 删除产品绑定关系参数
type ArgsDeleteBrandBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//产品ID
	// 可选，如果给与值，则认为本数据为直接绑定到对应产品
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
}

// DeleteBrandBind 删除产品绑定关系
func DeleteBrandBind(args *ArgsDeleteBrandBind) (err error) {
	data := GetBrandBindData(&ArgsGetBrandBindData{
		OrgID:     args.OrgID,
		BrandID:   args.BrandID,
		CompanyID: args.CompanyID,
		ProductID: args.ProductID,
	})
	if data.ID < 1 {
		return
	}
	err = brandBindDB.Delete().NeedSoft(true).AddWhereID(data.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteBrandBindCache(args.OrgID, args.BrandID, args.CompanyID, args.ProductID)
	return
}
