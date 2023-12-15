package RouterReport

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorBadRequest 反馈错误但不带JSON重新标定
func ErrorBadRequest(c *gin.Context, code string, codeMsg string) {
	res := DataType{
		Status: false,
		Code:   code,
		Msg:    codeMsg,
	}
	c.JSON(http.StatusBadRequest, &res)
	c.Abort()
}

// ErrorLog 反馈错误并抛出日志
func ErrorLog(c *gin.Context, message string, err error, code string, codeMsg string) {
	//输出日志
	if err != nil {
		CoreLog.Error(c.Request.RequestURI, " ", getLogMsg(c, message), " ", err)
	} else {
		CoreLog.Error(c.Request.RequestURI, " ", getLogMsg(c, message))
	}
	//反馈数据
	res := DataType{
		Status: false,
		Code:   code,
		Msg:    codeMsg,
	}
	c.JSON(http.StatusOK, &res)
	c.Abort()
}

// WarnLog 反馈警告并抛出错误
func WarnLog(c *gin.Context, message string, err error, code string, codeMsg string) {
	//输出日志
	if err != nil {
		CoreLog.Warn(c.Request.RequestURI, " ", getLogMsg(c, message), " ", err)
	} else {
		CoreLog.Warn(c.Request.RequestURI, " ", getLogMsg(c, message))
	}
	//反馈数据
	res := DataType{
		Status: false,
		Code:   code,
		Msg:    codeMsg,
	}
	c.JSON(http.StatusOK, &res)
	c.Abort()
}
