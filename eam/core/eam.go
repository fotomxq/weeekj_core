package EAMCore

import (
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	ERPProduct "github.com/fotomxq/weeekj_core/v5/erp/product"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

// ArgsGetCoreList 获取设备列表参数
type ArgsGetCoreList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//关联库存批次ID
	WarehouseBatchID int64 `db:"warehouse_batch_id" json:"warehouseBatchID" check:"id" empty:"true"`
	//采购订单来源
	ERPPurchaseOrderID int64 `db:"erp_purchase_order_id" json:"erpPurchaseOrderID" check:"id" empty:"true"`
	//使用状态
	// 0: 未使用; 1: 已使用; 2: 已报废; 3: 已闲置; 4 维修中
	Status int `db:"status" json:"status"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetCoreList 获取设备列表
func GetCoreList(args *ArgsGetCoreList) (dataList []FieldsEAM, dataCount int64, err error) {
	dataCount, err = coreDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "price"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetSearchQuery([]string{"code"}, args.Code).SetIDQuery("warehouse_batch_id", args.WarehouseBatchID).SetIDQuery("erp_purchase_order_id", args.ERPPurchaseOrderID).SetIntQuery("status", args.Status).SetSearchQuery([]string{"org_name", "product_name", "location", "remark"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getCoreData(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetCore 查看设备详情参数
type ArgsGetCore struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetCore 查看设备详情
func GetCore(args *ArgsGetCore) (data FieldsEAM, err error) {
	data = getCoreData(args.ID)
	if data.ID < 1 {
		err = fmt.Errorf("can not find core by id")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = fmt.Errorf("can not find core by id and org_id")
		return
	}
	return
}

// ArgsGetCoreByCode 通过编码查询设备参数
type ArgsGetCoreByCode struct {
	//编码
	// ID二选一操作
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetCoreByCode 通过编码查询设备
func GetCoreByCode(args *ArgsGetCoreByCode) (data FieldsEAM, err error) {
	//获取设备信息
	err = coreDB.Get().SetFieldsOne([]string{"id"}).GetByCodeAndOrgID(args.Code, args.OrgID).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	if data.ID < 1 {
		err = fmt.Errorf("can not find core by code")
		return
	}
	data = getCoreData(data.ID)
	if data.ID < 1 {
		err = fmt.Errorf("can not find core by id")
		return
	}
	return
}

// ArgsCreateCore 新建设备参数
type ArgsCreateCore struct {
	//编码
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//关联库存批次ID
	WarehouseBatchID int64 `db:"warehouse_batch_id" json:"warehouseBatchID" check:"id" empty:"true"`
	//采购订单来源
	ERPPurchaseOrderID int64 `db:"erp_purchase_order_id" json:"erpPurchaseOrderID" check:"id" empty:"true"`
	//使用状态
	// 0: 未使用; 1: 已使用; 2: 已报废; 3: 已闲置; 4 维修中
	Status int `db:"status" json:"status"`
	//单价金额
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//质保过期时间
	// 根据入库时间+产品质保时间计算
	WarrantyAt time.Time `db:"warranty_at" json:"warrantyAt"`
	//存放位置
	Location string `db:"location" json:"location" check:"des" min:"1" max:"600" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"3000" empty:"true"`
}

// CreateCore 创建设备
func CreateCore(args *ArgsCreateCore) (id int64, err error) {
	//获取组织信息
	var orgName string
	if args.OrgID > 0 {
		orgData := OrgCore.GetOrgByID(args.OrgID)
		orgName = orgData.Name
	}
	//获取商品信息
	productData := ERPProduct.GetProductByIDNoErr(args.ProductID)
	//创建数据
	id, err = coreDB.Insert().SetFields([]string{"code", "org_id", "org_name", "product_id", "product_name", "warehouse_batch_id", "erp_purchase_order_id", "status", "price", "warranty_at", "location", "remark"}).Add(map[string]any{
		"code":                  args.Code,
		"org_id":                args.OrgID,
		"org_name":              orgName,
		"product_id":            args.ProductID,
		"product_name":          productData.Title,
		"warehouse_batch_id":    args.WarehouseBatchID,
		"erp_purchase_order_id": args.ERPPurchaseOrderID,
		"status":                args.Status,
		"price":                 args.Price,
		"warranty_at":           args.WarrantyAt,
		"location":              args.Location,
		"remark":                args.Remark,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateCore 修改设备信息参数
type ArgsUpdateCore struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//编码
	// ID二选一操作
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//使用状态
	// 0: 未使用; 1: 已使用; 2: 已报废; 3: 已闲置; 4 维修中
	Status int `db:"status" json:"status"`
	//单价金额
	Price int64 `db:"price" json:"price" check:"int64Than0" empty:"true"`
	//质保过期时间
	// 根据入库时间+产品质保时间计算
	WarrantyAt time.Time `db:"warranty_at" json:"warrantyAt"`
	//存放位置
	Location string `db:"location" json:"location" check:"des" min:"1" max:"600" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"3000" empty:"true"`
}

// UpdateCore 修改设备信息
func UpdateCore(args *ArgsUpdateCore) (err error) {
	if args.ID < 1 && args.Code == "" {
		err = fmt.Errorf("id and code can not be empty")
		return
	}
	if args.ID < 1 && args.Code != "" {
		//获取设备信息
		data, _ := GetCoreByCode(&ArgsGetCoreByCode{
			Code:  args.Code,
			OrgID: args.OrgID,
		})
		if data.ID < 1 {
			err = fmt.Errorf("can not find core by code")
			return
		}
		args.ID = data.ID
	}
	err = coreDB.Update().SetFields([]string{"status", "price", "warranty_at", "location", "remark"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]any{
		"status":      args.Status,
		"price":       args.Price,
		"warranty_at": args.WarrantyAt,
		"location":    args.Location,
		"remark":      args.Remark,
	})
	if err != nil {
		return
	}
	deleteCoreCache(args.ID)
	//反馈
	return
}

// ArgsUpdateCoreStatus 更新设备状态参数
type ArgsUpdateCoreStatus struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//编码
	// ID二选一操作
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//使用状态
	// 0: 未使用; 1: 已使用; 2: 已报废; 3: 已闲置; 4 维修中
	Status int `db:"status" json:"status"`
}

// UpdateCoreStatus 更新设备状态
func UpdateCoreStatus(args *ArgsUpdateCoreStatus) (err error) {
	if args.ID < 1 && args.Code == "" {
		err = fmt.Errorf("id and code can not be empty")
		return
	}
	if args.ID < 1 && args.Code != "" {
		//获取设备信息
		data, _ := GetCoreByCode(&ArgsGetCoreByCode{
			Code:  args.Code,
			OrgID: args.OrgID,
		})
		if data.ID < 1 {
			err = fmt.Errorf("can not find core by code")
			return
		}
		args.ID = data.ID
	}
	err = coreDB.Update().SetFields([]string{"status"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]any{
		"status": args.Status,
	})
	if err != nil {
		return
	}
	deleteCoreCache(args.ID)
	//反馈
	return
}

// ArgsDeleteCore 删除设备参数
type ArgsDeleteCore struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id" empty:"true"`
	//编码
	// ID二选一操作
	Code string `db:"code" json:"code" check:"des" min:"1" max:"50" empty:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteCore 删除设备
func DeleteCore(args *ArgsDeleteCore) (err error) {
	if args.ID < 1 && args.Code == "" {
		err = fmt.Errorf("id and code can not be empty")
		return
	}
	if args.ID < 1 && args.Code != "" {
		//获取设备信息
		data, _ := GetCoreByCode(&ArgsGetCoreByCode{
			Code:  args.Code,
			OrgID: args.OrgID,
		})
		if data.ID < 1 {
			err = fmt.Errorf("can not find core by code")
			return
		}
		args.ID = data.ID
	}
	err = coreDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).ExecNamed(nil)
	if err != nil {
		return
	}
	deleteCoreCache(args.ID)
	//反馈
	return
}

// getCoreData 获取设备数据
func getCoreData(id int64) (data FieldsEAM) {
	cacheMark := getCoreCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := coreDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "code", "org_id", "org_name", "product_id", "product_name", "warehouse_batch_id", "erp_purchase_order_id", "status", "price", "warranty_at", "location", "remark"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheCoreTime)
	return
}

// 缓冲
func getCoreCacheMark(id int64) string {
	return fmt.Sprint("eam:core.id.", id)
}

func deleteCoreCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getCoreCacheMark(id))
}
