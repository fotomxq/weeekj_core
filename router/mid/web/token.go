package RouterMidWeb

import (
	BaseToken2 "gitee.com/weeekj/weeekj_core/v5/base/token2"
	"github.com/gin-gonic/gin"
)

// 获取tokenID
func getTokenID(c *gin.Context) int64 {
	tokenID, b := c.Get("tokenID")
	if !b {
		return 0
	}
	return tokenID.(int64)
}

// 获取会话数据包
func getTokenInfo(c *gin.Context) (data BaseToken2.FieldsToken) {
	tokenID := getTokenID(c)
	if tokenID < 1 {
		return
	}
	data = BaseToken2.GetByID(tokenID)
	return
}
