package Router2Params

import (
	Router2Mid "gitee.com/weeekj/weeekj_core/v5/router2/mid"
	"github.com/gin-gonic/gin"
)

// 识别和获取头部上下文
func getContext(c any) *gin.Context {
	return Router2Mid.GetContext(c)
}

// getContextBodyByte 尝试获取上下文的body byte
func getContextBodyByte(c any) (dataByte []byte, b bool) {
	return Router2Mid.GetContextBodyByte(c)
}
