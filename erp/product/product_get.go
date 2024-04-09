package ERPProduct

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"strings"
)

// ArgsGetProductList 获取产品列表参数
type ArgsGetProductList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//规格类型
	// -1 跳过; 0 盒装; 1 袋装; 3 散装; 4 瓶装
	PackType int `db:"pack_type" json:"packType"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索编码
	SearchCode string `json:"searchCode" check:"search" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetProductList 获取产品列表
func GetProductList(args *ArgsGetProductList) (dataList []FieldsProduct, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.PackType > -1 {
		where = where + " AND pack_type = :pack_type"
		maps["pack_type"] = args.PackType
	}
	if args.SearchCode != "" {
		where = where + " AND (code ILIKE '%' || :search_code || '%')"
		maps["search_code"] = args.SearchCode
	}
	if args.Search != "" {
		args.Search = strings.ReplaceAll(args.Search, " ", "")
		where = where + " AND (REPLACE(title, ' ', '') ILIKE '%' || :search || '%' OR REPLACE(title_des, ' ', '') ILIKE '%' || :search || '%' OR REPLACE(des, ' ', '') ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_product"
	var rawList []FieldsProduct
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "weight", "tip_price", "tip_tax_price"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getProductByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetProductListV2 获取产品列表V2参数
type ArgsGetProductListV2 struct {
	//分页
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//规格类型
	// -1 跳过; 0 盒装; 1 袋装; 3 散装; 4 瓶装
	PackType int `db:"pack_type" json:"packType"`
	//所属供应商ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//所属品牌ID
	BrandID int64 `db:"brand_id" json:"brandID" check:"id" empty:"true"`
	//所使用模板ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索编码
	SearchCode string `json:"searchCode" check:"search" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetProductListV2 获取产品列表V2
func GetProductListV2(args *ArgsGetProductListV2) (dataList []FieldsProduct, dataCount int64, err error) {
	var brandIDList pq.Int64Array
	var sortIDList pq.Int64Array
	if args.TemplateID > 0 {
		templateBindList, _, _ := GetTemplateBindList(&ArgsGetTemplateBindList{
			Pages: CoreSQL2.ArgsPages{
				Page: 1,
				Max:  30,
				Sort: "id",
				Desc: false,
			},
			OrgID:      args.OrgID,
			TemplateID: args.TemplateID,
			CategoryID: -1,
			BrandID:    -1,
			IsRemove:   false,
		})
		for _, v := range templateBindList {
			if v.BrandID > 0 {
				brandIDList = append(brandIDList, v.BrandID)
			}
			if v.CategoryID > 0 {
				sortIDList = append(sortIDList, v.CategoryID)
			}
		}
	}
	var companyIDList pq.Int64Array
	var productIDList pq.Int64Array
	if args.BrandID > 0 {
		brandBindList, _, _ := GetBrandBindList(&ArgsGetBrandBindList{
			Pages: CoreSQL2.ArgsPages{
				Page: 1,
				Max:  30,
				Sort: "id",
				Desc: false,
			},
			OrgID:     args.OrgID,
			BrandID:   args.BrandID,
			CompanyID: -1,
			ProductID: -1,
			IsRemove:  false,
		})
		for _, v := range brandBindList {
			if v.CompanyID > 0 {
				companyIDList = append(companyIDList, v.CompanyID)
			}
			if v.ProductID > 0 {
				productIDList = append(productIDList, v.ProductID)
			}
		}
	}
	if len(brandIDList) > 0 {
		for _, v := range brandIDList {
			vBrandBindList, _, _ := GetBrandBindList(&ArgsGetBrandBindList{
				Pages: CoreSQL2.ArgsPages{
					Page: 1,
					Max:  30,
					Sort: "id",
					Desc: false,
				},
				OrgID:     args.OrgID,
				BrandID:   v,
				CompanyID: -1,
				ProductID: -1,
				IsRemove:  false,
			})
			for _, v2 := range vBrandBindList {
				if v2.CompanyID > 0 {
					companyIDList = append(companyIDList, v2.CompanyID)
				}
				if v2.ProductID > 0 {
					productIDList = append(productIDList, v2.ProductID)
				}
			}
		}
	}
	var rawList []FieldsProduct
	dataCount, err = productDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetSearchQuery([]string{"code"}, args.SearchCode).SetSearchQuery([]string{"title", "title_des", "des"}, args.Search).SetIDQuery("org_id", args.OrgID).SetIDQuery("sort_id", args.SortID).SetIDsQuery("tags", args.Tags).SetIntQuery("pack_type", args.PackType).SetIDQuery("company_id", args.CompanyID).SetIDsQuery("sort_id", sortIDList).SelectList("").ResultAndCount(&rawList)
	if err != nil {
		return
	}
	for _, v := range rawList {
		dataList = append(dataList, getProductByID(v.ID))
	}
	return
}

// GetProductName 获取产品名称
func GetProductName(id int64) string {
	data := getProductByID(id)
	return data.Title
}

// ArgsGetProductByID 获取指定产品参数
type ArgsGetProductByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetProductByID 获取指定产品
func GetProductByID(args *ArgsGetProductByID) (data FieldsProduct, err error) {
	data = getProductByID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsProduct{}
		err = errors.New("no data")
		return
	}
	return
}

func GetProductByIDNoErr(id int64) (data FieldsProduct) {
	data = getProductByID(id)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		data = FieldsProduct{}
		return
	}
	return
}

// ArgsGetProductMore 获取多个产品ID参数
type ArgsGetProductMore struct {
	//ID
	IDs []int64 `json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
}

// GetProductMore 获取多个产品ID
func GetProductMore(args *ArgsGetProductMore) (dataList []FieldsProduct) {
	for _, v := range args.IDs {
		vData, _ := GetProductByID(&ArgsGetProductByID{
			ID:    v,
			OrgID: args.OrgID,
		})
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetProductByCode 获取指定产品参数
type ArgsGetProductByCode struct {
	//组织ID
	// 注意如果给-1则无效，必须给与0或对应组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//商品条形码编码
	Code string `db:"code" json:"code"`
}

// GetProductByCode 获取指定产品
func GetProductByCode(args *ArgsGetProductByCode) (data FieldsProduct, err error) {
	data = getProductByCode(args.OrgID, args.Code)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsProduct{}
		err = errors.New("no data")
		return
	}
	return
}

// GetProductBySN 获取指定产品
func GetProductBySN(orgID int64, sn string) (data FieldsProduct) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM erp_product WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND sn = $2", orgID, sn)
	if err != nil || data.ID < 1 {
		return
	}
	data = getProductByID(data.ID)
	return
}

// getProductByID 获取产品
func getProductByID(id int64) (data FieldsProduct) {
	cacheMark := getProductCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, company_id, company_name, sort_id, tags, sn, code, pin_yin, en_name, manufacturer_name, title, title_des, des, cover_file_ids, expire_hour, weight, size_w, size_h, size_z, pack_type, pack_unit_name, pack_unit, tip_price, tip_tax_price, is_discount, currency, cost_price, tax, tax_cost_price, rebate_price, params FROM erp_product WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheProductTime)
	return
}

func getProductByCode(orgID int64, code string) (data FieldsProduct) {
	cacheMark := getProductCodeCacheMark(orgID, code)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM erp_product WHERE org_id = $1 AND code = $2 AND delete_at < to_timestamp(1000000)", orgID, code)
	if err != nil {
		return
	}
	data = getProductByID(data.ID)
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheProductTime)
	return
}

// checkProductCode 检查是否存在相同的code
func checkProductCode(orgID int64, code string) bool {
	cacheMark := getProductCodeCacheMark(orgID, code)
	var data FieldsProduct
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		if !CoreSQL.CheckTimeHaveData(data.DeleteAt) && orgID == data.OrgID {
			return false
		}
	}
	var count int64
	err := Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM erp_product WHERE org_id = $1 AND code = $2 AND delete_at < to_timestamp(1000000)", orgID, code)
	if err != nil || count < 1 {
		return true
	}
	return false
}
