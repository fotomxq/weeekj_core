package ERPPermanentAssets

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	"time"
)

// ArgsCreateCheck 创建清查记录参数
type ArgsCreateCheck struct {
	//创建时间/盘点时间
	CreateAt time.Time `db:"create_at" json:"createAt" check:"defaultTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//处置后总价值
	AllPrice int64 `db:"all_price" json:"allPrice" check:"price"`
	//处置后资产单价
	PerPrice int64 `db:"per_price" json:"perPrice" check:"price"`
	//修正数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//存放地点
	SavePlace string `db:"save_place" json:"savePlace"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateCheck 创建清查记录
func CreateCheck(args *ArgsCreateCheck) (errCode string, err error) {
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	productData := getProductByID(args.ProductID)
	if productData.ID < 1 {
		errCode = "err_erp_permanent_assets_product_no_data"
		err = errors.New("product no data")
		return
	}
	if args.Count > productData.Count {
		errCode = "err_erp_permanent_assets_product_too_more_use_count"
		err = errors.New("product too more than use count")
		return
	}
	err = createLog(&argsCreateLog{
		CreateAt:     args.CreateAt,
		OrgID:        args.OrgID,
		OrgBindID:    args.OrgBindID,
		ProductID:    args.ProductID,
		Mode:         "check",
		UseName:      "",
		UseOrgBindID: 0,
		AllPrice:     args.AllPrice,
		PerPrice:     args.PerPrice,
		Count:        args.Count,
		SavePlace:    "",
		Des:          args.Des,
		Params:       args.Params,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	erpPermanentAssetsCheckNext := OrgCore.Config.GetConfigValIntNoErr(args.OrgID, "ERPPermanentAssetsCheckNext")
	if erpPermanentAssetsCheckNext < 1 {
		erpPermanentAssetsCheckNext = 30
	}
	err = updateProductCheck(&argsUpdateProductCheck{
		ID:          productData.ID,
		WaitCheckAt: CoreFilter.GetNowTimeCarbon().AddDays(erpPermanentAssetsCheckNext).Time,
		NowPerPrice: args.PerPrice,
		NowAllPrice: args.AllPrice,
		Count:       args.Count,
		SavePlace:   args.SavePlace,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	return
}
