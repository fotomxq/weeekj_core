package ServiceOrder

import (
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	"time"
)

// 完成订单的收尾处理
func subNatsOrderUpdateFinish(logAppend string, orderData *FieldsOrder) {
	//获取订单来源渠道的统计转化值
	orderSystemMarkKey := getOrderSystemMarkKey(orderData.SystemMark)
	//统计订单数据
	// 记录订单数量
	AnalysisAny2.AppendData("add", "service_order_finish_count", time.Time{}, orderData.OrgID, orderData.UserID, 0, orderSystemMarkKey, 0, 1)
	// 记录订单金额
	AnalysisAny2.AppendData("add", "service_order_finish_pay_price", time.Time{}, orderData.OrgID, orderData.UserID, 0, orderSystemMarkKey, 0, orderData.Price)
	// 记录订单关联的采购商（客户公司）信息
	if orderData.CompanyID > 0 {
		// 订单数量
		AnalysisAny2.AppendData("add", "service_order_company_client_count", time.Time{}, orderData.OrgID, orderData.UserID, orderData.CompanyID, 0, 0, 1)
		// 订单金额
		AnalysisAny2.AppendData("add", "service_order_company_client_price", time.Time{}, orderData.OrgID, orderData.UserID, orderData.CompanyID, 0, 0, orderData.Price)
	}
}
