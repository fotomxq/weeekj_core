package FinancePay

import (
	"errors"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsCheckFinishByIDs 检查一组是否完成？参数
type ArgsCheckFinishByIDs struct {
	//一组ID
	IDs []int64 `json:"ids"`
}

// DataCheckFinish 检查一组是否完成？数据
type DataCheckFinish struct {
	//ID
	ID int64 `db:"id" json:"id" check:"ids"`
	//当前状态
	Status int `db:"status" json:"status"`
	//是否完成
	IsFinish bool `db:"is_finish" json:"isFinish"`
	//失败代码
	FailedCode string `db:"failed_code" json:"failedCode"`
	//失败消息
	FailedMessage string `db:"failed_message" son:"failedMessage"`
}

// CheckFinishByIDs 检查一组是否完成？
func CheckFinishByIDs(args *ArgsCheckFinishByIDs) (dataList []DataCheckFinish, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, status, (status = 3) as is_finish, failed_code, failed_message FROM finance_pay WHERE id = ANY($1)", pq.Array(args.IDs))
	if err == nil {
		if len(dataList) < 1 {
			err = errors.New("data is empty")
			return
		}
		for k, _ := range dataList {
			checkFinishStatus(&dataList[k])
		}
	}
	return
}

type ArgsCheckFinishByID struct {
	//ID
	ID int64 `json:"id" check:"id"`
}

func CheckFinishByID(args *ArgsCheckFinishByID) (data DataCheckFinish, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, status, (status = 3) as is_finish, failed_code, failed_message FROM finance_pay WHERE id = $1", args.ID)
	if err != nil {
		return
	}
	CoreNats.PushDataNoErr("/finance/pay/finish_check_result", "", data.ID, "", data)
	return
}

func checkFinishStatus(data *DataCheckFinish) {
	switch data.Status {
	case 2:
		if data.FailedCode == "" {
			data.FailedCode = "failed"
			data.FailedMessage = "未知支付失败"
		}
	case 4:
		if data.FailedCode == "" {
			data.FailedCode = "remove"
			data.FailedMessage = "支付请求被删除"
		}
	case 5:
		if data.FailedCode == "" {
			data.FailedCode = "expire"
			data.FailedMessage = "支付时间超时"
		}
	case 6:
		if data.FailedCode == "" {
			data.FailedCode = "refund"
			data.FailedMessage = "发起了退款"
		}
	case 7:
		if data.FailedCode == "" {
			data.FailedCode = "refundAudit"
			data.FailedMessage = "正在退款中"
		}
	case 8:
		if data.FailedCode == "" {
			data.FailedCode = "refundFailed"
			data.FailedMessage = "退款发起失败"
		}
	case 9:
		if data.FailedCode == "" {
			data.FailedCode = "refundFinish"
			data.FailedMessage = "支付已经退款"
		}
	default:
		//不做处理
	}
}

// 检查商户是否开通独立财务支付权限
func checkOrgHaveFinancePayPermission(orgID int64) bool {
	return OrgCore.CheckOrgPermissionFunc(orgID, "finance_independent")
}
