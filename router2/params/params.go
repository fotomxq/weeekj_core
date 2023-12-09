package Router2Params

import (
	"encoding/json"
	"fmt"
	Router2Report "gitee.com/weeekj/weeekj_core/v5/router2/mid"
)

// GetJSON 获取和验证参数
func GetJSON(context any, params interface{}) (b bool) {
	//获取路由上下文
	routerCtx := getContext(context)
	//获取参数数据包
	if err := routerCtx.ShouldBindJSON(params); err != nil {
		bodyByte, b2 := getContextBodyByte(context)
		if !b2 {
			Router2Report.ReportErrorBadRequestLog(context, "get params struct", err, "report_params_lost")
			return
		}
		if err2 := json.Unmarshal(bodyByte, params); err2 != nil {
			Router2Report.ReportErrorBadRequestLog(context, "get params struct", err2, "report_params_lost")
			return
		}
	}
	//过滤参数
	if b = CheckJSON(context, params); !b {
		return
	}
	//反馈成功
	return
}

// CheckJSON 单纯过滤数据
func CheckJSON(context any, params interface{}) (b bool) {
	//获取路由上下文
	routerCtx := getContext(context)
	//过滤参数
	var errField string
	var errCode string
	errField, errCode, b = filterParams(routerCtx, params)
	if !b {
		Router2Report.ReportErrorBadRequestLog(context, fmt.Sprint("field:", errField, ",errCode:", errCode), nil, "report_params_error")
		return
	}
	return
}
