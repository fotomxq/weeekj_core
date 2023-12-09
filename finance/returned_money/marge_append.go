package FinanceReturnedMoney

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	FinancePayMod "gitee.com/weeekj/weeekj_core/v5/finance/pay/mod"
	OrgWorkTipMod "gitee.com/weeekj/weeekj_core/v5/org/work_tip/mod"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	ServiceCompany "gitee.com/weeekj/weeekj_core/v5/service/company"
	"github.com/golang-module/carbon"
	"time"
)

type argsAppendMarge struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//回款公司ID
	CompanyID int64 `db:"company_id" json:"companyID"`
	//是否回款, 否则为入账
	IsReturn bool `db:"isReturn" json:"isReturn"`
	//回款金额
	Price int64 `db:"price" json:"price"`
	//支付ID
	PayID int64 `db:"pay_id" json:"payID" check:"id" empty:"true"`
}

func appendMarge(args *argsAppendMarge) (err error) {
	//锁定
	appendMargeLock.Lock()
	defer appendMargeLock.Unlock()
	//获取公司设置
	companyData := GetCompanyByCompanyID(args.CompanyID, args.OrgID)
	if companyData.ID < 1 || CoreSQL.CheckTimeHaveData(companyData.DeleteAt) || args.OrgID != companyData.OrgID {
		//自动设置公司
		_, err = ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
			ID:    args.CompanyID,
			OrgID: args.OrgID,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("service company data is empty, company id: ", args.CompanyID, ", org id: ", args.OrgID, ", err: ", err))
			return
		}
		err = SetCompany(&ArgsSetCompany{
			OrgID:            args.OrgID,
			CompanyID:        args.CompanyID,
			NeedTakeAddMonth: 0,
			NeedTakeStartDay: 0,
			NeedTakeDay:      0,
			SellOrgBindID:    0,
			ReturnOrgBindID:  0,
			IsBan:            false,
			ReturnLocation:   "",
		})
		if err != nil {
			err = errors.New(fmt.Sprint("company not exist, auto create but failed, err: ", err, ", company id: ", args.CompanyID, ", org id: ", args.OrgID))
			return
		}
		//重新获取公司设置
		companyData = GetCompanyByCompanyID(args.CompanyID, -1)
	}
	//计算金额
	var needPrice int64 = 0
	var havePrice int64 = 0
	var haveAt carbon.Carbon
	if args.IsReturn {
		havePrice += args.Price
	} else {
		needPrice += args.Price
	}
	//获取当期数据或创建
	var nowMargeData FieldsMarge
	err = Router2SystemConfig.MainDB.Get(&nowMargeData, "SELECT id FROM finance_returned_money_marge WHERE company_id = $1 AND delete_at < to_timestamp(1000000) AND need_at >= NOW()", companyData.CompanyID)
	//如果存在数据，则修改；否则创建
	if err == nil && nowMargeData.ID > 0 {
		//获取当期数据详情
		nowMargeData = getMargeByID(nowMargeData.ID)
		//累加当期数据
		needPrice += nowMargeData.NeedPrice
		havePrice += nowMargeData.HavePrice
		//修改数据
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_returned_money_marge SET update_at = NOW(), need_price = :need_price, have_price = :have_price, have_at = :have_at, sell_org_bind_id = :sell_org_bind_id, return_org_bind_id = :return_org_bind_id WHERE id = :id", map[string]interface{}{
			"id":                 nowMargeData.ID,
			"org_id":             companyData.OrgID,
			"company_id":         companyData.CompanyID,
			"need_price":         needPrice,
			"have_price":         havePrice,
			"have_at":            haveAt.Time,
			"sell_org_bind_id":   companyData.SellOrgBindID,
			"return_org_bind_id": companyData.ReturnOrgBindID,
		})
		if err != nil {
			return
		}
	} else {
		//分析时间节点
		// 关键节点时间
		nextTakeAt := getCompanyNextTime(args.OrgID, companyData.CompanyID)
		//获取开始催款时间
		needTakeStartDay := 0
		if companyData.NeedTakeStartDay < 1 {
			needTakeStartDay = 0
		} else {
			if companyData.NeedTakeStartDay > 31 {
				needTakeStartDay = 31
			}
		}
		var startAt carbon.Carbon
		if needTakeStartDay > 0 {
			startAt = nextTakeAt.SetDay(companyData.NeedTakeStartDay)
		} else {
			startAt = nextTakeAt
		}
		//创建新数据
		var newMargeID int64
		newMargeID, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_returned_money_marge (org_id, company_id, start_at, need_price, need_at, have_price, have_at, sell_org_bind_id, return_org_bind_id) VALUES (:org_id, :company_id, :start_at, :need_price, :need_at, :have_price, :have_at, :sell_org_bind_id, :return_org_bind_id)", map[string]interface{}{
			"org_id":             companyData.OrgID,
			"company_id":         companyData.CompanyID,
			"start_at":           startAt.Time,
			"need_price":         needPrice,
			"need_at":            nextTakeAt.Time,
			"have_price":         havePrice,
			"have_at":            haveAt.Time,
			"sell_org_bind_id":   companyData.SellOrgBindID,
			"return_org_bind_id": companyData.ReturnOrgBindID,
		})
		if err != nil {
			return
		}
		if newMargeID < 1 {
			err = errors.New("insert no data")
			return
		}
		//获取当期数据详情
		nowMargeData = getMargeByID(newMargeID)
	}
	//获取该公司累计未还款
	_, _, allLastPrice := GetMargeNoTakePrice(companyData.CompanyID, CoreFilter.GetNowTimeCarbon().AddYear().Time)
	//检查是否已经还款
	//检查累计总数
	if allLastPrice < 1 {
		haveAt = CoreFilter.GetNowTimeCarbon()
	} else {
		haveAt = carbon.Carbon{}
	}
	var returnConfirmAt time.Time
	if !args.IsReturn {
		returnConfirmAt = time.Time{}
	} else {
		returnConfirmAt = nowMargeData.ReturnConfirmAt
	}
	//修改后续的回款时间
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE finance_returned_money_marge SET update_at = NOW(), have_at = :have_at, return_confirm_at = :return_confirm_at WHERE id = :id", map[string]interface{}{
		"id":                nowMargeData.ID,
		"have_at":           haveAt.Time,
		"return_confirm_at": returnConfirmAt,
	})
	if err != nil {
		return
	}
	//发送提醒消息
	if !args.IsReturn {
		var serviceCompanyData ServiceCompany.FieldsCompany
		serviceCompanyData, _ = ServiceCompany.GetCompanyID(&ServiceCompany.ArgsGetCompanyID{
			ID:    companyData.CompanyID,
			OrgID: -1,
		})
		if companyData.OrgID > 0 && companyData.ReturnOrgBindID > 0 {
			OrgWorkTipMod.AppendTip(&OrgWorkTipMod.ArgsAppendTip{
				OrgID:     companyData.OrgID,
				OrgBindID: companyData.ReturnOrgBindID,
				Msg:       fmt.Sprint(serviceCompanyData.Name, "需要催收一笔款项，请及时确认"),
				System:    "finance_returned_money",
				BindID:    nowMargeData.ID,
			})
		}
	}
	//删除缓冲
	deleteMargeCache(nowMargeData.ID)
	//获取公司统计信息，用于更新公司付款状态
	analysisCompanyData := GetMargeAnalysisCompany(companyData.OrgID, companyData.CompanyID)
	var companyResultStatus = 0
	if analysisCompanyData.AllLastPrice < 1 {
		if analysisCompanyData.AllReturnPrice > 0 {
			companyResultStatus = 3
		} else {
			companyResultStatus = 0
		}
	} else {
		if analysisCompanyData.Last1Price < 1 {
			companyResultStatus = 1
		} else {
			if analysisCompanyData.Last30Price < 1 {
				companyResultStatus = 4
			} else {
				if analysisCompanyData.Last60Price < 1 {
					companyResultStatus = 5
				} else {
					if analysisCompanyData.Last90Price < 1 {
						companyResultStatus = 6
					} else {
						if analysisCompanyData.Last365Price < 1 {
							companyResultStatus = 7
						} else {
							companyResultStatus = 8
						}
					}
				}
			}
		}
	}
	if companyData.ReturnStatus != companyResultStatus {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_returned_money_company SET update_at = NOW(), return_status = :return_status WHERE id = :id", map[string]interface{}{
			"id":            companyData.ID,
			"return_status": companyResultStatus,
		})
		if err != nil {
			CoreLog.Warn("finance returned money update company return status failed, company id: ", companyData.CompanyID, ", dest return status: ", companyResultStatus, ", err: ", err)
			err = nil
		} else {
			deleteCompanyCache(companyData.CompanyID)
		}
	}
	//如果存在支付ID，则触发nats
	if args.PayID > 0 {
		FinancePayMod.PushPayFinish(args.PayID)
	}
	//反馈
	return
}
