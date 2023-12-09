package MallCoreMod

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsGetProduct 获取指定商品参数
type ArgsGetProduct struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetProduct 获取指定商品
func GetProduct(args *ArgsGetProduct) (data FieldsCore, err error) {
	data = getProductByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

func GetProductNoErr(id int64, orgID int64) (data FieldsCore) {
	var err error
	data, err = GetProduct(&ArgsGetProduct{
		ID:    id,
		OrgID: orgID,
	})
	if err != nil {
		data = FieldsCore{}
		return
	}
	return
}

// GetProductTop 是否采用追溯的商品获取产品
func GetProductTop(args *ArgsGetProduct) (data FieldsCore, err error) {
	data, err = GetProduct(args)
	if err != nil {
		return
	}
	if data.ParentID < 1 {
		return
	}
	data = getProductByID(data.ParentID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// CheckProductIsVirtual 检查商品是否为虚拟商品
func CheckProductIsVirtual(id int64, orgID int64) (b bool) {
	data := getProductByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) {
		return
	}
	if data.IsVirtual {
		return true
	}
	virtual, b := data.Params.GetValBool("virtual")
	if b && virtual {
		return true
	}
	return
}

// GetProductListByWarehouseProductID 通过仓储产品ID获取商品
func GetProductListByWarehouseProductID(warehouseProductID int64) (dataList []FieldsCore) {
	var rawList []FieldsCore
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM mall_core WHERE delete_at < to_timestamp(1000000) AND warehouse_product_id = $1 AND parent_id = 0", warehouseProductID)
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

// 获取商品信息
func getProductByID(id int64) (data FieldsCore) {
	cacheMark := getProductCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, parent_id, org_id, product_type, is_virtual, sort, publish_at, sort_id, tags, code, title, title_des, des, cover_file_ids, des_files, currency, price_real, price_expire_at, price, integral, integral_price, integral_transport_free, user_sub_price, user_ticket, transport_id, address, warehouse_product_id, weight, count, buy_count, other_options, giving_tickets, params FROM mall_core WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 86400)
	return
}
