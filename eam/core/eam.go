package EAMCore

import (
	"fmt"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"time"
)

//获取设备列表

// ArgsGetCore 查看设备详情参数
type ArgsGetCore struct {
}

// GetCore 查看设备详情
func GetCore(args *ArgsGetCore) (data FieldsEAM, err error) {
	return
}

// ArgsGetCoreByCode 通过编码查询设备参数
type ArgsGetCoreByCode struct {
}

// GetCoreByCode 通过编码查询设备
func GetCoreByCode(args *ArgsGetCoreByCode) (data FieldsEAM, err error) {
	return

}

// ArgsCreateCore 新建设备参数
type ArgsCreateCore struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品名称
	ProductName string `db:"product_name" json:"productName" check:"des" min:"1" max:"300"`
	//挂出价格
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//所属分类ID
	CategoryID int64 `db:"category_id" json:"categoryID" check:"id"`
}

// CreateCore 创建设备
func CreateCore(args *ArgsCreateCore) (id int64, err error) {
	return
}

// ArgsUpdateCore 修改设备信息参数
type ArgsUpdateCore struct {
	//使用状态
	// 0: 未使用; 1: 已使用; 2: 已报废; 3: 已闲置; 4 维修中
	Status int `db:"status" json:"status"`
	//当前总金额
	Total int64 `db:"total" json:"total" check:"int64Than0"`
	//单价金额
	Price int64 `db:"price" json:"price" check:"int64Than0"`
	//质保过期时间
	// 根据入库时间+产品质保时间计算
	WarrantyAt time.Time `db:"warranty_at" json:"warrantyAt"`
	//存放位置
	Location string `db:"location" json:"location"`
	//备注
	Remark string `db:"remark" json:"remark"`
}

// UpdateCore 修改设备信息
func UpdateCore(args *ArgsUpdateCore) (err error) {
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
	if args.ID > 0 {
		err = coreDB.Delete().NeedSoft(true).AddWhereID(args.ID).AddWhereOrgID(args.OrgID).ExecNamed(nil)
		if err != nil {
			return
		}
		deleteCoreCache(args.ID)
	} else {
		err = coreDB.Delete().NeedSoft(true).SetWhereAnd("code", args.Code).AddWhereOrgID(args.OrgID).ExecNamed(nil)
	}
	return
}

// getCoreData 获取设备数据
func getCoreData(id int64) (data FieldsEAM) {
	cacheMark := getCoreCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := coreDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "code", "org_id", "product_id", "warehouse_batch_id", "status", "total", "price", "warranty_at", "location", "remark"}).GetByID(id).NeedLimit().Result(&data)
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
