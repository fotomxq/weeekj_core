package OrgCert2

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	FinancePayCreate "github.com/fotomxq/weeekj_core/v5/finance/pay_create"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsPayCert 支付费用参数
type ArgsPayCert struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	PaymentChannel CoreSQLFrom.FieldsFrom `db:"payment_channel" json:"paymentChannel"`
	//支付备注
	// 用户环节可根据实际业务需求开放此项
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// PayCert 支付费用
func PayCert(args *ArgsPayCert) (payData FinancePay.FieldsPayType, errCode string, err error) {
	//获取请求
	var certData FieldsCert
	certData = getCertByID(args.ID)
	if certData.ID < 1 || CoreSQL.CheckTimeHaveData(certData.DeleteAt) {
		errCode = "cert_not_exist"
		err = errors.New("cert not exist")
		return
	}
	//构建支付请求
	payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
		UserID:         args.UserID,
		OrgID:          args.OrgID,
		IsRefund:       false,
		Currency:       certData.Currency,
		Price:          certData.Price,
		PaymentChannel: args.PaymentChannel,
		ExpireAt:       CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
		Des:            args.Des,
	})
	if err != nil {
		return
	}
	//计算支付方式
	payFromSystem := fmt.Sprint(payData.PaymentChannel.System)
	if payData.PaymentChannel.Mark != "" {
		payFromSystem = payFromSystem + "_" + payData.PaymentChannel.Mark
	}
	certData.Params = CoreSQLConfig.Set(certData.Params, "paySystem", payFromSystem)
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), pay_failed = false, pay_id = :pay_id, params = :params WHERE id = :id", map[string]interface{}{
		"id":     certData.ID,
		"pay_id": payData.ID,
		"params": certData.Params,
	})
	if err != nil {
		errCode = "update"
		return
	}
	deleteCertCache(certData.ID)
	return
}

// 更新缴费成功
func updateCertPayFinish(payID int64) (err error) {
	//支付ID必须大于0
	if payID < 1 {
		return
	}
	//获取缴费记录
	var data FieldsCert
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_cert2 WHERE pay_id = $1 LIMIT 1", payID)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = getCertByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	//修改缴费状态
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), pay_at = NOW() WHERE id = :id", map[string]interface{}{
		"id": data.ID,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteCertCache(data.ID)
	//获取配置
	configData := getConfigByID(data.ConfigID)
	//请求自动审核
	if configData.AuditType == "auto" {
		pushNatsAutoAudit(data.ID)
	}
	//反馈
	return
}
