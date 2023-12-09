package RouterParams

import (
	"fmt"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	RouterReport "gitee.com/weeekj/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
)

// GetJSON 获取和验证参数
func GetJSON(c *gin.Context, params interface{}) (b bool) {
	//获取参数数据包
	if err := c.ShouldBindJSON(&params); err != nil {
		CoreLog.Warn(fmt.Sprint("url params is error, ", c.Request.RequestURI, ", err: ", err))
		RouterReport.ErrorBadRequest(c, "params_lost", "缺少必要的参数")
		return
	}
	//过滤参数
	if b = CheckJSON(c, params); !b {
		return
	}
	//反馈成功
	return
}

// CheckJSON 单纯过滤数据
func CheckJSON(c *gin.Context, params interface{}) (b bool) {
	//过滤参数
	var errField string
	var errCode string
	errField, errCode, b = filterParams(c, params)
	if !b {
		CoreLog.Warn(fmt.Sprint("url params is error, ", c.Request.RequestURI, ", field: ", errField, ", errCode: ", errCode))
		RouterReport.ErrorBadRequest(c, "params_error", fmt.Sprint("参数", errField, "验证[", errCode, "]失败"))
		return
	}
	return
}
