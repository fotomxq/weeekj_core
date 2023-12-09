package ERPWarehouse

import (
	"errors"
	AnalysisAny2 "gitee.com/weeekj/weeekj_core/v5/analysis/any2"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	ERPProduct "gitee.com/weeekj/weeekj_core/v5/erp/product"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"math"
	"time"
)

// argsAppendLog 创建日志参数
type argsAppendLog struct {
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//SN，商户下唯一，可注入其他外部系统SN
	SN string `db:"sn" json:"sn"`
	//动作类型
	// in 入库; out 出库; move_in 移动入库; move_out 移动出库
	Action string `db:"action" json:"action"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	UserID    int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id" empty:"true"`
	//所属仓库
	WarehouseID int64 `db:"warehouse_id" json:"warehouseID" check:"id"`
	//区域
	AreaID int64 `db:"area_id" json:"areaID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//变动数量
	Count int64 `db:"count" json:"count" check:"int64Than0" empty:"true"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"600" empty:"true"`
	//附加数据，可选，如果不存在将从产品资料获取
	// 过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	// 变动时产品价格
	PerPrice int64 `db:"per_price" json:"perPrice" check:"price" empty:"true"`
}

// appendLog 创建日志
func appendLog(args *argsAppendLog) (err error) {
	//获取产品数据
	if !CoreSQL.CheckTimeHaveData(args.ExpireAt) || args.PerPrice < 1 {
		productData := ERPProduct.GetProductByIDNoErr(args.ProductID)
		productCompanyData := ERPProduct.GetProductCompany(args.OrgID, args.ProductID, productData.CompanyID)
		if productData.ExpireHour > 0 && !CoreSQL.CheckTimeHaveData(args.ExpireAt) {
			args.ExpireAt = CoreFilter.GetNowTimeCarbon().AddHours(productData.ExpireHour).Time
		}
		if args.PerPrice < 1 {
			args.PerPrice = productCompanyData.TaxCostPrice
		}
	}
	//修正sn
	if args.SN == "" {
		args.SN = CoreFilter.GetRandStr4(30)
	}
	//检查sn
	if !checkLogSN(args.OrgID, args.SN) {
		err = errors.New("sn error")
		return
	}
	//创建数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO erp_warehouse_log (create_at, sn, expire_at, action, org_id, user_id, org_bind_id, warehouse_id, area_id, product_id, count, per_price, des) VALUES (:create_at,:sn,:expire_at,:action,:org_id, :user_id, :org_bind_id, :warehouse_id, :area_id, :product_id, :count, :per_price, :des)", map[string]interface{}{
		"create_at":    args.CreateAt,
		"sn":           args.SN,
		"expire_at":    args.ExpireAt,
		"action":       args.Action,
		"org_id":       args.OrgID,
		"user_id":      args.UserID,
		"org_bind_id":  args.OrgBindID,
		"warehouse_id": args.WarehouseID,
		"area_id":      args.AreaID,
		"product_id":   args.ProductID,
		"count":        args.Count,
		"per_price":    args.PerPrice,
		"des":          args.Des,
	})
	if err != nil {
		return
	}
	//统计支持
	switch args.Action {
	case "in":
		AnalysisAny2.AppendData("add", "erp_warehouse_store_product_price", CoreFilter.GetNowTime(), args.OrgID, 0, 0, 0, 0, args.PerPrice*args.Count)
		AnalysisAny2.AppendData("add", "erp_warehouse_store_count", CoreFilter.GetNowTime(), args.OrgID, 0, 0, 0, 0, args.Count)
		updateAnalysis(args.OrgID)
	case "out":
		AnalysisAny2.AppendData("reduce", "erp_warehouse_store_product_price", CoreFilter.GetNowTime(), args.OrgID, 0, 0, 0, 0, int64(math.Abs(float64(args.PerPrice*args.Count))))
		AnalysisAny2.AppendData("reduce", "erp_warehouse_store_count", CoreFilter.GetNowTime(), args.OrgID, 0, 0, 0, 0, int64(math.Abs(float64(args.Count))))
		updateAnalysis(args.OrgID)
	}
	//反馈
	return
}

func checkLogSN(orgID int64, sn string) bool {
	var data FieldsLog
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM erp_warehouse_log WHERE org_id = $1 AND sn = $2", orgID, sn)
	if err != nil || data.ID < 1 {
		return true
	}
	return false
}
