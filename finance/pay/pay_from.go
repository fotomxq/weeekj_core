package FinancePay

import (
	"fmt"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
)

// GetPayFrom 根据渠道来源，获取支付渠道同一化字符串
func GetPayFrom(payChannel CoreSQLFrom.FieldsFrom) string {
	return fmt.Sprint(payChannel.System, "_", payChannel.Mark)
}
