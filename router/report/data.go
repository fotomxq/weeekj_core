package RouterReport

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"github.com/gin-gonic/gin"
)

// Data 通用反馈单一数据
func Data(c *gin.Context, errMessage string, codeMsg string, err error, data interface{}) {
	if err != nil {
		BaseError(c, "data_empty", codeMsg)
		CoreLog.Warn(c.Request.RequestURI, ", ", getLogMsg(c, errMessage), err)
		return
	}
	BaseData(c, data)
}

func DataNoErr(c *gin.Context, errMessage string, codeMsg string, err error, data interface{}) {
	if err != nil {
		BaseError(c, "data_empty", codeMsg)
		return
	}
	BaseData(c, data)
}
