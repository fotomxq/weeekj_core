package ERPPurchase

import (
	"errors"
	"fmt"
	BaseApproverMod "github.com/fotomxq/weeekj_core/v5/base/approver/mod"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetPurchaseList 获取Purchase列表参数
type ArgsGetPurchaseList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetPurchaseList 获取Purchase列表
func GetPurchaseList(args *ArgsGetPurchaseList) (dataList []FieldsOrder, dataCount int64, err error) {
	dataCount, err = purchaseDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetIDQuery("org_id", args.OrgID).SetIDQuery("company_id", args.CompanyID).SetIntQuery("status", args.Status).SetSearchQuery([]string{"company_name", "remark"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getPurchaseByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetPurchaseByID 获取Purchase数据包参数
type ArgsGetPurchaseByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetPurchaseByID 获取Purchase数
func GetPurchaseByID(args *ArgsGetPurchaseByID) (data FieldsOrder, err error) {
	data = getPurchaseByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	if !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsOrder{}
		err = errors.New("no data")
		return
	}
	return
}

// ArgsCreatePurchase 创建Purchase参数
type ArgsCreatePurchase struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//提交用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
}

// CreatePurchase 创建Purchase
func CreatePurchase(args *ArgsCreatePurchase) (id int64, err error) {
	//创建数据
	id, err = purchaseDB.Insert().SetFields([]string{"status", "org_id", "org_bind_id", "user_id", "company_id", "company_name", "remark"}).Add(map[string]any{
		"status":       0,
		"org_id":       args.OrgID,
		"org_bind_id":  args.OrgBindID,
		"user_id":      args.UserID,
		"company_id":   args.CompanyID,
		"company_name": args.CompanyName,
		"remark":       args.Remark,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//nats 通知审批
	BaseApproverMod.PushRequest("erp_order", id, BaseApproverMod.ParamsPushRequest{
		OrgID:          args.OrgID,
		OrgBindID:      args.OrgBindID,
		UserID:         args.UserID,
		ForkCode:       "default",
		ApproverRemark: args.Remark,
	})
	//反馈
	return
}

// ArgsUpdatePurchase 修改Purchase参数
type ArgsUpdatePurchase struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//提交组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindID" check:"id"`
	//提交用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//供应商公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//供应商名称
	CompanyName string `db:"company_name" json:"companyName" check:"des" min:"1" max:"300" empty:"true"`
	//备注
	Remark string `db:"remark" json:"remark" check:"des" min:"1" max:"300" empty:"true"`
}

// UpdatePurchase 修改Purchase
func UpdatePurchase(args *ArgsUpdatePurchase) (err error) {
	//更新数据
	err = purchaseDB.Update().SetFields([]string{"org_bind_id", "user_id", "company_id", "company_name", "remark"}).NeedUpdateTime().AddWhereID(args.ID).AddWhereOrgID(args.OrgID).NamedExec(map[string]any{
		"org_bind_id":  args.OrgBindID,
		"user_id":      args.UserID,
		"company_id":   args.CompanyID,
		"company_name": args.CompanyName,
		"remark":       args.Remark,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deletePurchaseCache(args.ID)
	//反馈
	return
}

// ArgsAuditPurchase 审批Purchase参数
type ArgsAuditPurchase struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//审批状态
	// 0: 未审批; 1: 审批中; 2: 审批通过; 3: 审批拒绝
	Status int `db:"status" json:"status"`
}

// AuditPurchase 审批Purchase
func AuditPurchase(args *ArgsAuditPurchase) (err error) {
	//更新数据
	err = purchaseDB.Update().SetFields([]string{"status"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"status": args.Status,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deletePurchaseCache(args.ID)
	//反馈
	return
}

// ArgsDeletePurchase 删除Purchase参数
type ArgsDeletePurchase struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeletePurchase 删除Purchase
func DeletePurchase(args *ArgsDeletePurchase) (err error) {
	//删除数据
	err = purchaseDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deletePurchaseCache(args.ID)
	//反馈
	return
}

// getPurchaseByID 通过ID获取Purchase数据包
func getPurchaseByID(id int64) (data FieldsOrder) {
	cacheMark := getPurchaseCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := purchaseDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "status", "org_id", "org_bind_id", "user_id", "company_id", "company_name", "remark"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cachePurchaseTime)
	return
}

// 缓冲
func getPurchaseCacheMark(id int64) string {
	return fmt.Sprint("erp:purchase:id.", id)
}

func deletePurchaseCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getPurchaseCacheMark(id))
}
