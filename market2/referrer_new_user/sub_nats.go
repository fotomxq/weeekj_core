package Market2ReferrerNewUser

import (
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	Market2Log "github.com/fotomxq/weeekj_core/v5/market2/log"
	Market2LogMod "github.com/fotomxq/weeekj_core/v5/market2/log/mod"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//新注册用户
	CoreNats.SubDataByteNoErr("/user/core/create_user", subNatsNewUser)
}

// 新注册用户
func subNatsNewUser(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	//日志
	appendLog := "market2 referrer new user sub new user, "
	//解析参数
	var userData UserCore.FieldsUserType
	if err := CoreNats.ReflectDataByte(data, &userData); err != nil {
		CoreLog.Error(appendLog, "get params, ", err)
		return
	}
	//获取邀请人信息
	referrerUserID := userData.Infos.GetValInt64NoBool("referrerUserID")
	if referrerUserID < 1 {
		return
	}
	//获取前置条件并加以判断
	var marketReferrerLimitMinCount int64
	// 获取是否允许重复奖励
	var MarketReferrerNewUserRepeat bool
	if userData.OrgID > 0 {
		marketReferrerLimitMinCount = OrgCore.Config.GetConfigValInt64NoErr(userData.OrgID, "MarketReferrerLimitMinCount")
		MarketReferrerNewUserRepeat = OrgCore.Config.GetConfigValBoolNoErr(userData.OrgID, "MarketReferrerNewUserRepeat")
	} else {
		marketReferrerLimitMinCount, _ = BaseConfig.GetDataInt64("MarketReferrerLimitMinCount")
		MarketReferrerNewUserRepeat, _ = BaseConfig.GetDataBool("MarketReferrerLimitMinCount")
	}
	if !MarketReferrerNewUserRepeat {
		//检查邀请最少限制人数
		if marketReferrerLimitMinCount < 1 {
			return
		} else {
			if marketReferrerLimitMinCount > 1 {
				//获取当前的邀请人数
				referrerCount := Market2Log.GetLogCountByUserID("referrer_new_user", 0, userData.OrgID, -1, referrerUserID)
				if referrerCount+1 < marketReferrerLimitMinCount {
					return
				}
			}
		}
	} else {
		//检查邀请最少限制人数
		if marketReferrerLimitMinCount < 1 {
			return
		} else {
			if marketReferrerLimitMinCount > 1 {
				//获取当前的邀请人数
				referrerCount := Market2Log.GetLogCountByUserID("referrer_new_user", 0, userData.OrgID, -1, referrerUserID)
				if referrerCount%marketReferrerLimitMinCount != 0 {
					return
				}
			}
		}
	}
	//获取相关奖励配置
	var marketReferrerNewUserUserIntegral, marketReferrerNewUserDepositPrice, marketReferrerNewUserTicketConfigID, marketReferrerNewUserTicketCount, marketReferrerNewUserUserSubAddHour, marketReferrerNewUserDepositPriceAllLimit int64
	var marketReferrerNewUserDepositType string
	if userData.OrgID > 0 {
		marketReferrerNewUserUserIntegral = OrgCore.Config.GetConfigValInt64NoErr(userData.OrgID, "MarketReferrerNewUserUserIntegral")
		marketReferrerNewUserDepositType = OrgCore.Config.GetConfigValNoErr(userData.OrgID, "MarketReferrerNewUserDepositType")
		marketReferrerNewUserDepositPrice = OrgCore.Config.GetConfigValInt64NoErr(userData.OrgID, "MarketReferrerNewUserDepositPrice")
		marketReferrerNewUserTicketConfigID = OrgCore.Config.GetConfigValInt64NoErr(userData.OrgID, "MarketReferrerNewUserTicketConfigID")
		marketReferrerNewUserTicketCount = OrgCore.Config.GetConfigValInt64NoErr(userData.OrgID, "MarketReferrerNewUserTicketCount")
		marketReferrerNewUserUserSubAddHour = OrgCore.Config.GetConfigValInt64NoErr(userData.OrgID, "MarketReferrerNewUserUserSubAddHour")
		marketReferrerNewUserDepositPriceAllLimit = OrgCore.Config.GetConfigValInt64NoErr(userData.OrgID, "MarketReferrerNewUserDepositPriceAllLimit")
	} else {
		marketReferrerNewUserUserIntegral, _ = BaseConfig.GetDataInt64("MarketReferrerNewUserUserIntegral")
		marketReferrerNewUserDepositType = BaseConfig.GetDataStringNoErr("MarketReferrerNewUserDepositType")
		marketReferrerNewUserDepositPrice, _ = BaseConfig.GetDataInt64("MarketReferrerNewUserDepositPrice")
		marketReferrerNewUserTicketConfigID, _ = BaseConfig.GetDataInt64("MarketReferrerNewUserTicketConfigID")
		marketReferrerNewUserTicketCount, _ = BaseConfig.GetDataInt64("MarketReferrerNewUserTicketCount")
		marketReferrerNewUserUserSubAddHour, _ = BaseConfig.GetDataInt64("MarketReferrerNewUserUserSubAddHour")
		marketReferrerNewUserDepositPriceAllLimit = BaseConfig.GetDataInt64NoErr("MarketReferrerNewUserDepositPriceAllLimit")
	}
	//修正奖励金额
	var marketReferrerNewUserDepositPriceFix string
	if userData.OrgID > 0 {
		marketReferrerNewUserDepositPriceFix = OrgCore.Config.GetConfigValNoErr(userData.OrgID, "MarketReferrerNewUserDepositPriceFix")
	} else {
		marketReferrerNewUserDepositPriceFix = BaseConfig.GetDataStringNoErr("MarketReferrerNewUserDepositPriceFix")
	}
	if marketReferrerNewUserDepositPriceFix != "" {
		marketReferrerNewUserDepositPriceFixPrice := getMarketReferrerNewUserDepositPriceFix(marketReferrerNewUserDepositPriceFix)
		if marketReferrerNewUserDepositPriceFixPrice > -1 {
			marketReferrerNewUserDepositPrice = marketReferrerNewUserDepositPriceFixPrice
		}
	}
	//奖励金额限制
	if marketReferrerNewUserDepositPriceAllLimit > 0 {
		if marketReferrerNewUserDepositPrice > marketReferrerNewUserDepositPriceAllLimit {
			marketReferrerNewUserDepositPrice = 0
		}
		//递减配置
		marketReferrerNewUserDepositPriceAllLimit -= marketReferrerNewUserDepositPrice
		err := BaseConfig.UpdateByMark(&BaseConfig.ArgsUpdateByMark{
			UpdateHash: "",
			Mark:       "MarketReferrerNewUserDepositPriceAllLimit",
			Value:      fmt.Sprint(marketReferrerNewUserDepositPriceAllLimit),
		})
		if err != nil {
			CoreLog.Error(appendLog, "update MarketReferrerNewUserDepositPriceAllLimit, ", err)
			return
		}
	}
	//触发奖励行为
	Market2LogMod.AppendLog(Market2LogMod.ArgsAppendLog{
		OrgID:                userData.OrgID,
		OrgBindID:            -1,
		UserID:               referrerUserID,
		BindID:               0,
		BindUserID:           userData.ID,
		GivingUserIntegral:   marketReferrerNewUserUserIntegral,
		GivingDepositType:    marketReferrerNewUserDepositType,
		GivingDepositPrice:   marketReferrerNewUserDepositPrice,
		GivingTicketConfigID: marketReferrerNewUserTicketConfigID,
		GivingTicketCount:    marketReferrerNewUserTicketCount,
		GivingUserSubAddHour: marketReferrerNewUserUserSubAddHour,
		Action:               "referrer_new_user",
		Des:                  "邀请新用户奖励",
		Params:               nil,
		NoReplaceByFrom:      true,
		SourceOrgID:          userData.OrgID,
	})
}
