package ERPPermanentAssets

import (
	"errors"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"math"
	"time"
)

// ArgsCreateDelete 销毁记录参数
type ArgsCreateDelete struct {
	//创建时间/盘点时间
	CreateAt time.Time `db:"create_at" json:"createAt" check:"defaultTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//增加或减少数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//存放地点
	SavePlace string `db:"save_place" json:"savePlace"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateDelete 销毁记录
func CreateDelete(args *ArgsCreateDelete) (errCode string, err error) {
	productData := getProductByID(args.ProductID)
	if productData.ID < 1 {
		errCode = "err_erp_permanent_assets_product_no_data"
		err = errors.New("product no data")
		return
	}
	if args.Count > productData.Count-productData.UseCount {
		errCode = "err_erp_permanent_assets_product_too_more_delete"
		err = errors.New("no more count")
		return
	}
	targetDeleteCount := productData.Count - args.Count
	if targetDeleteCount < productData.UseCount {
		errCode = "err_erp_permanent_assets_product_too_more_delete_use"
		err = errors.New("no more use count")
		return
	}
	err = createLog(&argsCreateLog{
		CreateAt:     args.CreateAt,
		OrgID:        args.OrgID,
		OrgBindID:    args.OrgBindID,
		ProductID:    args.ProductID,
		Mode:         "delete",
		UseName:      "",
		UseOrgBindID: 0,
		AllPrice:     int64(math.Abs(float64(productData.NowPerPrice * args.Count))),
		PerPrice:     productData.NowPerPrice,
		Count:        args.Count,
		SavePlace:    "",
		Des:          args.Des,
		Params:       args.Params,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	err = updateProductCount(&argsUpdateProductCount{
		ID:    productData.ID,
		Count: targetDeleteCount,
	})
	if err != nil {
		errCode = "err_update"
		return
	}
	return
}
