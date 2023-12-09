package ERPPermanentAssets

import (
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetLogList 获取清查列表参数
type ArgsGetLogList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID" check:"id" empty:"true"`
	//模式
	// in 入库; take 领取使用; return 归还入库; check 清查库存; fix 维护; delete 销毁
	Mode string `db:"mode" json:"mode" check:"mark" empty:"true"`
	//操作主体描述
	UseName string `db:"use_name" json:"useName" check:"name" empty:"true"`
	//实际使用人
	UseOrgBindID int64 `db:"use_org_bind_id" json:"useOrgBindID" check:"id" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetLogList 获取清查列表
func GetLogList(args *ArgsGetLogList) (dataList []FieldsLog, dataCount int64, err error) {
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
	if args.Mode != "" {
		where = where + " AND mode = :mode"
		maps["mode"] = args.Mode
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
	tableName := "erp_permanent_assets_log"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		fmt.Sprint("SELECT id, create_at, org_id, org_bind_id, product_id, mode, use_name, use_org_bind_id, all_price, per_price, count, save_place, des, params ", "FROM "+tableName+" WHERE "+where),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	for k := 0; k < len(dataList); k++ {
		dataList[k].Des = ""
	}
	return
}

// GetLog 获取清查详情
func GetLog(id int64, orgID int64) (data FieldsLog) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, product_id, mode, use_name, use_org_bind_id, all_price, per_price, count, save_place, des, params FROM erp_permanent_assets_log WHERE id = $1 AND ($2 < 0 OR org_id = $2)", id, orgID)
	if err != nil {
		return
	}
	return
}

// getLogLastByProduct 获取指定产品的最后一条日志
func getLogLastByProduct(productID int64, mode string) (data FieldsLog) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, product_id, mode, use_name, use_org_bind_id, all_price, per_price, count, save_place, des, params FROM erp_permanent_assets_log WHERE product_id = $1 AND mode = $2 ORDER BY create_at DESC LIMIT 1", productID, mode)
	if err != nil {
		return
	}
	return
}

// getLogFirstInBetween 获取指定时间段最初的数据
func getLogFirstInBetween(productID int64, mode string, startAt, endAt time.Time) (data FieldsLog) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, product_id, mode, use_name, use_org_bind_id, all_price, per_price, count, save_place, des, params FROM erp_permanent_assets_log WHERE product_id = $1 AND ($2 = '' OR mode = $2) AND create_at >= $3 AND create_at <= $4 ORDER BY create_at LIMIT 1", productID, mode, startAt, endAt)
	if err != nil {
		return
	}
	return
}

// getLogLastInBetween 获取指定时间段最后的数据
func getLogLastInBetween(productID int64, mode string, startAt, endAt time.Time) (data FieldsLog) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, product_id, mode, use_name, use_org_bind_id, all_price, per_price, count, save_place, des, params FROM erp_permanent_assets_log WHERE product_id = $1 AND ($2 = '' OR mode = $2) AND create_at >= $3 AND create_at <= $4 ORDER BY create_at DESC LIMIT 1", productID, mode, startAt, endAt)
	if err != nil {
		return
	}
	return
}

// getLogLastIDByProductAndBefore 获取指定产品的记录ID上一条记录
func getLogLastIDByProductAndBefore(productID int64, mode string, logID int64) (data FieldsLog) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, org_bind_id, product_id, mode, use_name, use_org_bind_id, all_price, per_price, count, save_place, des, params FROM erp_permanent_assets_log WHERE product_id = $1 AND mode = $2 AND id < $3 ORDER BY create_at DESC LIMIT 1", productID, mode, logID)
	if err != nil {
		return
	}
	return
}

// getLogSUMInBetween 获取指定时间范围的总和
func getLogSUMInBetween(productID int64, mode string, startAt, endAt time.Time) (data FieldsLog) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT SUM(all_price) AS all_price, SUM(per_price) AS per_price, SUM(count) AS count FROM erp_permanent_assets_log WHERE product_id = $1 AND ($2 = '' OR mode = $2) AND create_at >= $3 AND create_at <= $4", productID, mode, startAt, endAt)
	if err != nil {
		return
	}
	return
}

// 获取指定时间段的所有日志
func getLogListByOrgAndMode(orgID int64, mode string, startAt, endAt time.Time) (dataList []FieldsLog) {
	err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, org_id, org_bind_id, product_id, mode, use_name, use_org_bind_id, all_price, per_price, count, save_place, des, params FROM erp_permanent_assets_log WHERE org_id = $1 AND ($2 = '' OR mode = $2) AND create_at >= $3 AND create_at <= $4", orgID, mode, startAt, endAt)
	if err != nil {
		return
	}
	return
}

// argsCreateLog 创建新的清查参数
type argsCreateLog struct {
	//创建时间/盘点时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//模式
	// in 入库; take 领取使用; return 归还入库; check 清查库存; fix 维护; delete 销毁
	Mode string `db:"mode" json:"mode"`
	//操作主体描述
	UseName string `db:"use_name" json:"useName" check:"name" empty:"true"`
	//实际使用人
	UseOrgBindID int64 `db:"use_org_bind_id" json:"useOrgBindID" check:"id" empty:"true"`
	//处置后总价值
	AllPrice int64 `db:"all_price" json:"allPrice" check:"price"`
	//处置后资产单价
	PerPrice int64 `db:"per_price" json:"perPrice" check:"price"`
	//增加或减少数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//存放地点
	SavePlace string `db:"save_place" json:"savePlace"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// createLog 创建新的清查
func createLog(args *argsCreateLog) (err error) {
	//修正参数
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_permanent_assets_log (create_at, org_id, org_bind_id, product_id, mode, use_name, use_org_bind_id, all_price, per_price, count, save_place, des, params) VALUES (:create_at, :org_id, :org_bind_id, :product_id, :mode, :use_name, :use_org_bind_id, :all_price, :per_price, :count, :save_place, :des, :params)", args)
	if err != nil {
		return
	}
	//反馈
	return
}
