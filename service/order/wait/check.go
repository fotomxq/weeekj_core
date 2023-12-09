package ServiceOrderWait

import (
	"errors"
	ServiceOrderWaitFields "gitee.com/weeekj/weeekj_core/v5/service/order/wait_fields"
)

// ArgsCheckOrder 检查订单推送状态参数
type ArgsCheckOrder struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，检测
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	// 可选，检测
	UserID int64 `db:"user_id" json:"userID" check:"id"`
}

// CheckOrder 检查订单推送状态
func CheckOrder(args *ArgsCheckOrder) (orderID int64, errCode, errMsg string, err error) {
	//获取订单数据
	var data ServiceOrderWaitFields.FieldsWait
	data, err = getCreateWait(args.ID)
	if err != nil || data.ID < 1 {
		//不存在则说明已经超出期限
		err = errors.New("time expire")
		return
	}
	if (args.OrgID > 0 && data.OrgID != args.OrgID) || (args.UserID > 0 && data.UserID != args.UserID) {
		err = errors.New("time expire")
		return
	}
	//存在订单说明成功
	if data.OrderID > 0 {
		//记录订单ID
		orderID = data.OrderID
		//反馈
		return
	} else {
		if data.ErrCode != "" {
			errCode = data.ErrCode
			errMsg = data.ErrMsg
			err = errors.New("error")
			return
		}
	}
	//等待移走中
	err = errors.New("data wait")
	return
}
