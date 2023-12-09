package FinanceReturnedMoney

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	OrgWorkTipMod "gitee.com/weeekj/weeekj_core/v5/org/work_tip/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceCompany "gitee.com/weeekj/weeekj_core/v5/service/company"
)

// ArgsUpdateMargeConfirm 确认催收款参数
type ArgsUpdateMargeConfirm struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//催收负责人
	ReturnOrgBindID int64 `db:"return_org_bind_id" json:"returnOrgBindID"`
}

// UpdateMargeConfirm 确认催收款
func UpdateMargeConfirm(args *ArgsUpdateMargeConfirm) (err error) {
	//获取融合表信息
	margeData := getMargeByID(args.ID)
	if margeData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//获取催收人信息
	if margeData.ReturnOrgBindID < 1 || margeData.ReturnOrgBindID != args.ReturnOrgBindID {
		err = errors.New("no data")
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_returned_money_marge SET update_at = NOW(), return_confirm_at = NOW() WHERE id = :id", map[string]interface{}{
		"id": margeData.ID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteMargeCache(args.ID)
	//反馈
	return
}

// ArgsSendMargeReturn 催促收款单参数
type ArgsSendMargeReturn struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// SendMargeReturn 催促收款单
func SendMargeReturn(args *ArgsSendMargeReturn) (err error) {
	//获取融合表信息
	margeData := getMargeByID(args.ID)
	if margeData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//获取催收人信息
	var serviceCompanyData ServiceCompany.FieldsCompany
	serviceCompanyData, _ = ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
		ID:    margeData.CompanyID,
		OrgID: -1,
	})
	if margeData.ReturnOrgBindID > 0 && serviceCompanyData.ID > 0 && margeData.OrgID > 0 {
		OrgWorkTipMod.AppendTip(&OrgWorkTipMod.ArgsAppendTip{
			OrgID:     margeData.OrgID,
			OrgBindID: margeData.ReturnOrgBindID,
			Msg:       fmt.Sprint(serviceCompanyData.Name, "需要催收一笔款项，请尽快完成催收处理！"),
			System:    "finance_returned_money",
			BindID:    margeData.ID,
		})
	}
	//反馈
	return
}
