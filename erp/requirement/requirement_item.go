package ERPRequirement

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetRequirementItemList 获取RequirementItem列表参数
type ArgsGetRequirementItemList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//关联头ID
	RequisitionID int64 `db:"requisition_id" json:"requisitionID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetRequirementItemList 获取RequirementItem列表
func GetRequirementItemList(args *ArgsGetRequirementItemList) (dataList []FieldsRequisitionItem, dataCount int64, err error) {
	dataCount, err = requirementItemDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("requisition_id", args.RequisitionID).SetIDQuery("product_id", args.ProductID).SetIDQuery("company_id", args.CompanyID).SetSearchQuery([]string{"companyName"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getRequirementItemByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetRequirementItemByID 获取RequirementItem数据包参数
type ArgsGetRequirementItemByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetRequirementItemByID 获取RequirementItem数
func GetRequirementItemByID(args *ArgsGetRequirementItemByID) (data FieldsRequisitionItem, err error) {
	data = getRequirementItemByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreateRequirementItem 创建RequirementItem参数
type ArgsCreateRequirementItem struct {
	//关联头ID
	RequisitionID int64 `db:"requisition_id" json:"requisitionID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品价格
	Price int64 `db:"price" json:"price" check:"price"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
}

// CreateRequirementItem 创建RequirementItem
func CreateRequirementItem(args *ArgsCreateRequirementItem) (id int64, err error) {
	//创建数据
	id, err = requirementItemDB.Insert().SetFields([]string{"requisition_id", "product_id", "price", "count", "company_id", "company_name"}).Add(map[string]any{
		"requisition_id": args.RequisitionID,
		"product_id":     args.ProductID,
		"price":          args.Price,
		"count":          args.Count,
		"company_id":     args.CompanyID,
		"company_name":   args.CompanyName,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateRequirementItem 修改RequirementItem参数
type ArgsUpdateRequirementItem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//关联头ID
	RequisitionID int64 `db:"requisition_id" json:"requisitionID" check:"id"`
	//产品ID
	ProductID int64 `db:"product_id" json:"productID" check:"id"`
	//产品价格
	Price int64 `db:"price" json:"price" check:"price"`
	//数量
	Count int64 `db:"count" json:"count" check:"int64Than0"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
}

// UpdateRequirementItem 修改RequirementItem
func UpdateRequirementItem(args *ArgsUpdateRequirementItem) (err error) {
	//更新数据
	err = requirementItemDB.Update().SetFields([]string{"requisition_id", "product_id", "price", "count", "remark", "company_id", "company_name"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"requisition_id": args.RequisitionID,
		"product_id":     args.ProductID,
		"price":          args.Price,
		"count":          args.Count,
		"company_id":     args.CompanyID,
		"company_name":   args.CompanyName,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteRequirementItemCache(args.ID)
	//反馈
	return
}

// ArgsDeleteRequirementItem 删除RequirementItem参数
type ArgsDeleteRequirementItem struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteRequirementItem 删除RequirementItem
func DeleteRequirementItem(args *ArgsDeleteRequirementItem) (err error) {
	//删除数据
	err = requirementItemDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteRequirementItemCache(args.ID)
	//反馈
	return
}

// getRequirementItemByID 通过ID获取RequirementItem数据包
func getRequirementItemByID(id int64) (data FieldsRequisitionItem) {
	cacheMark := getRequirementItemCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := requirementItemDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "requisition_id", "product_id", "price", "count", "company_id", "company_name"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheRequirementItemTime)
	return
}

// 缓冲
func getRequirementItemCacheMark(id int64) string {
	return fmt.Sprint("erp:requirement_item:id.", id)
}

func deleteRequirementItemCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getRequirementItemCacheMark(id))
}
