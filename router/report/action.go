package RouterReport

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"github.com/gin-gonic/gin"
)

//反馈action

func ActionCreateNoData(c *gin.Context, logMsg string, codeMsg string, err error) {
	if err != nil {
		BaseError(c, "create_failed", codeMsg)
		CoreLog.Warn(c.Request.RequestURI, ", ", getLogMsg(c, logMsg), err)
		return
	}
	BaseSuccess(c)
}

func ActionCreate(c *gin.Context, logMsg string, codeMsg string, err error, data interface{}) {
	if err != nil {
		BaseError(c, "create_failed", codeMsg)
		CoreLog.Warn(c.Request.RequestURI, ", ", getLogMsg(c, logMsg), err)
		return
	}
	BaseData(c, data)
}

func ActionUpdate(c *gin.Context, logMsg string, codeMsg string, err error) {
	if err != nil {
		BaseError(c, "update_failed", codeMsg)
		CoreLog.Warn(c.Request.RequestURI, ", ", getLogMsg(c, logMsg), err)
		return
	}
	BaseSuccess(c)
}

func ActionDelete(c *gin.Context, logMsg string, codeMsg string, err error) {
	if err != nil {
		BaseError(c, "delete_failed", codeMsg)
		CoreLog.Warn(c.Request.RequestURI, ", ", getLogMsg(c, logMsg), err)
		return
	}
	BaseSuccess(c)
}
