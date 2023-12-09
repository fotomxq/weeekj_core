package FinanceReturnedMoney

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
	"github.com/tidwall/gjson"
)

func subNats() {
	//新增回款记录
	CoreNats.SubDataByteNoErr("/finance/return_money/log", subNatsNewLog)
}

// 新增回款日志
func subNatsNewLog(_ *nats.Msg, action string, logID int64, _ string, data []byte) {
	//检查参数
	if logID < 1 {
		return
	}
	orgID := gjson.GetBytes(data, "orgID").Int()
	companyID := gjson.GetBytes(data, "companyID").Int()
	isReturn := gjson.GetBytes(data, "isReturn").Bool()
	price := gjson.GetBytes(data, "price").Int()
	payID := gjson.GetBytes(data, "payID").Int()
	//检查行为模式
	switch action {
	case "new":
		if err := appendMarge(&argsAppendMarge{
			OrgID:     orgID,
			CompanyID: companyID,
			IsReturn:  isReturn,
			Price:     price,
			PayID:     payID,
		}); err != nil {
			CoreLog.Error("finance return money sub nats new log, append marge, ", err)
		}
	}
}
