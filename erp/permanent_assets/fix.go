package ERPPermanentAssets

import (
	"errors"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

// ArgsCreateFix 创建维护处理参数
type ArgsCreateFix struct {
	//创建时间/盘点时间
	CreateAt time.Time `db:"create_at" json:"createAt" check:"defaultTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//操作人
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"price"`
	//资产ID
	ProductID int64 `db:"product_id" json:"productID" check:"price"`
	//修正数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//备注
	Des string `db:"des" json:"des" check:"des" min:"0" max:"3000" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateFix 创建维护处理
func CreateFix(args *ArgsCreateFix) (errCode string, err error) {
	productData := getProductByID(args.ProductID)
	if productData.ID < 1 {
		errCode = "err_erp_permanent_assets_return_more"
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
		Mode:         "fix",
		UseName:      "",
		UseOrgBindID: 0,
		AllPrice:     productData.NowPerPrice * args.Count,
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
	return
}
