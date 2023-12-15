package FinanceReturnedMoney

import (
	"errors"
	"fmt"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ServiceCompany "github.com/fotomxq/weeekj_core/v5/service/company"
)

// ArgsAppendLog 添加一个回款记录参数
type ArgsAppendLog struct {
	//回款单号
	SN string `db:"sn" json:"sn"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//关联订单ID
	OrderID int64 `db:"order_id" json:"orderID" check:"id" empty:"true"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID" check:"id" empty:"true"`
	//关联其他第三方模块
	BindSystem string `db:"bind_system" json:"bindSystem" check:"mark" empty:"true"`
	BindID     int64  `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	BindMark   string `db:"bind_mark" json:"bindMark"`
	//是否回款, 否则为入账
	IsReturn bool `db:"is_return" json:"isReturn" check:"bool" empty:"true"`
	//回款金额
	Price int64 `db:"price" json:"price" check:"price"`
	//备注历史
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// AppendLog 添加一个回款记录
// 如果needAt在未来时间，则不断叠加费用
func AppendLog(args *ArgsAppendLog) (errCode string, err error) {
	//获取公司
	companyData, _ := ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
		ID:    args.CompanyID,
		OrgID: args.OrgID,
	})
	//if companyData.ID < 1 || !CoreFilter.EqID2(args.OrgID, companyData.OrgID) {
	if companyData.ID < 1 {
		errCode = "err_pay_no_company"
		err = errors.New(fmt.Sprint("no company data, company id: ", args.CompanyID, ", find id: ", companyData.ID, ", org id: ", args.OrgID, ", find id: ", companyData.OrgID))
		return
	}
	args.OrgID = companyData.OrgID
	//检查是否为退款处理
	haveReturn := false
	if args.IsReturn {
		if args.PayID > 0 {
			var findHaveReturnID int64
			_ = Router2SystemConfig.MainDB.Get(&findHaveReturnID, "SELECT id FROM finance_returned_money_log WHERE company_id = $1 AND pay_id = $2 AND is_return = false LIMIT 1", args.CompanyID, args.PayID)
			if findHaveReturnID > 0 {
				haveReturn = true
			}
			_, _ = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_returned_money_log SET have_refund = true WHERE id = :id", map[string]interface{}{
				"id": findHaveReturnID,
			})
		}
	}
	//添加到日志
	var newLogID int64
	newLogID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_returned_money_log (sn, org_id, company_id, order_id, pay_id, bind_system, bind_id, bind_mark, is_return, have_refund, price, des) VALUES (:sn, :org_id, :company_id, :order_id, :pay_id, :bind_system, :bind_id, :bind_mark, :is_return, :have_refund, :price, :des)", map[string]interface{}{
		"sn":          args.SN,
		"org_id":      args.OrgID,
		"company_id":  args.CompanyID,
		"order_id":    args.OrderID,
		"pay_id":      args.PayID,
		"bind_system": args.BindSystem,
		"bind_id":     args.BindID,
		"bind_mark":   args.BindMark,
		"is_return":   args.IsReturn,
		"have_refund": haveReturn,
		"price":       args.Price,
		"des":         args.Des,
	})
	if err != nil {
		errCode = "err_insert"
		return
	}
	//触发消息列队，加入回款汇总表
	CoreNats.PushDataNoErr("/finance/return_money/log", "new", newLogID, "", map[string]interface{}{
		"orgID":     companyData.OrgID,
		"companyID": companyData.ID,
		"isReturn":  args.IsReturn,
		"price":     args.Price,
		"payID":     args.PayID,
	})
	//反馈
	return
}
