package FinanceReturnedMoney

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceCompany "gitee.com/weeekj/weeekj_core/v5/service/company"
)

type ArgsSetCompany struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//催款间隔月份
	NeedTakeAddMonth int `db:"need_take_add_month" json:"needTakeAddMonth" check:"intThan0" empty:"true"`
	//标记的应回款开始日
	// 支持: 0 月初、1-28对应日、-1 月底模式
	NeedTakeStartDay int `db:"need_take_start_day" json:"needTakeStartDay" check:"intThan0" empty:"true"`
	//每个月几号回款
	// 支持: 0 月初、1-28对应日、-1 月底模式
	NeedTakeDay int `db:"need_take_day" json:"needTakeDay" check:"intThan0" empty:"true"`
	//销售人员
	SellOrgBindID int64 `db:"sell_org_bind_id" json:"sellOrgBindID" check:"id" empty:"true"`
	//催收负责人
	ReturnOrgBindID int64 `db:"return_org_bind_id" json:"returnOrgBindID" check:"id" empty:"true"`
	//是否坏账
	IsBan bool `db:"is_ban" json:"isBan" check:"bool"`
	//催款路线
	ReturnLocation string `db:"return_location" json:"returnLocation" check:"des" min:"0" max:"600" empty:"true"`
}

func SetCompany(args *ArgsSetCompany) (err error) {
	//锁定
	setCompanyLock.Lock()
	defer setCompanyLock.Unlock()
	//获取公司
	companyData, _ := ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
		ID:    args.CompanyID,
		OrgID: args.OrgID,
	})
	if companyData.ID < 1 || !CoreFilter.EqID2(args.OrgID, companyData.OrgID) {
		err = errors.New("no company data")
		return
	}
	args.OrgID = companyData.OrgID
	//设置的回款时间不能
	//查询公司信息
	var data FieldsCompany
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, company_id FROM finance_returned_money_company WHERE company_id = $1 AND ($2 < 1 OR org_id = $2)", args.CompanyID, args.OrgID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_returned_money_company SET update_at = NOW(), delete_at = to_timestamp(0), need_take_add_month = :need_take_add_month, need_take_start_day = :need_take_start_day, need_take_day = :need_take_day, sell_org_bind_id = :sell_org_bind_id, return_org_bind_id = :return_org_bind_id, is_ban = :is_ban, return_location = :return_location WHERE id = :id", map[string]interface{}{
			"id":                  data.ID,
			"need_take_add_month": args.NeedTakeAddMonth,
			"need_take_start_day": args.NeedTakeStartDay,
			"need_take_day":       args.NeedTakeDay,
			"sell_org_bind_id":    args.SellOrgBindID,
			"return_org_bind_id":  args.ReturnOrgBindID,
			"is_ban":              args.IsBan,
			"return_location":     args.ReturnLocation,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update finance_returned_money_company data, ", err))
			return
		}
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "finance_returned_money_company", "INSERT INTO finance_returned_money_company (org_id, company_id, need_take_add_month, need_take_start_day, need_take_day, sell_org_bind_id, return_org_bind_id, is_ban, return_location, return_status) VALUES (:org_id, :company_id, :need_take_add_month, :need_take_start_day, :need_take_day, :sell_org_bind_id, :return_org_bind_id, :is_ban, :return_location, 0)", map[string]interface{}{
			"org_id":              args.OrgID,
			"company_id":          companyData.ID,
			"need_take_add_month": args.NeedTakeAddMonth,
			"need_take_start_day": args.NeedTakeStartDay,
			"need_take_day":       args.NeedTakeDay,
			"sell_org_bind_id":    args.SellOrgBindID,
			"return_org_bind_id":  args.ReturnOrgBindID,
			"is_ban":              args.IsBan,
			"return_location":     args.ReturnLocation,
		}, &data)
		if err != nil {
			err = errors.New(fmt.Sprint("insert finance_returned_money_company data, ", err))
			return
		}
	}
	deleteCompanyCache(data.CompanyID)
	data = getCompanyID(data.CompanyID)
	//修改最新一期的数据
	var margeData FieldsMarge
	err = Router2SystemConfig.MainDB.Get(&margeData, "SELECT id FROM finance_returned_money_marge WHERE delete_at < to_timestamp(1000000) AND org_id = $1 AND company_id = $2 ORDER BY need_at DESC LIMIT 1", data.OrgID, data.CompanyID)
	if err == nil && margeData.ID > 0 {
		//更新数据
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_returned_money_marge SET update_at = NOW(), sell_org_bind_id = :sell_org_bind_id, return_org_bind_id = :return_org_bind_id WHERE id = :id", map[string]interface{}{
			"id":                 margeData.ID,
			"sell_org_bind_id":   data.SellOrgBindID,
			"return_org_bind_id": data.ReturnOrgBindID,
		})
		if err != nil {
			return
		}
		//删除缓冲
		deleteMargeCache(margeData.ID)
	} else {
		err = nil
	}
	//反馈
	return
}
