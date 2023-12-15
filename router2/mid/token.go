package Router2Mid

import (
	BaseToken2 "github.com/fotomxq/weeekj_core/v5/base/token2"
	"github.com/gin-gonic/gin"
)

// GetTokenID 获取tokenID
func GetTokenID(c *gin.Context) int64 {
	tokenID, b := c.Get("tokenID")
	if !b {
		return 0
	}
	return tokenID.(int64)
}

// GetTokenInfo 获取会话数据包
func GetTokenInfo(c *gin.Context) (data BaseToken2.FieldsToken) {
	tokenID := GetTokenID(c)
	if tokenID < 1 {
		return
	}
	data = BaseToken2.GetByID(tokenID)
	return
}
