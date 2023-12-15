package RouterReport

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	"github.com/gin-gonic/gin"
)

// DataList 通用反馈列表方案
func DataList(c *gin.Context, errMessage string, codeMsg string, err error, dataList interface{}, dataCount int64) {
	if err != nil {
		BaseError(c, "data_empty", codeMsg)
		CoreLog.Warn(c.Request.RequestURI, ", ", getLogMsg(c, errMessage), err)
		return
	}
	BaseDataList(c, dataCount, dataList)
}
