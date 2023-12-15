package ERPPermanentAssets

import (
	"errors"
	"fmt"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetUseList 获取使用记录列表参数
type ArgsGetUseList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//操作主体描述
	UseName string `db:"use_name" json:"useName" check:"name" empty:"true"`
	//实际使用人
	UseOrgBindID int64 `db:"use_org_bind_id" json:"useOrgBindID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetUseList 获取使用记录列表
func GetUseList(args *ArgsGetUseList) (dataList []FieldsUse, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.OrgBindID > -1 {
		where = where + " AND org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.ProductID > -1 {
		where = where + " AND product_id = :product_id"
		maps["product_id"] = args.ProductID
	}
	if args.UseName != "" {
		where = where + " AND use_name = :use_name"
		maps["use_name"] = args.UseName
	}
	if args.UseOrgBindID > -1 {
		where = where + " AND use_org_bind_id = :use_org_bind_id"
		maps["use_org_bind_id"] = args.UseOrgBindID
	}
	if args.Search != "" {
		where = where + " AND (des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_permanent_assets_use"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		fmt.Sprint("SELECT id, create_at, return_at, org_id, org_bind_id, product_id, use_name, use_org_bind_id, count, take_count, return_count, des, params ", "FROM "+tableName+" WHERE "+where),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "return_at", "count", "take_count", "return_count"},
	)
	if err != nil {
		return
	}
	for k, _ := range dataList {
		dataList[k].Des = ""
	}
	return
}

// GetUse 查询使用记录
func GetUse(id int64, orgID int64) (data FieldsUse) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, return_at, org_id, org_bind_id, product_id, use_name, use_org_bind_id, count, take_count, return_count, des, params FROM erp_permanent_assets_use WHERE id = $1 AND ($2 < 0 OR org_id = $2)", id, orgID)
	if err != nil {
		return
	}
	return
}

// ArgsCreateUse 领域使用参数
type ArgsCreateUse struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//实际使用主体（部门）
	UseName string `db:"use_name" json:"useName" check:"name" empty:"true"`
	//实际使用人
	UseOrgBindID int64 `db:"use_org_bind_id" json:"useOrgBindID"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"0" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateUse 领域使用
func CreateUse(args *ArgsCreateUse) (errCode string, err error) {
	//修正参数
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//获取产品
	productData := getProductByID(args.ProductID)
	if productData.ID < 1 {
		errCode = "err_erp_permanent_assets_product_no_data"
		err = errors.New("product no data")
		return
	}
	if productData.Count-productData.UseCount < args.Count {
		errCode = "err_erp_permanent_assets_product_no_more"
		err = errors.New("product no more")
		return
	}
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_permanent_assets_use (create_at, return_at, org_id, org_bind_id, product_id, use_name, use_org_bind_id, count, take_count, return_count, des, params) VALUES (:create_at, to_timestamp(0), :org_id, :org_bind_id, :product_id, :use_name, :use_org_bind_id, :count, :count, 0, :des, :params)", args)
	if err != nil {
		errCode = "err_insert"
		return
	}
	//记录日志
	err = createLog(&argsCreateLog{
		CreateAt:     args.CreateAt,
		OrgID:        args.OrgID,
		OrgBindID:    args.OrgBindID,
		ProductID:    args.ProductID,
		Mode:         "take",
		UseName:      args.UseName,
		UseOrgBindID: args.UseOrgBindID,
		AllPrice:     args.Count * productData.NowPerPrice,
		PerPrice:     productData.NowPerPrice,
		Count:        args.Count,
		SavePlace:    "",
		Des:          args.Des,
		Params:       nil,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//更新产品数量
	err = updateProductUseCount(&argsUpdateProductUseCount{
		ID:    productData.ID,
		Count: productData.UseCount + args.Count,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//反馈
	return
}

// ArgsCreateUseReturn 归还使用参数
type ArgsCreateUseReturn struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//归还时间
	ReturnAt time.Time `db:"return_at" json:"returnAt" check:"defaultTime"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//归还数量
	ReturnCount int64 `db:"return_count" json:"returnCount" check:"int64Than0"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"0" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateUseReturn 归还使用
func CreateUseReturn(args *ArgsCreateUseReturn) (errCode string, err error) {
	//修正参数
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//获取日志
	useData := GetUse(args.ID, args.OrgID)
	if useData.ID < 1 {
		errCode = "err_no_data"
		err = errors.New("use no data")
		return
	}
	if args.ReturnCount > useData.TakeCount {
		errCode = "err_erp_permanent_assets_return_more"
		err = errors.New("use no data")
		return
	}
	//获取产品
	productData := getProductByID(useData.ProductID)
	if productData.ID < 1 {
		errCode = "err_erp_permanent_assets_product_no_data"
		err = errors.New("product no data")
		return
	}
	if productData.UseCount < args.ReturnCount {
		errCode = "err_erp_permanent_assets_product_too_more"
		err = errors.New("product return count more")
		return
	}
	//写入数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE erp_permanent_assets_use SET return_at = :return_at, return_count = :return_count, count = :count, params = :params WHERE id = :id AND org_id = :org_id", map[string]interface{}{
		"id":           args.ID,
		"org_id":       args.OrgID,
		"return_at":    args.ReturnAt,
		"return_count": args.ReturnCount,
		"count":        useData.TakeCount - args.ReturnCount,
		"params":       args.Params,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	//记录日志
	err = createLog(&argsCreateLog{
		CreateAt:     args.ReturnAt,
		OrgID:        args.OrgID,
		OrgBindID:    args.OrgBindID,
		ProductID:    useData.ProductID,
		Mode:         "return",
		UseName:      useData.UseName,
		UseOrgBindID: useData.UseOrgBindID,
		AllPrice:     args.ReturnCount * productData.NowPerPrice,
		PerPrice:     productData.NowPerPrice,
		Count:        args.ReturnCount,
		SavePlace:    "",
		Des:          args.Des,
		Params:       nil,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//更新产品数量
	err = updateProductUseCount(&argsUpdateProductUseCount{
		ID:    productData.ID,
		Count: productData.UseCount - args.ReturnCount,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//反馈
	return
}
