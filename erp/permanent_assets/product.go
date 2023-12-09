package ERPPermanentAssets

import (
	"errors"
	"fmt"
	CoreCache "gitee.com/weeekj/weeekj_core/v5/core/cache"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	OrgCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetProductList 获取固定资产列表参数
type ArgsGetProductList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//录入操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//计划盘点人
	CheckOrgBindID int64 `db:"check_org_bind_id" json:"checkOrgBindID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//资产条码
	Code string `db:"code" json:"code" check:"mark" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetProductList 获取固定资产列表
func GetProductList(args *ArgsGetProductList) (dataList []FieldsProduct, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.OrgBindID > -1 {
		where = where + " AND org_bind_id = :org_bind_id"
		maps["org_bind_id"] = args.OrgBindID
	}
	if args.CheckOrgBindID > -1 {
		where = where + " AND check_org_bind_id = :check_org_bind_id"
		maps["check_org_bind_id"] = args.CheckOrgBindID
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags && :tags"
		maps["tags"] = args.Tags
	}
	if args.Code != "" {
		where = where + " AND code = :code"
		maps["code"] = args.Code
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "erp_permanent_assets_product"
	var rawList []FieldsProduct
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		fmt.Sprint("SELECT id ", "FROM "+tableName+" WHERE "+where),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getProductByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		vData.Des = ""
		dataList = append(dataList, vData)
	}
	return
}

// GetProduct 获取指定固定资产
func GetProduct(id int64, orgID int64) (data FieldsProduct) {
	data = getProductByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		data = FieldsProduct{}
		return
	}
	return
}

func GetProductName(id int64) (name string) {
	data := getProductByID(id)
	return data.Name
}

// GetProductListBySortID 获取分类下所有产品
func GetProductListBySortID(orgID int64, sortID int64) (dataList []FieldsProduct) {
	where := "org_id = :org_id AND sort_id = :sort_id AND delete_at < to_timestamp(1000000)"
	maps := map[string]any{
		"org_id":  orgID,
		"sort_id": sortID,
	}
	tableName := "erp_permanent_assets_product"
	var rawList []FieldsProduct
	err := CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		fmt.Sprint("SELECT id ", "FROM "+tableName+" WHERE "+where),
		maps,
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getProductByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		vData.Des = ""
		dataList = append(dataList, vData)
	}
	return
}

// GetProductListByCreateAt 获取指定时间范围创建的产品
func GetProductListByCreateAt(orgID int64, startAt time.Time, endAt time.Time) (dataList []FieldsProduct) {
	where := "org_id = :org_id AND create_at >= :start_at AND create_at <= :end_at AND delete_at < to_timestamp(1000000)"
	maps := map[string]any{
		"org_id":   orgID,
		"start_at": startAt,
		"end_at":   endAt,
	}
	tableName := "erp_permanent_assets_product"
	var rawList []FieldsProduct
	err := CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		fmt.Sprint("SELECT id ", "FROM "+tableName+" WHERE "+where),
		maps,
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getProductByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		vData.Des = ""
		dataList = append(dataList, vData)
	}
	return
}

// 获取产品的残值率
func getProductResidualRate(id int64) (residualRate float64) {
	data := getProductByID(id)
	if data.ID < 1 {
		return
	}
	residualRateInt64 := data.Params.GetValInt64NoBool("ERPPermanentAssetsAutoExpireP")
	if residualRateInt64 > 0 {
		residualRate = float64(residualRateInt64) / 10000
	} else {
		globConfigInt64 := OrgCore.Config.GetConfigValInt64NoErr(data.OrgID, "ERPPermanentAssetsAutoExpireP")
		if globConfigInt64 > 0 {
			residualRate = float64(globConfigInt64) / 10000
		}
	}
	if residualRate < 0 {
		residualRate = 0
	}
	if residualRate > 1 {
		residualRate = 1
	}
	return
}

// ArgsCreateProduct 创建固定资产参数
type ArgsCreateProduct struct {
	//创建时间/盘点时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"defaultTime" empty:"true"`
	//使用年限
	UseExpireYear int64 `db:"use_expire_year" json:"useExpireYear" check:"int64Than0" empty:"true"`
	//使用月份限
	UseExpireMonth int64 `db:"use_expire_month" json:"useExpireMonth" check:"int64Than0" empty:"true"`
	//下一次盘点时间
	WaitCheckAt time.Time `db:"wait_check_at" json:"waitCheckAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//录入操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//计划盘点人
	CheckOrgBindID int64 `db:"check_org_bind_id" json:"checkOrgBindID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//资产名称
	Name string `db:"name" json:"name" check:"name"`
	//资产条码
	Code string `db:"code" json:"code" check:"mark" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//购买单价
	BuyPerPrice int64 `db:"buy_per_price" json:"buyPerPrice" check:"price" empty:"true"`
	//购买总价
	BuyAllPrice int64 `db:"buy_all_price" json:"buyAllPrice" check:"price" empty:"true"`
	//当前数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//存放地点
	SavePlace string `db:"save_place" json:"savePlace" check:"des" min:"1" max:"300" empty:"true"`
	//约定使用部门名称
	// 可以使用log借出逻辑，或者在这里直接指定部门名称
	PlanUseOrgGroupName string `db:"plan_use_org_group_name" json:"planUseOrgGroupName" check:"des" min:"1" max:"300" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateProduct 创建固定资产
func CreateProduct(args *ArgsCreateProduct) (err error) {
	//修正参数
	if args.Tags == nil {
		args.Tags = pq.Int64Array{}
	}
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//检索code是否重复？
	if !checkProductCode(args.OrgID, args.Code, 0) {
		err = errors.New("code replace")
		return
	}
	//写入数据
	var newProductID int64
	newProductID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_permanent_assets_product (create_at, expire_at, use_expire_year, use_expire_month, wait_check_at, org_id, org_bind_id, check_org_bind_id, sort_id, tags, name, code, cover_file_id, des, buy_per_price, buy_all_price, now_per_price, now_all_price, count, use_count, save_place, plan_use_org_group_name, params) VALUES (:create_at, :expire_at, :use_expire_year, :use_expire_month, :wait_check_at, :org_id, :org_bind_id, :check_org_bind_id, :sort_id, :tags, :name, :code, :cover_file_id, :des, :buy_per_price, :buy_all_price, :buy_per_price, :buy_all_price, :count, 0, :save_place, :plan_use_org_group_name, :params)", args)
	if err != nil {
		return
	}
	//创建日志
	err = createLog(&argsCreateLog{
		CreateAt:     args.CreateAt,
		OrgID:        args.OrgID,
		OrgBindID:    args.OrgBindID,
		ProductID:    newProductID,
		Mode:         "in",
		UseName:      "",
		UseOrgBindID: 0,
		AllPrice:     args.BuyAllPrice,
		PerPrice:     args.BuyPerPrice,
		Count:        args.Count,
		SavePlace:    args.SavePlace,
		Des:          args.Des,
		Params:       nil,
	})
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateProduct 修改固定资产参数
type ArgsUpdateProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"defaultTime" empty:"true"`
	//使用年限
	UseExpireYear int64 `db:"use_expire_year" json:"useExpireYear" check:"int64Than0" empty:"true"`
	//使用月份限
	UseExpireMonth int64 `db:"use_expire_month" json:"useExpireMonth" check:"int64Than0" empty:"true"`
	//下一次盘点时间
	WaitCheckAt time.Time `db:"wait_check_at" json:"waitCheckAt" check:"defaultTime"`
	//计划盘点人
	CheckOrgBindID int64 `db:"check_org_bind_id" json:"checkOrgBindID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签ID列
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//资产名称
	Name string `db:"name" json:"name" check:"name"`
	//资产条码
	Code string `db:"code" json:"code" check:"mark" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//购买单价
	BuyPerPrice int64 `db:"buy_per_price" json:"buyPerPrice" check:"price" empty:"true"`
	//购买总价
	BuyAllPrice int64 `db:"buy_all_price" json:"buyAllPrice" check:"price" empty:"true"`
	//当前资产单价
	NowPerPrice int64 `db:"now_per_price" json:"nowPerPrice" check:"price" empty:"true"`
	//当前总价值
	NowAllPrice int64 `db:"now_all_price" json:"nowAllPrice" check:"price" empty:"true"`
	//当前数量
	Count int64 `db:"count" json:"count" check:"int64Than0" empty:"true"`
	//存放地点
	SavePlace string `db:"save_place" json:"savePlace" check:"des" min:"1" max:"300" empty:"true"`
	//约定使用部门名称
	// 可以使用log借出逻辑，或者在这里直接指定部门名称
	PlanUseOrgGroupName string `db:"plan_use_org_group_name" json:"planUseOrgGroupName" check:"des" min:"1" max:"300" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateProduct 修改固定资产
func UpdateProduct(args *ArgsUpdateProduct) (err error) {
	//修正参数
	if args.Tags == nil {
		args.Tags = pq.Int64Array{}
	}
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//检索code是否重复？
	if !checkProductCode(args.OrgID, args.Code, args.ID) {
		err = errors.New("code replace")
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_permanent_assets_product SET update_at = NOW(), expire_at = :expire_at, use_expire_year = :use_expire_year, use_expire_month = :use_expire_month, wait_check_at = :wait_check_at, check_org_bind_id = :check_org_bind_id, sort_id = :sort_id, tags = :tags, name = :name, code = :code, cover_file_id = :cover_file_id, des = :des, buy_per_price = :buy_per_price, buy_all_price = :buy_all_price, now_per_price = :now_per_price, now_all_price = :now_all_price, count = :count, save_place = :save_place, plan_use_org_group_name = :plan_use_org_group_name, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteProductCache(args.ID)
	//反馈
	return
}

// ArgsDeleteProduct 删除固定资产参数
type ArgsDeleteProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteProduct 删除固定资产
func DeleteProduct(args *ArgsDeleteProduct) (err error) {
	//删除数据
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "erp_permanent_assets_product", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteProductCache(args.ID)
	//反馈
	return
}

func checkProductCode(orgID int64, code string, nowID int64) (b bool) {
	if code == "" {
		return true
	}
	var data FieldsProduct
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM erp_permanent_assets_product WHERE org_id = $1 AND code = $2 AND delete_at < to_timestamp(1000000)", orgID, code)
	if err != nil || data.ID < 1 {
		return true
	}
	if nowID > 0 && data.ID == nowID {
		return true
	}
	return
}

func getProductByID(id int64) (data FieldsProduct) {
	cacheMark := getProductCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, use_expire_year, wait_check_at, org_id, org_bind_id, check_org_bind_id, sort_id, tags, name, code, cover_file_id, des, buy_per_price, buy_all_price, now_per_price, now_all_price, count, use_count, save_place, params FROM erp_permanent_assets_product WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, CoreCache.CacheTime2Day)
	return
}

func getProductCacheMark(id int64) string {
	return fmt.Sprint("erp:permanent:assets:product:id:", id)
}

func deleteProductCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getProductCacheMark(id))
}

// argsUpdateProductCheck 修改固定资产参数
type argsUpdateProductCheck struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//下一次盘点时间
	WaitCheckAt time.Time `db:"wait_check_at" json:"waitCheckAt"`
	//当前资产单价
	NowPerPrice int64 `db:"now_per_price" json:"nowPerPrice"`
	//当前总价值
	NowAllPrice int64 `db:"now_all_price" json:"nowAllPrice"`
	//当前数量
	Count int64 `db:"count" json:"count"`
	//存放地点
	SavePlace string `db:"save_place" json:"savePlace"`
}

// updateProductCheck 修改固定资产，检查专用
func updateProductCheck(args *argsUpdateProductCheck) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_permanent_assets_product SET update_at = NOW(), wait_check_at = :wait_check_at, now_per_price = :now_per_price, now_all_price = :now_all_price, count = :count, save_place = :save_place WHERE id = :id", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteProductCache(args.ID)
	//反馈
	return
}

// argsUpdateProductCount 修改固定资产参数
type argsUpdateProductCount struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//当前数量
	Count int64 `db:"count" json:"count"`
}

// updateProductCount 修改固定资产，变更数量
func updateProductCount(args *argsUpdateProductCount) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_permanent_assets_product SET update_at = NOW(), count = :count WHERE id = :id", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteProductCache(args.ID)
	//如果当前数量为0，则自动删除产品
	if args.Count < 1 {
		_ = DeleteProduct(&ArgsDeleteProduct{
			ID:    args.ID,
			OrgID: -1,
		})
	}
	//反馈
	return
}

// argsUpdateProductUseCount 修改固定资产参数
type argsUpdateProductUseCount struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//当前数量
	Count int64 `db:"count" json:"count"`
}

// updateProductUseCount 修改固定资产，变更使用数量
func updateProductUseCount(args *argsUpdateProductUseCount) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE erp_permanent_assets_product SET update_at = NOW(), use_count = :count WHERE id = :id", args)
	if err != nil {
		return
	}
	//清理缓冲
	deleteProductCache(args.ID)
	//反馈
	return
}
